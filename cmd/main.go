package main;

import (
	"fmt";
	"time";
	"context";
	// "crypto";
	"crypto/rand";
	"crypto/rsa";

	sch "github.com/marekgalovic/simplechain/pkg";
)

func main() {
	var addrA sch.WalletAddress
	if _, err := rand.Read(addrA[:]); err != nil {
		panic(err)
	}
	pKeyA, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	walletA := sch.NewWallet(addrA, pKeyA)

	var addrB sch.WalletAddress
	if _, err := rand.Read(addrB[:]); err != nil {
		panic(err)
	}
	pKeyB, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	walletB := sch.NewWallet(addrB, pKeyB)

	tA, err := walletA.CreateTransaction(walletB.Address(), 10)
	if err != nil {
		panic(err)
	}

	tB, err := walletB.CreateTransaction(walletA.Address(), 5)
	if err != nil {
		panic(err)
	}

	var firstBlock [32]byte
	b := sch.NewBlock(firstBlock)
	b.AddTransaction(tA)
	b.AddTransaction(tB)

	m := sch.NewMiner(16)
	defer m.Stop()
	s := time.Now()
	m.Mine(context.Background(), b)
	fmt.Println(time.Since(s))
	
	fmt.Println(tA.Verify(walletA.PublicKey()))
	fmt.Println(tB.Verify(walletB.PublicKey()))
}