package cpu

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/vitelabs/go-vite/common/helper"
	"github.com/vitelabs/go-vite/common/types"
	"github.com/vitelabs/go-vite/crypto"
	"github.com/vitelabs/go-vite/pow"
	"github.com/vitelabs/powclient/log15"
	"golang.org/x/crypto/blake2b"
	"math/big"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var (
	MAXTHREADNUM = runtime.NumCPU()
	log          = log15.New("module", "service_cpu")
)

func GetPowNonce(threshold *big.Int, dataHash types.Hash) (*string, error) {
	log.Info("start time.no", nil, time.Now())
	if threshold == nil {
		return nil, errors.New("threshold can't be nil")
	}
	if threshold.BitLen() > 256 {
		return nil, errors.New("threshold too long")
	}

	data := dataHash.Bytes()
	target256 := helper.LeftPadBytes(threshold.Bytes(), 32)

	var (
		pend       sync.WaitGroup
		resultChan = make(chan string)
		abort      = make(chan struct{})
		stop       = make(chan struct{})
		result     = ""
	)

	threadnum := MAXTHREADNUM / 2

	if threadnum <= 0 {
		threadnum = 1
	}
	for i := 0; i < threadnum; i++ {
		pend.Add(1)
		go func(index int) {
			defer pend.Done()
		MINE:
			for {
				select {
				case <-abort:
					log.Info("abort", "index", index)
					break MINE
				default:
					nonce := crypto.GetEntropyCSPRNG(8)
					out := powHash256(nonce, data)
					if pow.QuickGreater(out, target256) {
						hexNonce := strconv.FormatUint(binary.LittleEndian.Uint64(nonce), 16)
						log.Info("success", nil, fmt.Sprint("hexNonce:%v index:%v\n", hexNonce, index))
						resultChan <- hexNonce
						break MINE
					}
				}
			}
		}(i)
	}

	timer := time.AfterFunc(60*time.Second, func() {
		log.Info("timeout", nil, time.Now())
		close(stop)
	})

	select {
	case <-stop:
		close(abort)
	case result = <-resultChan:
		log.Info("success select", nil, result)
		close(abort)
	}
	pend.Wait()
	timer.Stop()
	log.Info("end time.no", nil, time.Now())
	if result == "" {
		return nil, errors.New("get pow nonce error")
	}

	return &result, nil
}

func powHash256(nonce []byte, data []byte) []byte {
	hash, _ := blake2b.New256(nil)
	hash.Write(nonce)
	hash.Write(data)
	out := hash.Sum(nil)
	return out
}
