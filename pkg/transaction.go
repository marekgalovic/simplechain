package simplechain;

import (
	"io";
	"time";
	"hash";
	"crypto";
	"crypto/rand";
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
	timestamp time.Time
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

	hash := sha256HashPool.Get().(hash.Hash)
	defer sha256HashPool.Put(hash)
	hash.Reset()
	err := this.write(hash, false)
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
	hash := sha256HashPool.Get().(hash.Hash)
	defer sha256HashPool.Put(hash)
	hash.Reset()
	if err := this.write(hash, false); err != nil {
		return err
	}

	return rsa.VerifyPSS(key, crypto.SHA256, hash.Sum(nil), this.signature, nil)
}

func (this *Transaction) write(w io.Writer, includeSignature bool) error {
	if _, err := w.Write(this.src[:]); err != nil {
		return err
	}
	if _, err := w.Write(this.dst[:]); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, this.amount); err != nil {
		return err
	}
	tsBytes, err := this.timestamp.MarshalBinary()
	if err != nil {
		return err
	}
	if _, err := w.Write(tsBytes); err != nil {
		return err
	}

	if includeSignature {
		if _, err := w.Write(this.signature); err != nil {
			return err
		}
	}

	return nil
}