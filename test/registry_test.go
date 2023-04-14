package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/werbenhu/registry"
)

func Test_registryNew(t *testing.T) {
	r := registry.New([]registry.IOption{
		registry.OptId("testid"),
		registry.OptBind("127.0.0.1:7370"),
		registry.OptBindAdvertise("127.0.0.1:7370"),
		registry.OptRegistries(""),
		registry.OptAddr("127.0.0.1:9000"),
		registry.OptAdvertise("127.0.0.1:9000"),
	})

	assert.NotNil(t, r)
}

func Test_registryServe(t *testing.T) {
	r := registry.New([]registry.IOption{
		registry.OptId("testid"),
		registry.OptBind("127.0.0.1:7370"),
		registry.OptBindAdvertise("127.0.0.1:7370"),
		registry.OptRegistries(""),
		registry.OptAddr("127.0.0.1:9000"),
		registry.OptAdvertise("127.0.0.1:9000"),
	})
	err := r.Serve()
	assert.Nil(t, err)
	r.Close()
}

func Test_registryServeErr(t *testing.T) {
	r := registry.New([]registry.IOption{
		registry.OptId("testid"),
		registry.OptBind("127.0.0.1:abcd"),
		registry.OptBindAdvertise("127.0.0.1:7370"),
		registry.OptRegistries(""),
		registry.OptAddr("127.0.0.1:9000"),
		registry.OptAdvertise("127.0.0.1:9000"),
	})
	err := r.Serve()
	assert.NotNil(t, err)
	assert.Equal(t, registry.ErrParsePort, err)
	r.Close()
}

func Test_registryOnMemberJoin(t *testing.T) {
}

func Test_registryOnMemberLeave(t *testing.T) {
}

func Test_registryOnMemberUpdate(t *testing.T) {
}

func Test_SRouteMatch(t *testing.T) {
}

func Test_SRouteMembers(t *testing.T) {
}
