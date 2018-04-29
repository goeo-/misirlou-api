package ripple_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"zxq.co/x/ripple"
)

var c = &ripple.Client{
	Token: os.Getenv("RIPPLE_TOKEN"),
}

func tokenskip(t *testing.T) {
	if os.Getenv("RIPPLE_TOKEN") == "" {
		t.Skip("test requires token, skipping...")
	}
}

func TestUser(t *testing.T) {
	t.Run("Nyo", func(t *testing.T) {
		u, err := c.User(1009)
		assert.NoError(t, err)
		assert.Equal(t, u.ID, 1009)
		assert.Equal(t, u.Country, "IT")
	})
	t.Run("NotFound", func(t *testing.T) {
		u, err := c.User(-1009)
		assert.NoError(t, err)
		assert.Nil(t, u)
	})
	t.Run("Self", func(t *testing.T) {
		tokenskip(t)
		u, err := c.User(ripple.Self)
		assert.NoError(t, err)
		assert.NotNil(t, u)
		t.Log(u)
	})
}
