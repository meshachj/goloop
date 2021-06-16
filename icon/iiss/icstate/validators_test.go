package icstate

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/icon-project/goloop/icon/iiss/icutils"
	"github.com/icon-project/goloop/module"
)

type ownerToNodeMapper struct {
	o2n map[string]module.Address
	n2o map[string]module.Address
}

func (m *ownerToNodeMapper) add(owner module.Address, node module.Address) {
	if owner == nil || node == nil {
		return
	}
	if m.o2n == nil {
		m.o2n = make(map[string]module.Address)
	}
	if m.n2o == nil {
		m.n2o = make(map[string]module.Address)
	}
	m.o2n[icutils.ToKey(owner)] = node
	m.o2n[icutils.ToKey(node)] = owner
}

func (m *ownerToNodeMapper) GetNodeByOwner(owner module.Address) module.Address {
	node, ok := m.o2n[icutils.ToKey(owner)]
	if !ok {
		return nil
	}
	return node
}

func (m *ownerToNodeMapper) GetOwnerByNode(node module.Address) module.Address {
	owner, ok := m.n2o[icutils.ToKey(node)]
	if !ok {
		return nil
	}
	return owner
}

func newDummyOwnerToNodeMapper(size int) OwnerToNodeMappable {
	m := new(ownerToNodeMapper)
	for i := 0; i < size; i++ {
		owner := newDummyAddress(i)
		node := newDummyAddress(i + 100)
		m.add(owner, node)
	}
	return m
}

func newDummyValidatorsData(size int) *validatorsData {
	snapshots := newDummyPRepSnapshots(size)
	m := newDummyOwnerToNodeMapper(size)
	vd := new(validatorsData)
	vd.init(snapshots, m, size)
	return vd
}

func TestValidatorsData_init(t *testing.T) {
	size := 10
	snapshots := newDummyPRepSnapshots(size)
	m := newDummyOwnerToNodeMapper(size)

	type a struct {
		m int
		n string
	}

	b := a{1, "hello"}
	c := a{2, "w"}
	b.m = c.m
	b.n = c.n

	vd := validatorsData{}
	vd.init(snapshots, m, size)
	assert.Equal(t, size, vd.Len())

	for i := 0; i < size; i++ {
		snapshot := snapshots.Get(i).(*PRepSnapshot)
		node := m.GetNodeByOwner(snapshot.Owner())

		node2 := vd.Get(i)
		assert.True(t, node.Equal(node2))
		assert.Equal(t, i, vd.IndexOf(node))
	}

	assert.Equal(t, size, vd.NextPRepSnapshotIndex())
}

func TestValidatorsData_clone(t *testing.T) {
	size := 22
	vd := newDummyValidatorsData(size)
	vd2 := vd.clone()
	assert.True(t, vd.equal(&vd2))
	vss := ValidatorsSnapshot{}
	vss.validatorsData = vd2

	hash := vd.Hash()
	assert.Zero(t, bytes.Compare(vd.Hash(), vd2.Hash()))
	assert.Equal(t, 32, len(hash))
	assert.True(t, vd2.equal(&vss.validatorsData))
}

func TestValidatorsSnapshot_RLPEncodeDecode(t *testing.T) {
	state := newDummyState(false)

	size := 10
	nodes := make([]module.Address, size)
	for i := 0; i < size; i++ {
		nodes[i] = newDummyAddress(i)
	}

	vd := newValidatorsData(nodes)
	vss := &ValidatorsSnapshot{
		validatorsData: vd,
	}
	assert.Equal(t, size, vss.Len())

	err := state.SetValidatorsSnapshot(vss)
	assert.NoError(t, err)

	state = flushAndNewState(state, false)

	vss2 := state.GetValidatorsSnapshot()
	assert.NotNil(t, vss2)

	assert.True(t, vss.Equal(vss2))
	assert.Equal(t, size, vss2.Len())
}

func TestNewValidatorStateWithSnapshot(t *testing.T) {
	var snapshot *ValidatorsSnapshot
	vs := NewValidatorsStateWithSnapshot(snapshot)
	assert.Zero(t, vs.Len())
}