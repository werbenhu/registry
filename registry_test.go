package registry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_registryNew(t *testing.T) {
	r := New([]IOption{
		OptId("testid"),
		OptBind("127.0.0.1:7370"),
		OptBindAdvertise("127.0.0.1:7370"),
		OptRegistries(""),
		OptAddr("127.0.0.1:9000"),
		OptAdvertise("127.0.0.1:9000"),
	})

	assert.NotNil(t, r)
	assert.Equal(t, "127.0.0.1:7370", r.opt.Bind)
	assert.Equal(t, "127.0.0.1:7370", r.opt.BindAdvertise)
	assert.Equal(t, "", r.opt.Registries)
	assert.Equal(t, "127.0.0.1:9000", r.opt.Addr)
	assert.Equal(t, "127.0.0.1:9000", r.opt.Advertise)
	assert.NotNil(t, r.serf)
	assert.NotNil(t, r.api)
}

func Test_registryServe(t *testing.T) {
	r1 := New([]IOption{
		OptId("testid"),
		OptBind("127.0.0.1:7370"),
		OptBindAdvertise("127.0.0.1:7370"),
		OptRegistries(""),
		OptAddr("127.0.0.1:9000"),
		OptAdvertise("127.0.0.1:9000"),
	})
	err := r1.Serve()
	assert.Nil(t, err)
	r1.Close()
}

func Test_registryServeErr(t *testing.T) {
	r1 := New([]IOption{
		OptId("testid"),
		OptBind("127.0.0.1:7370"),
		OptBindAdvertise("127.0.0.1:7370"),
		OptRegistries(""),
		OptAddr("127.0.0.1:9000"),
		OptAdvertise("127.0.0.1:9000"),
	})
	err := r1.Serve()
	assert.NotNil(t, err)
	assert.Equal(t, ErrParseAddrToHostPort, err)
	r1.Close()

	r2 := New([]IOption{
		OptId("testid"),
		OptBind("127.0.0.1:abcd"),
		OptBindAdvertise("127.0.0.1:7370"),
		OptRegistries(""),
		OptAddr("127.0.0.1:9000"),
		OptAdvertise("127.0.0.1:9000"),
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
