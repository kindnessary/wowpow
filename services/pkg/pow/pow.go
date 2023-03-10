package pow

import (
	"bytes"
	"crypto/sha256"
	"log"
	"math"
	"math/big"
)

type Processor struct {
	initialData []byte
	difficulty  byte
	target      *big.Int
}

func New(difficulty byte, initialData []byte) *Processor {
	target := big.NewInt(1)
	target.Lsh(target, 256-uint(difficulty))

	return &Processor{
		initialData: initialData,
		difficulty:  difficulty,
		target:      target,
	}
}

func (p *Processor) initNonce(nonce uint64) []byte {
	data := bytes.Join(
		[][]byte{
			p.initialData,
			ToHex(nonce),
			ToHex(uint64(p.difficulty)),
		},
		[]byte{},
	)
	return data
}

func (p *Processor) Calculate() uint64 {
	var intHash big.Int
	var hash [32]byte

	var result uint64
	for i := uint64(0); i < math.MaxInt64; i++ {
		data := p.initNonce(i)
		hash = sha256.Sum256(data)

		intHash.SetBytes(hash[:])

		if intHash.Cmp(p.target) == -1 {
			result = i
			break
		}
	}

	log.Printf("calculated hash: \r%x\n", hash)

	return result
}

func (p *Processor) Validate(nonce uint64) bool {
	var intHash big.Int

	data := p.initNonce(nonce)

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(p.target) == -1
}
