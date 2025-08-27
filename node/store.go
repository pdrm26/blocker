package node

import (
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/pdrm26/blocker/proto"
	"github.com/pdrm26/blocker/types"
)

type TXHash = string
type TXStorer interface {
	Put(*proto.Transaction) error
	Get(TXHash) (*proto.Transaction, error)
}

type MemoryTXStore struct {
	lock sync.RWMutex
	txx  map[TXHash]*proto.Transaction
}

func NewMemoryTXStore() *MemoryTXStore {
	return &MemoryTXStore{
		txx: make(map[TXHash]*proto.Transaction),
	}
}

func (s *MemoryTXStore) Put(tx *proto.Transaction) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	hash := hex.EncodeToString(types.HashTransaction(tx))

	s.txx[hash] = tx
	return nil
}

func (s *MemoryTXStore) Get(txHash TXHash) (*proto.Transaction, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	tx, ok := s.txx[txHash]
	if !ok {
		return nil, fmt.Errorf("could not find a tx with txHash: %s", txHash)
	}

	return tx, nil
}

type UTXOStorer interface {
	Put(*UTXO) error
	Get(TXHash) (*UTXO, error)
}
type MemoryUTXOStore struct {
	lock   sync.RWMutex
	blocks map[string]*UTXO
}

func NewMemoryUTXOStore() *MemoryUTXOStore {
	return &MemoryUTXOStore{
		blocks: make(map[string]*UTXO),
	}
}

func (s *MemoryUTXOStore) Put(utxo *UTXO) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	key := fmt.Sprintf("%s_%d", utxo.Hash, utxo.OutIndex)
	s.blocks[key] = utxo
	return nil
}

func (s *MemoryUTXOStore) Get(hash BlockHash) (*UTXO, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	utxo, ok := s.blocks[hash]
	if !ok {
		return nil, fmt.Errorf("UTXO with hash [%s] does not exist", hash)
	}

	return utxo, nil
}

type BlockHash = string
type BlockStorer interface {
	Put(*proto.Block) error
	Get(BlockHash) (*proto.Block, error)
}

type MemoryBlockStore struct {
	lock   sync.RWMutex
	blocks map[string]*proto.Block
}

func NewMemoryBlockStore() *MemoryBlockStore {
	return &MemoryBlockStore{
		blocks: make(map[string]*proto.Block),
	}
}

func (s *MemoryBlockStore) Put(block *proto.Block) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	hash := hex.EncodeToString(types.HashBlock(block))
	s.blocks[hash] = block
	return nil
}

func (s *MemoryBlockStore) Get(hash BlockHash) (*proto.Block, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	block, ok := s.blocks[hash]
	if !ok {
		return nil, fmt.Errorf("block with hash [%s] does not exist", hash)
	}

	return block, nil
}
