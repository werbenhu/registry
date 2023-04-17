package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/werbenhu/registry"
)

func Test_NewService(t *testing.T) {
	s := registry.NewService("test_id", "test_group", "172.0.0.1:80")
	assert.NotNil(t, s)
	assert.Equal(t, s.Id, "test_id")
	assert.Equal(t, s.Group, "test_group")
	assert.Equal(t, s.Addr, "172.0.0.1:8777")
}

func Test_NewMember(t *testing.T) {
	m := registry.NewMember("test_id", "172.0.0.1:7031", "172.0.0.2:7031", "172.0.0.1:7030", "test_group", "172.0.0.1:80")
	assert.NotNil(t, m)
	assert.Equal(t, m.Id, "test_id")
	assert.Equal(t, m.Service.Group, "test_group")
	assert.Equal(t, m.Bind, "172.0.0.1:7031")
	assert.Equal(t, m.Advertise, "172.0.0.2:7031")
	assert.Equal(t, m.Registries, "172.0.0.1:7030")
	assert.Equal(t, m.Service.Addr, "172.0.0.1:80")
	assert.Equal(t, m.Service.Id, "test_id")
}

func Test_NewSimpleMember(t *testing.T) {
	m := registry.NewSimpleMember("test_id", "172.0.0.1:7031", "172.0.0.2:7031")
	assert.NotNil(t, m)
	assert.Equal(t, m.Id, "test_id")
	assert.Equal(t, m.Service.Group, "")
	assert.Equal(t, m.Bind, "172.0.0.1:7031")
	assert.Equal(t, m.Advertise, "172.0.0.2:7031")
	assert.Equal(t, m.Registries, "")
	assert.Equal(t, m.Service.Addr, "")
	assert.Equal(t, m.Service.Id, "test_id")
}

func Test_MemberIsSelf(t *testing.T) {
	self := registry.NewSimpleMember("test_id", "172.0.0.1:7031", "172.0.0.2:7031")
	assert.NotNil(t, self)

	m := registry.NewMember("test_id", "172.0.0.1:7031", "172.0.0.2:7031", "172.0.0.1:7030", "test_group", "172.0.0.1:80")
	assert.NotNil(t, m)
	assert.Equal(t, self.IsSelf(m), true)

	m = registry.NewMember("test_id_xxx", "172.0.0.1:7031", "172.0.0.2:7031", "172.0.0.1:7030", "test_group", "172.0.0.1:80")
	assert.NotNil(t, m)
	assert.Equal(t, self.IsSelf(m), false)
}

func Test_MemberSetTag(t *testing.T) {
	m := registry.NewSimpleMember("test_id", "172.0.0.1:7031", "172.0.0.2:7031")
	assert.NotNil(t, m)

	m.SetTag(registry.TagAddr, "172.0.0.1:80")
	m.SetTag(registry.TagGroup, "test_group")
	m.SetTag(registry.TagReplicas, "10000")

	tags := m.GetTags()
	assert.Equal(t, map[string]string{
		registry.TagAddr:     "172.0.0.1:80",
		registry.TagGroup:    "test_group",
		registry.TagReplicas: "10000",
	}, tags)
}

func Test_MemberGetTag(t *testing.T) {
	m := registry.NewSimpleMember("test_id", "172.0.0.1:7031", "172.0.0.2:7031")
	assert.NotNil(t, m)

	m.SetTag(registry.TagAddr, "172.0.0.1:80")
	m.SetTag(registry.TagGroup, "test_group")
	m.SetTag(registry.TagReplicas, "10000")

	val, ok := m.GetTag(registry.TagAddr)
	assert.Equal(t, true, ok)
	assert.Equal(t, "172.0.0.1:80", val)

	val, ok = m.GetTag(registry.TagGroup)
	assert.Equal(t, true, ok)
	assert.Equal(t, "test_group", val)

	val, ok = m.GetTag(registry.TagReplicas)
	assert.Equal(t, true, ok)
	assert.Equal(t, "10000", val)

	val, ok = m.GetTag("wrong_tag")
	assert.Equal(t, false, ok)
	assert.Equal(t, "", val)
}

func Test_MemberSetTags(t *testing.T) {
	m := registry.NewSimpleMember("test_id", "172.0.0.1:7031", "172.0.0.2:7031")
	assert.NotNil(t, m)

	m.SetTags(map[string]string{
		registry.TagAddr:     "172.0.0.1:80",
		registry.TagGroup:    "test_group",
		registry.TagReplicas: "10000",
	})

	assert.Equal(t, m.Service.Addr, "172.0.0.1:80")
	assert.Equal(t, m.Service.Group, "test_group")

	tags := m.GetTags()
	assert.Equal(t, map[string]string{
		registry.TagAddr:     "172.0.0.1:80",
		registry.TagGroup:    "test_group",
		registry.TagReplicas: "10000",
	}, tags)
}

func Test_MemberGetTags(t *testing.T) {
	m := registry.NewMember("test_id", "172.0.0.1:7031", "172.0.0.2:7031", "172.0.0.1:7030", "test_group", "172.0.0.1:80")
	assert.NotNil(t, m)

	tags := m.GetTags()
	assert.Equal(t, map[string]string{
		registry.TagAddr:     "172.0.0.1:80",
		registry.TagGroup:    "test_group",
		registry.TagReplicas: "10000",
	}, tags)
}
