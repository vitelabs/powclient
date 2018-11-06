package cpu

import (
	"errors"
	"github.com/vitelabs/go-vite/common/helper"
	"github.com/vitelabs/go-vite/common/types"
	"github.com/vitelabs/go-vite/crypto"
	"github.com/vitelabs/go-vite/pow"
	"golang.org/x/crypto/blake2b"
	"math/big"
)

func GetPowNonce(threshold *big.Int, dataHash types.Hash) ([]byte, error) {

	if threshold == nil {
		return nil, errors.New("threshold can't be nil")
	}
	if threshold.BitLen() > 256 {
		return nil, errors.New("threshold too long")
	}

	data := dataHash.Bytes()
	target256 := helper.LeftPadBytes(threshold.Bytes(), 32)
	for {
		nonce := crypto.GetEntropyCSPRNG(8)
		out := powHash256(nonce, data)
		if pow.QuickGreater(out, target256) {
			return nonce, nil
		}
	}
	return nil, errors.New("get pow nonce error")
}

func powHash256(nonce []byte, data []byte) []byte {
	hash, _ := blake2b.New256(nil)
	hash.Write(nonce)
	hash.Write(data)
	out := hash.Sum(nil)
	return out
}
