package server

import (
	"bytes"
	"pow/pkg/config"
	"pow/pkg/tcp"
	"testing"
)

func TestServer_prepareResponseValidatePoW(t *testing.T) {
	type fields struct {
		config *config.Config
	}
	type args struct {
		requestData []byte
		puzzle      []byte
	}
	type wantData struct {
		code byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   wantData
	}{
		{
			name: "valid request",
			args: args{
				puzzle:      []byte{'p', 'u', 'z'},
				requestData: []byte{3, 0, 0, 0, 0, 186, 116, 7, 185, 4, 248, 254, 124, 114, 125, 184, 111, 210, 145, 243, 151, 198, 76, 184, 143, 38, 58, 79, 16, 71, 95, 5, 246, 45, 145, 175, 176},
			},
			fields: fields{
				config: &config.Config{
					PuzzleSize: 3,
				},
			},
			want: wantData{
				code: tcp.CodeRequestValidPoW,
			},
		},

		{
			name: "invalid request (validation error)",
			args: args{
				puzzle:      []byte{'p', 'u', 'z'},
				requestData: []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			},
			fields: fields{
				config: &config.Config{
					PuzzleSize: 3,
				},
			},
			want: wantData{
				code: tcp.CodeRequestInvalidPoW,
			},
		},
		{
			name: "invalid request(short request)",
			args: args{
				puzzle:      []byte{'p', 'u', 'z'},
				requestData: []byte{1},
			},
			fields: fields{
				config: &config.Config{
					PuzzleSize: 3,
				},
			},
			want: wantData{
				code: tcp.CodeRequestInvalidPoW,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				config: tt.fields.config,
			}
			response := s.prepareResponseValidatePoW(tt.args.requestData, tt.args.puzzle)
			if response[0] != tt.want.code {
				t.Errorf("code doesn't match")
			}
		})
	}
}

func TestServer_prepareRequestGeneratePuzzle(t *testing.T) {
	type fields struct {
		config *config.Config
	}
	tests := []struct {
		name          string
		fields        fields
		wantPuzzleLen int
		wantResponse  []byte
		wantErr       bool
	}{
		{
			name: "valid response",
			fields: fields{
				config: &config.Config{
					PuzzleSize: 10,
					TargetBits: 10,
				},
			},
			wantPuzzleLen: 10,
			wantResponse:  []byte{tcp.CodeRequestReturnPuzzle, 10, 10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				config: tt.fields.config,
			}
			puzzle, response, err := s.prepareRequestGeneratePuzzle()
			if (err != nil) != tt.wantErr {
				t.Errorf("prepareRequestGeneratePuzzle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.wantResponse = append(tt.wantResponse, puzzle...)
			if len(puzzle) != tt.wantPuzzleLen {
				t.Errorf("puzzle len error")
			}
			if !bytes.Equal(tt.wantResponse, response) {
				t.Errorf("prepareRequestGeneratePuzzle() wantResponse = %v, want %v", response, tt.wantResponse)
			}
		})
	}
}
