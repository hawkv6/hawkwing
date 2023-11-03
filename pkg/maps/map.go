package maps

import (
	"fmt"
	"sync"

	"github.com/cilium/ebpf"
	"github.com/hawkv6/hawkwing/pkg/bpf"
)

type Map struct {
	spec *ebpf.MapSpec
	m    *ebpf.Map
	lock sync.RWMutex
	path string
}

// OpenOrCreate opens or creates the map.
//
// Returns:
//   - An error if the map could not be opened or created.
func (m *Map) OpenOrCreate() error {
	m.lock.Lock()
	defer m.lock.Unlock()

	return m.openOrCreate(true)
}

func (m *Map) openOrCreate(pin bool) error {
	if m.m != nil {
		return nil
	}

	if m.spec == nil {
		return fmt.Errorf("cannot create map without spec")
	}

	m.setPath()

	if pin {
		m.spec.Pinning = ebpf.PinByName
	}

	em, err := bpf.CreateMap(m.spec, m.path)
	if err != nil {
		return err
	}

	m.m = em

	return nil
}

// Open opens the map.
//
// Returns:
//   - An error if the map could not be opened.
func (m *Map) Open() error {
	m.lock.Lock()
	defer m.lock.Unlock()

	return m.open()
}

func (m *Map) open() error {
	if m.m != nil {
		return nil
	}

	em, err := ebpf.LoadPinnedMap(m.path, nil)
	if err != nil {
		return fmt.Errorf("failed to load pinned map: %w", err)
	}

	m.m = em

	return nil
}

func (m *Map) setPath() {
	if m.path == "" {
		m.path = bpf.BpffsRoot
	}
}

// Lookup looks up the value corresponding to the given key in the map.
//
// Parameters:
//   - key: The key to lookup.
//   - value: The value to store the result in.
//
// Returns:
//   - An error if the map could not be opened or the lookup failed.
func (m *Map) Lookup(key interface{}, value interface{}) error {
	if err := m.open(); err != nil {
		return err
	}

	m.lock.RLock()
	defer m.lock.RUnlock()

	err := m.m.Lookup(key, value)
	if err != nil {
		return fmt.Errorf("failed to lookup map %s: %w", m.spec.Name, err)
	}

	return nil
}

// UpdateInner updates the inner map corresponding to the given outer key with
// the given inner key and value.
//
// Parameters:
//   - outerKey: The key of the outer map.
//   - innerKey: The key of the inner map.
//   - innerValue: The value to update the inner map with.
//
// Returns:
//   - An error if the inner map could not be updated.
func (m *Map) UpdateInner(outerKey interface{}, innerKey interface{}, innerValue interface{}) error {
	var err error

	m.lock.Lock()
	defer m.lock.Unlock()

	var innerMapID ebpf.MapID
	err = m.m.Lookup(outerKey, &innerMapID)
	if err != nil {
		return fmt.Errorf("could not find inner map in outer map: %s", err)
	}

	innerMap, err := ebpf.NewMapFromID(innerMapID)
	if err != nil {
		return fmt.Errorf("could not create inner map from ID: %s", err)
	}

	err = innerMap.Put(innerKey, innerValue)
	if err != nil {
		return fmt.Errorf("could not update inner map: %s", err)
	}

	return nil
}
