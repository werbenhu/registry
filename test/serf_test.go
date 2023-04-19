package test

import (
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/werbenhu/registry"
)

func Test_NewSerf(t *testing.T) {
	member := registry.NewMember(
		"test_id",
		"172.0.0.1:7730",
		"172.0.0.1:7730",
		"",
		"test_group",
		"172.0.0.1:80",
	)

	serf := registry.NewSerf(member)
	assert.NotNil(t, serf)
	serf.Stop()
}

func Test_SerfLocalMember(t *testing.T) {
	member := registry.NewMember(
		"test_id",
		"127.0.0.1:7730",
		"127.0.0.1:7730",
		"",
		"test_group",
		"127.0.0.1:80",
	)

	serf := registry.NewSerf(member)
	assert.NotNil(t, serf)
	err := serf.Start()
	assert.Nil(t, err)

	local := serf.LocalMember()
	assert.NotNil(t, local)
	assert.EqualValues(t, member, local)

	serf.Stop()
}

func Test_SerfMembers(t *testing.T) {

	member1 := registry.NewMember(
		"test_id1",
		"127.0.0.1:7730",
		"127.0.0.1:7730",
		"",
		"test_group",
		"127.0.0.1:80",
	)

	serf1 := registry.NewSerf(member1)
	assert.NotNil(t, serf1)
	err := serf1.Start()
	assert.Nil(t, err)
	time.Sleep(sleepTime)

	member2 := registry.NewMember(
		"test_id2",
		"127.0.0.1:7731",
		"127.0.0.1:7731",
		"127.0.0.1:7730",
		"test_group",
		"127.0.0.1:81",
	)
	serf2 := registry.NewSerf(member2)
	assert.NotNil(t, serf2)
	err = serf2.Start()
	assert.Nil(t, err)
	time.Sleep(sleepTime)

	ms1 := serf1.Members()
	ms2 := serf2.Members()
	assert.Len(t, ms1, 2)
	assert.Len(t, ms2, 2)

	sort.Slice(ms1, func(i, j int) bool { return ms1[i].Id < ms1[j].Id })
	sort.Slice(ms2, func(i, j int) bool { return ms2[i].Id < ms2[j].Id })
	assert.Equal(t, true, assert.ObjectsAreEqual(ms1, ms2))

	member2.Registries = ""
	expert := []*registry.Member{
		member1, member2,
	}
	assert.Equal(t, true, assert.ObjectsAreEqual(expert, ms1))

	serf1.Stop()
	serf2.Stop()
}
