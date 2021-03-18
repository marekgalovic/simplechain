package main;

import (
	"fmt";
	"time";
	"context";
	// "crypto";
	"crypto/rand";
	"crypto/rsa";
	// "crypto/sha256";

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

	tA, err := walletA.Send(walletB.Address(), 10)
	if err != nil {
		panic(err)
	}

	tB, err := walletB.Send(walletA.Address(), 5)
	if err != nil {
		panic(err)
	}

	fmt.Println(tA.Verify(walletA.PublicKey()))
	fmt.Println(tB.Verify(walletB.PublicKey()))

	var firstBlock [32]byte
	b := sch.NewBlock(firstBlock)
	b.AddTransaction(tA)
	b.AddTransaction(tB)

	m := sch.NewMiner(16)
	defer m.Stop()
	s := time.Now()
	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()
	go func() {
		time.Sleep(20 * time.Second)
		ctxCancel()
	}()
	if err := m.MineBlock(ctx, b); err != nil {
		panic(err)
	}
	fmt.Println(time.Since(s))

	fmt.Println(b.Hash())
	fmt.Println(b.Nonce())
	
}