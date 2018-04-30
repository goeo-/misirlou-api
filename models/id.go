package models

import (
	"encoding/base64"
	"reflect"
	"sync/atomic"
	"time"

	"github.com/jinzhu/gorm"
)

// ID is the ID of a single resource. IDs are similar to Snowflake IDs, in that
// they are composed of a timestamp and a "step number".
//
// The first 15 least significant bits are reserved to the step number. This
// is generally a number incremented with each generation of the ID, to ensure
// uniqueness. Having 15 bits as a step number effectively allows to have 32768
// unique ID generations for each timestamp.
//
// The other 49 most significant bits are a unix timestamp in nanoseconds
// divided by 2^19. This allows for a roughly 0,5 millisecond precision, aside
// from being much faster than division thanks to bitwise shifting.
// This allows the timestamp to only roll over at 9353 years from the UNIX
// epoch, or 2^49 / (1e9*60*60*24*365.25 / 2^19) years, at which point if
// humans are still around and not dead due to a global catastrophe they
// probably are not still using something written 9000 years before. To put it
// in comparison, the 7th millenium BC is when the paleolithic ended and the
// neolithic started in China.
type ID uint64

var step uint32

// these are the bits reserved to each part of the ID, used for Bitwise AND'ing.
const (
	stepBits ID = 1<<15 - 1
	timeBits ID = ^stepBits
)

// GenerateID creates a new ID.
func GenerateID() ID {
	// Get timestamp in nanoseconds and shift 4 bits to the right. Reset the
	// least signif. 15 bits of the result (they are used for the step number).
	// This is thus a division by 2^19.
	id := (ID(time.Now().UnixNano()) >> 4) & timeBits
	// atomically increment step number and retrieve it, while also removing
	// all bits that are not the last 15.
	s := ID(atomic.AddUint32(&step, 1)) & stepBits
	return id | s
}

// Binary returns the binary representation of ID, like MarshalBinary. The only
// advantage is that this does not return an error and can be used directly.
func (i ID) Binary() []byte {
	return []byte{
		byte(i >> 56),
		byte(i >> 48),
		byte(i >> 40),
		byte(i >> 32),
		byte(i >> 24),
		byte(i >> 16),
		byte(i >> 8),
		byte(i),
	}
}

// MarshalBinary converts the ID into its binary big endian form. It implements
// encoding.BinaryMarshaler.
func (i ID) MarshalBinary() ([]byte, error) {
	return i.Binary(), nil
}

// UnmarshalBinary decodes b into the ID. It assumed that b is at least 8 bytes
// long, that it is encoded in big endian, and that i is not nil. It implements
// encoding.BinaryUnmarshaler.
func (i *ID) UnmarshalBinary(b []byte) error {
	if i == nil || len(b) < 8 {
		return nil
	}
	*i = ID(b[7]) |
		ID(b[6])<<8 |
		ID(b[5])<<16 |
		ID(b[4])<<24 |
		ID(b[3])<<32 |
		ID(b[2])<<40 |
		ID(b[1])<<48 |
		ID(b[0])<<56
	return nil
}

// MarshalText encodes i into its text representation. It is a 11-byte long
// string composed only of base64 characters.
func (i ID) MarshalText() ([]byte, error) {
	src, _ := i.MarshalBinary()
	dst := make([]byte, 11)
	base64.RawURLEncoding.Encode(dst, src)
	return dst, nil
}

// UnmarshalText decodes i from its text representation. text must be at least
// 11 bytes long and i must not be nil.
func (i *ID) UnmarshalText(text []byte) error {
	if i == nil || len(text) < 11 {
		return nil
	}
	text = text[:11]
	dst := make([]byte, 8)
	_, err := base64.RawURLEncoding.Decode(dst, text)
	if err != nil {
		return err
	}
	return i.UnmarshalBinary(dst)
}

func (i ID) String() string {
	s, _ := i.MarshalText()
	return string(s)
}

// Time retrieves the time of creation of the ID, by parsing the timestamp in
// the first 49 bits of the ID.
func (i ID) Time() time.Time {
	return time.Unix(0, int64((i<<4)&timeBits))
}

// Register gorm callback for generating IDs when they are blank.
func init() {
	gorm.DefaultCallback.Create().Before("gorm:before_create").
		Register("snowflake:generate_id", generateID)
}

var idType = reflect.TypeOf(ID(0))

func generateID(scope *gorm.Scope) {
	if scope.HasError() {
		return
	}
	field, ok := scope.FieldByName("ID")
	if !ok {
		return
	}
	if !field.IsBlank {
		return
	}
	if !idType.AssignableTo(field.Field.Type()) {
		return
	}
	err := field.Set(GenerateID())
	if err != nil {
		scope.Err(err)
	}
}
