package pow

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"math"
	"math/big"
	"pow/pkg/log"
)

func ProofOfWork(puzzle []byte, targetBits byte) ([]byte, []byte) {
	if puzzle == nil {
		return nil, nil
	}

	target := big.NewInt(1)
	target.Lsh(target, uint(256-int(targetBits)))

	var nonce uint32
	for nonce < math.MaxUint32 {
		nonceBytes, err := getNonceBytes(nonce)
		if err != nil {
			return nil, nil
		}

		hash := getHash(puzzle, nonceBytes)

		var hashInt big.Int
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(target) == -1 {
			return nonceBytes, hash[:]
		} else {
			nonce++
		}
	}

	return nil, nil
}

func getHash(data, nonceBytes []byte) [32]byte {
	return sha256.Sum256(bytes.Join([][]byte{data, nonceBytes}, []byte{}))
}

func getNonceBytes(nonce uint32) ([]byte, error) {
	bufNonce := new(bytes.Buffer)
	if err := binary.Write(bufNonce, binary.LittleEndian, nonce); err != nil {
		log.Errorf("error convert nonce to bytes")
		return nil, err
	}

	return bufNonce.Bytes(), nil
}

func GeneratePuzzle(puzzleSize int) ([]byte, error) {
	ret := make([]byte, puzzleSize)
	if _, err := rand.Read(ret); err != nil {
		return nil, err
	}

	return ret, nil
}

func ValidateProofOfWork(puzzle, nonceBytes, hash []byte, targetBits int) bool {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	hashBytes := getHash(puzzle, nonceBytes)

	return bytes.Equal(hashBytes[:], hash)
}
