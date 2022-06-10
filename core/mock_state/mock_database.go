package mock_state

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
)

type MockDatabase struct {
	state.Database
	usedTries map[*MockSecureTrie]struct{} // record opened tries
	readCode  map[common.Hash][]byte       // addrHash => contract code
}

func NewMockDatabase(db state.Database) *MockDatabase {
	return &MockDatabase{
		Database:  db,
		usedTries: make(map[*MockSecureTrie]struct{}, 0),
		readCode:  make(map[common.Hash][]byte, 0),
	}
}

func (m *MockDatabase) OpenTrie(root common.Hash) (state.Trie, error) {
	tr, err := m.Database.OpenTrie(root)
	if err != nil {
		return nil, err
	}
	mockTrie := NewMockSecureTrie(tr, root)
	m.usedTries[mockTrie] = struct{}{}
	return mockTrie, nil
}

func (m *MockDatabase) OpenStorageTrie(addrHash, root common.Hash) (state.Trie, error) {
	tr, err := m.Database.OpenStorageTrie(addrHash, root)
	if err != nil {
		return nil, err
	}
	mockTrie := NewMockSecureTrie(tr, root)
	m.usedTries[mockTrie] = struct{}{}
	return mockTrie, nil
}

func (m *MockDatabase) CopyTrie(t state.Trie) state.Trie {
	switch t := t.(type) {
	case *MockSecureTrie:
		cpy := *t
		return &cpy
	default:
		panic(fmt.Errorf("unknown trie type %T", t))
	}
}

func (m *MockDatabase) ContractCode(addrHash, codeHash common.Hash) ([]byte, error) {
	result, err := m.Database.ContractCode(addrHash, codeHash)
	if result != nil && len(result) > 0 {
		m.readCode[addrHash] = result
	}
	return result, err
}

type SimpleHashSet struct {
	nodes map[string][]byte
}

// Put stores a new node in the set
func (db *SimpleHashSet) Put(key []byte, value []byte) error {
	if _, ok := db.nodes[string(key)]; ok {
		return nil
	}
	keystr := string(key)
	db.nodes[keystr] = common.CopyBytes(value)

	return nil
}

// Delete removes a node from the set
func (db *SimpleHashSet) Delete(key []byte) error {
	return nil
}

func (m *MockDatabase) GetReadStorageKeyProof() ([][]byte, error) {
	nodeSet := &SimpleHashSet{nodes: make(map[string][]byte, 0)}
	for t := range m.usedTries {
		for keyStr := range t.ReadStorageKey {
			err := t.Prove([]byte(keyStr), 0, nodeSet)
			if err != nil {
				return nil, err
			}
		}
	}
	result := make([][]byte, 0)
	for _, raw := range nodeSet.nodes {
		result = append(result, raw)
	}
	for _, code := range m.readCode {
		result = append(result, code)
	}
	return result, nil
}
