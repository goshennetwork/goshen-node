package mock_state

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
)

type MockSecureTrie struct {
	state.Trie
	AddrHash common.Hash
	ReadAddr map[common.Address]struct{}
	ReadKey  map[string]struct{}
}

func NewMockSecureTrie(trie state.Trie, addrHash common.Hash) *MockSecureTrie {
	result := &MockSecureTrie{Trie: trie,
		AddrHash: addrHash,
		ReadAddr: make(map[common.Address]struct{}, 0),
		ReadKey:  make(map[string]struct{}, 0)}
	return result
}

func (m *MockSecureTrie) TryGet(key []byte) ([]byte, error) {
	if len(key) == common.AddressLength {
		addr := common.BytesToAddress(key)
		if _, ok := m.ReadAddr[addr]; !ok {
			m.ReadAddr[addr] = struct{}{}
		}
	} else {
		m.ReadKey[string(key)] = struct{}{}
	}
	return m.Trie.TryGet(key)
}
