package simplechain;

import (
	"time";
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

func (this *Wallet) Send(dst WalletAddress, amount uint64) (*Transaction, error) {
	transaction := &Transaction {
		src: this.Address(),
		dst: dst,
		amount: amount,
		timestamp: time.Now(),
	}

	if err := transaction.Sign(this.key); err != nil {
		return nil, err
	}

	return transaction, nil
}