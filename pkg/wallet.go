package simplechain;

import (
	"crypto/rsa";
)

type WalletAddress [64]byte

type Wallet struct {
	address WalletAddress
	key *rsa.PrivateKey
}

func NewWallet(address WalletAddress, key *rsa.PrivateKey) *Wallet {
	return &Wallet {
		address: address,
		key: key,
	}
}

func (this *Wallet) Address() WalletAddress {
	return this.address
}

func (this *Wallet) PublicKey() *rsa.PublicKey {
	return this.key.Public().(*rsa.PublicKey)
}

func (this *Wallet) CreateTransaction(dst WalletAddress, amount uint64) (*Transaction, error) {
	transaction := &Transaction {
		src: this.Address(),
		dst: dst,
		amount: amount,
	}

	if err := transaction.Sign(this.key); err != nil {
		return nil, err
	}

	return transaction, nil
}