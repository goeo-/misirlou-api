package models

import (
	"math/rand"
	"reflect"
	"testing"
	"time"
)

type fataller interface {
	Fatalf(format string, values ...interface{})
}

func equal(f fataller, got, want interface{}) {
	if !reflect.DeepEqual(want, got) {
		if w, ok := want.([]byte); ok {
			want = string(w)
		}
		if g, ok := got.([]byte); ok {
			got = string(g)
		}
		f.Fatalf("want %v got %v", want, got)
	}
}

func TestID(t *testing.T) {
	t.Run("Generate", func(t *testing.T) {
		n := ID(time.Now().UnixNano()) >> 19
		id := GenerateID()
		id >>= 15
		if n != id && n+1 != id {
			t.Fatalf("Timestamp is not what it should be: got %d want %d", id, n)
		}
	})
	t.Run("Binary", func(t *testing.T) {
		for i := 0; i < 5000; i++ {
			id := ID(rand.Uint64())
			b, _ := id.MarshalBinary()
			id2 := new(ID)
			id2.UnmarshalBinary(b)
			equal(t, id, *id2)
		}
	})
	t.Run("Text", func(t *testing.T) {
		for i := 0; i < 5000; i++ {
			id := ID(rand.Uint64())
			b, _ := id.MarshalText()
			id2 := new(ID)
			id2.UnmarshalText(b)
			equal(t, id, *id2)
		}
	})
	t.Run("String", func(t *testing.T) {
		const i = 9187681986
		s1 := ID(i).String()
		s2, _ := ID(i).MarshalText()
		equal(t, s1, string(s2))
	})
	t.Run("Time", func(t *testing.T) {
		want := time.Now().Truncate(time.Millisecond)
		id := GenerateID()
		got := id.Time().Truncate(time.Millisecond)
		if !got.Equal(want) && !got.Equal(want.Add(-time.Millisecond)) {
			t.Fatalf("got %v want %v", got, want)
		}
	})
}

func BenchmarkGenerateID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateID()
	}
}

func BenchmarkBinaryMarshalID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ID(1976898619).MarshalBinary()
	}
}

func BenchmarkBinaryUnmarshalID(b *testing.B) {
	id := new(ID)
	for i := 0; i < b.N; i++ {
		id.UnmarshalBinary([]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'})
	}
}

func BenchmarkTextMarshalID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ID(1976898619).MarshalText()
	}
}

func BenchmarkTextUnmarshalID(b *testing.B) {
	id := new(ID)
	for i := 0; i < b.N; i++ {
		id.UnmarshalText([]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k'})
	}
}
