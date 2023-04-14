package registry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_registryNew(t *testing.T) {
	r := New([]IOption{
		OptId("testid"),
		OptAddr("127.0.0.1:7370"),
		OptAdvertise("127.0.0.1:7370"),
		OptRegistries(""),
		OptApiAddr("127.0.0.1:9000"),
		OptService("127.0.0.1:9000"),
	})

	assert.NotNil(t, r)
	assert.Equal(t, "127.0.0.1:7370", r.opt.Addr)
	assert.Equal(t, "127.0.0.1:7370", r.opt.Advertise)
	assert.Equal(t, "", r.opt.Registries)
	assert.Equal(t, "127.0.0.1:9000", r.opt.ApiAddr)
	assert.Equal(t, "127.0.0.1:9000", r.opt.Service)
	assert.NotNil(t, r.serf)
	assert.NotNil(t, r.api)
}

func Test_registryServe(t *testing.T) {
	r1 := New([]IOption{
		OptId("testid"),
		OptAddr("127.0.0.1:7370"),
		OptAdvertise("127.0.0.1:7370"),
		OptRegistries(""),
		OptApiAddr("127.0.0.1:9000"),
		OptService("127.0.0.1:9000"),
	})
	err := r1.Serve()
	assert.Nil(t, err)
	r1.Close()
}

func Test_registryServeErr(t *testing.T) {
	r1 := New([]IOption{
		OptId("testid"),
		OptAddr("127.0.0.1"),
		OptAdvertise("127.0.0.1:7370"),
		OptRegistries(""),
		OptApiAddr("127.0.0.1:9000"),
		OptService("127.0.0.1:9000"),
	})
	err := r1.Serve()
	assert.NotNil(t, err)
	assert.Equal(t, ErrParseAddrToHostPort, err)
	r1.Close()

	r2 := New([]IOption{
		OptId("testid"),
		OptAddr("127.0.0.1:abcd"),
		OptAdvertise("127.0.0.1:7370"),
		OptRegistries(""),
		OptApiAddr("127.0.0.1:9000"),
		OptService("127.0.0.1:9000"),
	})
	err = r2.Serve()
	assert.NotNil(t, err)
	assert.Equal(t, ErrParsePort, err)
	r2.Close()
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
