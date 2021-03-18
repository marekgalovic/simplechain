package simplechain;

import (
	"hash";
	"crypto/sha256";
	"errors";
)

var (
	ErrBlockIsFull error = errors.New("Block is full")
)

type Block struct {
	prevBlock [sha256.Size]byte
	transactions []*Transaction
	nonce [8]byte
}

func NewBlock(prevBlock [sha256.Size]byte) *Block {
	return &Block {
		prevBlock: prevBlock,
		transactions: make([]*Transaction, 0, BLOCK_SIZE),
	}
}

func (this *Block) AddTransaction(transaction *Transaction) error {
	if len(this.transactions) >= BLOCK_SIZE {
		return ErrBlockIsFull
	}

	this.transactions = append(this.transactions, transaction)
	return nil
}

func (this *Block) writeHash(hash hash.Hash) error {
	if _, err := hash.Write(this.prevBlock[:]); err != nil {
		return err
	}

	for _, t := range this.transactions {
		if err := t.writeHash(hash); err != nil {
			return err
		}
	}

	return nil
}