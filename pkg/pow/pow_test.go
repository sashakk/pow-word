package pow

import (
	"bytes"
	"testing"
)

func TestGeneratePuzzle(t *testing.T) {
	tests := []struct {
		name       string
		puzzleSize int
		wantErr    bool
	}{
		{
			name:       "generate puzzle",
			puzzleSize: 20,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GeneratePuzzle(tt.puzzleSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("GeneratePuzzle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.puzzleSize {
				t.Errorf("size of puzzle doesn't match")
			}
		})
	}
}

func TestProofOfWork(t *testing.T) {
	tests := []struct {
		name       string
		puzzle     []byte
		hash       []byte
		nonceBytes []byte
		targetBits byte
	}{
		{
			name:       "dummy test that everything work as expected",
			puzzle:     []byte{'p', 'u', 'z', 'z', 'l', 'e'},
			hash:       []byte{0, 0, 9, 235, 79, 38, 254, 169, 208, 134, 193, 24, 69, 111, 200, 13, 186, 141, 240, 198, 112, 124, 210, 0, 123, 157, 62, 161, 129, 239, 161, 245},
			nonceBytes: []byte{54, 196, 1, 0},
			targetBits: 20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nonceBytes, hash := ProofOfWork(tt.puzzle, tt.targetBits)
			if !bytes.Equal(hash, tt.hash) {
				t.Error("getHash doesn't much")
				t.Fail()
			}
			if !bytes.Equal(nonceBytes, tt.nonceBytes) {
				t.Error("nonce doesn't much")
				t.Fail()
			}

			if !ValidateProofOfWork(tt.puzzle, nonceBytes, hash, int(tt.targetBits)) {
				t.Error("can't validate")
				t.Fail()
			}
		})
	}
}
