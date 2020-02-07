package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	service, err := New("testservice", "develop-local")
	assert.NoError(t, err)

	pong, err := service.RedisClient.Ping().Result()
	assert.Equal(t, "PONG", pong)
	assert.NoError(t, err)

	err = service.RegisterService(func() error {
		return nil
	})
	assert.NoError(t, err)
}
