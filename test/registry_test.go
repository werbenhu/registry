package test

import (
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/werbenhu/chash"
	"github.com/werbenhu/registry"
)

var (
	sleepTime = time.Millisecond * 100
)

func Test_RegistryNew(t *testing.T) {
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

func Test_RegistryServe(t *testing.T) {
	r := registry.New([]registry.IOption{
		registry.OptId("testid"),
		registry.OptBind("127.0.0.1:7370"),
		registry.OptBindAdvertise("127.0.0.1:7370"),
		registry.OptRegistries(""),
		registry.OptAddr("127.0.0.1:9000"),
		registry.OptAdvertise("127.0.0.1:9000"),
	})
	go r.Serve()
	time.Sleep(sleepTime)
	r.Close()
}

func Test_RegistryServeErr(t *testing.T) {
	r := registry.New([]registry.IOption{
		registry.OptId("testid"),
		registry.OptBind("127.0.0.1:abcd"),
		registry.OptBindAdvertise("127.0.0.1:7370"),
		registry.OptRegistries(""),
		registry.OptAddr("127.0.0.1:9000"),
		registry.OptAdvertise("127.0.0.1:9000"),
	})

	assert.Panics(t, func() {
		r.Serve()
	})
	r.Close()
}

func Test_RegistryOnMemberJoin(t *testing.T) {
	r := registry.New([]registry.IOption{
		registry.OptId("testid"),
		registry.OptBind("127.0.0.1:7370"),
		registry.OptBindAdvertise("127.0.0.1:7370"),
		registry.OptRegistries(""),
		registry.OptAddr("127.0.0.1:9000"),
		registry.OptAdvertise("127.0.0.1:9000"),
	})
	go r.Serve()
	time.Sleep(sleepTime)

	serviceId := "testid"
	serviceGroup := "testgroup"
	serviceAddr := "127.0.0.1:80"

	member := registry.NewMember(
		serviceId,
		"127.0.0.1:8370",
		"127.0.0.1:8370",
		"127.0.0.1:7370",
		serviceGroup,
		serviceAddr,
	)
	err := r.OnMemberJoin(member)
	assert.Nil(t, err)

	service, err := r.Match(serviceGroup, "xxx")
	assert.Nil(t, err)
	assert.NotNil(t, service)
	assert.Equal(t, serviceId, service.Id)
	assert.Equal(t, serviceGroup, service.Group)
	assert.Equal(t, serviceAddr, service.Addr)
	r.Close()
}

func Test_RegistryOnMemberLeave(t *testing.T) {
	r := registry.New([]registry.IOption{
		registry.OptId("testid"),
		registry.OptBind("127.0.0.1:7370"),
		registry.OptBindAdvertise("127.0.0.1:7370"),
		registry.OptRegistries(""),
		registry.OptAddr("127.0.0.1:9000"),
		registry.OptAdvertise("127.0.0.1:9000"),
	})
	go r.Serve()
	time.Sleep(sleepTime)

	serviceId := "testid"
	serviceGroup := "testgroup"
	serviceAddr := "127.0.0.1:80"

	member := registry.NewMember(
		serviceId,
		"127.0.0.1:8370",
		"127.0.0.1:8370",
		"127.0.0.1:7370",
		serviceGroup,
		serviceAddr,
	)
	err := r.OnMemberJoin(member)
	assert.Nil(t, err)

	err = r.OnMemberLeave(member)
	assert.Nil(t, err)

	service, err := r.Match(serviceGroup, "xxx")
	assert.Nil(t, service)
	assert.Equal(t, chash.ErrNoResultMatched, err)
	r.Close()
}

func Test_RegistryOnMemberUpdate(t *testing.T) {
	r := registry.New([]registry.IOption{
		registry.OptId("testid"),
		registry.OptBind("127.0.0.1:7370"),
		registry.OptBindAdvertise("127.0.0.1:7370"),
		registry.OptRegistries(""),
		registry.OptAddr("127.0.0.1:9000"),
		registry.OptAdvertise("127.0.0.1:9000"),
	})
	go r.Serve()
	time.Sleep(sleepTime)

	serviceId := "testid"
	serviceGroup := "testgroup"
	serviceAddr := "127.0.0.1:80"

	member := registry.NewMember(
		serviceId,
		"127.0.0.1:8370",
		"127.0.0.1:8370",
		"127.0.0.1:7370",
		serviceGroup,
		serviceAddr,
	)

	err := r.OnMemberJoin(member)
	assert.Nil(t, err)
	service, err := r.Match(serviceGroup, "xxx")
	assert.Nil(t, err)
	assert.Equal(t, serviceAddr, service.Addr)

	member.Service.Addr = "127.0.0.1:81"
	err = r.OnMemberUpdate(member)
	assert.Nil(t, err)

	service, err = r.Match(serviceGroup, "xxx")
	assert.Nil(t, err)
	assert.NotNil(t, service)
	assert.Equal(t, serviceId, service.Id)
	assert.Equal(t, serviceGroup, service.Group)
	assert.Equal(t, "127.0.0.1:81", service.Addr)
	r.Close()
}

func Test_RegistryMatch(t *testing.T) {
	r := registry.New([]registry.IOption{
		registry.OptId("testid"),
		registry.OptBind("127.0.0.1:7370"),
		registry.OptBindAdvertise("127.0.0.1:7370"),
		registry.OptRegistries(""),
		registry.OptAddr("127.0.0.1:9000"),
		registry.OptAdvertise("127.0.0.1:9000"),
	})
	go r.Serve()
	time.Sleep(sleepTime)

	serviceGroup := "testgroup"
	member1 := registry.NewMember(
		"testid1",
		"127.0.0.1:8370",
		"127.0.0.1:8370",
		"127.0.0.1:7370",
		serviceGroup,
		"127.0.0.1:80",
	)
	err := r.OnMemberJoin(member1)
	assert.Nil(t, err)

	member2 := registry.NewMember(
		"testid2",
		"127.0.0.1:8371",
		"127.0.0.1:8371",
		"127.0.0.1:7370",
		serviceGroup,
		"127.0.0.1:81",
	)
	err = r.OnMemberJoin(member2)
	assert.Nil(t, err)

	service, err := r.Match(serviceGroup, "werben")
	assert.Nil(t, err)
	assert.NotNil(t, service)
	assert.Equal(t, "testid1", service.Id)
	assert.Equal(t, serviceGroup, service.Group)
	assert.Equal(t, "127.0.0.1:80", service.Addr)

	service, err = r.Match(serviceGroup, "1testid2")
	assert.Nil(t, err)
	assert.NotNil(t, service)
	assert.Equal(t, "testid2", service.Id)
	assert.Equal(t, serviceGroup, service.Group)
	assert.Equal(t, "127.0.0.1:81", service.Addr)
	r.Close()
}

func Test_RegistryMembers(t *testing.T) {

	r := registry.New([]registry.IOption{
		registry.OptId("registy-id"),
		registry.OptBind("127.0.0.1:7370"),
		registry.OptBindAdvertise("127.0.0.1:7370"),
		registry.OptRegistries(""),
		registry.OptAddr("127.0.0.1:9000"),
		registry.OptAdvertise("127.0.0.1:9000"),
	})
	go r.Serve()
	time.Sleep(sleepTime)

	serviceGroup := "testgroup"

	member1 := registry.NewMember(
		"testid1",
		"127.0.0.1:8370",
		"127.0.0.1:8370",
		"127.0.0.1:7370",
		serviceGroup,
		"127.0.0.1:80",
	)
	err := r.OnMemberJoin(member1)
	assert.Nil(t, err)

	member2 := registry.NewMember(
		"testid2",
		"127.0.0.1:8371",
		"127.0.0.1:8371",
		"127.0.0.1:7371",
		serviceGroup,
		"127.0.0.1:81",
	)
	err = r.OnMemberJoin(member2)
	assert.Nil(t, err)

	services := r.Members(serviceGroup)
	sort.Slice(services, func(i int, j int) bool {
		return services[i].Id < services[j].Id
	})
	assert.NotNil(t, services)
	assert.Len(t, services, 2)

	assert.EqualValues(t, []*registry.Service{
		&member1.Service, &member2.Service,
	}, services)
	r.Close()
}
