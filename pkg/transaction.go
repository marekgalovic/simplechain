package simplechain;

import (
	"hash";
	"crypto";
	"crypto/rand";
	"crypto/sha256";
	"crypto/rsa";
	"encoding/binary";
	"errors";
)

var (
	ErrTransactionAlreadySigned error = errors.New("Transaction is already signed")
)

type Transaction struct {
	src WalletAddress
	dst WalletAddress
	amount uint64
	signature []byte
}

func (this *Transaction) Src() WalletAddress {
	return this.src
}

func (this *Transaction) Dst() WalletAddress {
	return this.dst
}

func (this *Transaction) Amount() uint64 {
	return this.amount
}

func (this *Transaction) Signature() []byte {
	return this.signature
}

func (this *Transaction) Sign(key *rsa.PrivateKey) error {
	if len(this.signature) > 0 {
		return ErrTransactionAlreadySigned
	}

	hash := sha256.New()
	err := this.writeHash(hash)
	if err != nil {
		return err
	}

	this.signature, err = rsa.SignPSS(rand.Reader, key, crypto.SHA256, hash.Sum(nil), nil)
	if err != nil {
		return err
	}

	return nil
}

func (this *Transaction) Verify(key *rsa.PublicKey) error {
	hash := sha256.New()
	if err := this.writeHash(hash); err != nil {
		return err
	}

	return rsa.VerifyPSS(key, crypto.SHA256, hash.Sum(nil), this.signature, nil)
}

func (this *Transaction) writeHash(hash hash.Hash) error {
	if _, err := hash.Write(this.src[:]); err != nil {
		return err
	}
	if _, err := hash.Write(this.dst[:]); err != nil {
		return err
	}
	if err := binary.Write(hash, binary.LittleEndian, this.amount); err != nil {
		return err
	}

	return nil
}