package mock_state

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/trie"
)

var emptyHash = common.Hash{}

type MockDatabase struct {
	state.Database
	usedTries map[*MockSecureTrie]struct{} // record opened tries
	ReadCode  map[common.Hash][]byte       // addrHash => contract code
}

func NewMockDatabase(db state.Database) *MockDatabase {
	return &MockDatabase{
		Database:  db,
		usedTries: make(map[*MockSecureTrie]struct{}, 0),
		ReadCode:  make(map[common.Hash][]byte, 0),
	}
}

func (m *MockDatabase) OpenTrie(root common.Hash) (state.Trie, error) {
	tr, err := m.Database.OpenTrie(root)
	if err != nil {
		return nil, err
	}
	mockTrie := NewMockSecureTrie(tr, emptyHash)
	m.usedTries[mockTrie] = struct{}{}
	return mockTrie, nil
}

func (m *MockDatabase) OpenStorageTrie(addrHash, root common.Hash) (state.Trie, error) {
	tr, err := m.Database.OpenStorageTrie(addrHash, root)
	if err != nil {
		return nil, err
	}
	mockTrie := NewMockSecureTrie(tr, addrHash)
	m.usedTries[mockTrie] = struct{}{}
	return mockTrie, nil
}

func (m *MockDatabase) CopyTrie(t state.Trie) state.Trie {
	switch t := t.(type) {
	case *MockSecureTrie:
		switch T:=t.Trie.(type) {
		case *trie.SecureTrie:
			return NewMockSecureTrie(T.Copy(),common.Hash{})
		default:
			panic(fmt.Errorf("unknow trie type %T",T))
		}
	default:
		panic(fmt.Errorf("unknown trie type %T", t))
	}
	panic(1)
}

func (m *MockDatabase) ContractCode(addrHash, codeHash common.Hash) ([]byte, error) {
	result, err := m.Database.ContractCode(addrHash, codeHash)
	if result != nil && len(result) > 0 {
		m.ReadCode[addrHash] = result
	}
	return result, err
}

func (m *MockDatabase) GetAllKey() (map[common.Address]struct{}, map[common.Hash][]string, map[common.Hash]bool) {
	result1 := make(map[common.Address]struct{}, 0)
	result2 := make(map[common.Hash][]string, 0)
	result3 := make(map[common.Hash]bool, 0)
	for t := range m.usedTries {
		for addr := range t.ReadAddr {
			result1[addr] = struct{}{}
		}
		if t.AddrHash != emptyHash {
			if result2[t.AddrHash] == nil {
				result2[t.AddrHash] = make([]string, 0)
			}
			for key := range t.ReadKey {
				result2[t.AddrHash] = append(result2[t.AddrHash], key)
			}
		}
		for key := range t.GetUsedNodeKey() {
			result3[key] = true
		}
	}
	return result1, result2, result3
}
