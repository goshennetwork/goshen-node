package mock_state

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
)

type MockSecureTrie struct {
	state.Trie
	RootHash       common.Hash
	ReadStorageKey map[string]struct{}
}

func NewMockSecureTrie(trie state.Trie, rootHash common.Hash) *MockSecureTrie {
	return &MockSecureTrie{Trie: trie, RootHash: rootHash, ReadStorageKey: make(map[string]struct{}, 0)}
}

func (m *MockSecureTrie) TryGet(key []byte) ([]byte, error) {
	m.ReadStorageKey[string(key)] = struct{}{}
	return m.Trie.TryGet(key)
}
