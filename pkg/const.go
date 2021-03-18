package simplechain;

import (
	"sync";
	"bytes";
	"crypto/sha256";
)

const BLOCK_SIZE int = 100
const NONCE_SIZE int = 4
const LEADING_ZERO_BYTES int = 3

var sha256HashPool = sync.Pool {
	New: func() interface{} {
		return sha256.New()
	},
}

var bufferPool = sync.Pool {
	New: func() interface{} {
		return bytes.NewBuffer(nil)
	},
}