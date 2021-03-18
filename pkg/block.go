package simplechain;

import (
	"io";
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
	nonce [NONCE_SIZE]byte
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

func (this *Block) Nonce() [NONCE_SIZE]byte {
	return this.nonce
}

func (this *Block) Hash() ([]byte, error) {
	hash := sha256HashPool.Get().(hash.Hash)
	defer sha256HashPool.Put(hash)
	hash.Reset()
	if err := this.write(hash); err != nil {
		return nil, err
	}
	return hash.Sum(nil), nil
}

func (this *Block) setNonce(nonce []byte) {
	for i, v := range nonce {
		this.nonce[i] = v
	}	
}

func (this *Block) write(w io.Writer) error {
	if _, err := w.Write(this.prevBlock[:]); err != nil {
		return err
	}

	for _, t := range this.transactions {
		if err := t.write(w, true); err != nil {
			return err
		}
	}

	if _, err := w.Write(this.nonce[:]); err != nil {
		return err
	}

	return nil
}