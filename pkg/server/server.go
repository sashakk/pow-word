package server

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"pow/pkg/config"
	"pow/pkg/log"
	"pow/pkg/pow"
	"pow/pkg/quote"
	"pow/pkg/tcp"
	"time"
)

const (
	BufSize = 512
)

type Server struct {
	listener net.Listener
	config   *config.Config
}

func NewServer(config *config.Config) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
	if err != nil {
		return nil, err
	}
	return &Server{
		listener: listener,
		config:   config,
	}, nil
}

func (s *Server) Close() {
	if err := s.listener.Close(); err != nil {
		log.Errorf("error while closing connection: %s", err)
	}
}

func (s *Server) Serve() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Errorf("error while accepting connection: %s", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) writeResponse(conn net.Conn, data []byte) error {
	if _, err := conn.Write(data); err != nil {
		return fmt.Errorf("error while write wantResponse: %w", err)
	}
	return nil
}

func (s *Server) handleConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Errorf("error while closing connection: %s", err)
		}
	}()

	if err := conn.SetDeadline(time.Now().Add(time.Duration(s.config.ConnTimeout) * time.Millisecond)); err != nil {
		log.Errorf("error setting deadline: %s", err)
		return
	}

	reader := bufio.NewReader(conn)
	requestData := make([]byte, BufSize)
	var puzzle []byte

	for {
		if _, err := reader.Read(requestData); err == io.EOF {
			log.Infof("connection interrupt")
			return
		} else if err != nil {
			log.Infof("error while reading: %s", err)
			return
		}

		code := requestData[tcp.IndexRequestCode]
		var err error

		switch code {
		case tcp.CodeRequestPuzzle:
			var response []byte
			puzzle, response, err = s.prepareRequestGeneratePuzzle()
			if err = s.writeResponse(conn, response); err != nil {
				log.Errorf("error while write RequestPuzzle: %s", err)
				return
			}
		case tcp.CodeRequestValidatePoW:
			response := s.prepareResponseValidatePoW(requestData, puzzle)
			if err = s.writeResponse(conn, response); err != nil {
				log.Errorf("error while write RequestValidatePoW: %s", err)
				return
			}
		default:
			log.Errorf("invalid request code")
			return
		}
	}
}

func (s *Server) prepareResponseValidatePoW(requestData, puzzle []byte) []byte {
	if len(requestData) < tcp.EndNonce+tcp.SizeOfHash {
		return []byte{tcp.CodeRequestInvalidPoW}
	}

	nonceBytes := requestData[tcp.StartNonce:tcp.EndNonce]
	hash := requestData[tcp.StartHash : tcp.StartHash+tcp.SizeOfHash]
	log.Infof("hash: %v, nonceBytes: %v", hash, nonceBytes)
	if pow.ValidateProofOfWork(puzzle, nonceBytes, hash, s.config.TargetBits) {
		q := quote.GetQuote()
		return bytes.Join([][]byte{{tcp.CodeRequestValidPoW, byte(len(q))}, []byte(q)}, []byte{})
	}
	return []byte{tcp.CodeRequestInvalidPoW}
}

func (s *Server) prepareRequestGeneratePuzzle() ([]byte, []byte, error) {
	puzzle, err := pow.GeneratePuzzle(s.config.PuzzleSize)
	log.Infof("Puzzle: %v", puzzle)
	if err != nil {
		return nil, nil, fmt.Errorf("error while generatePuzzle: %w", err)
	}
	response := bytes.Join([][]byte{
		{tcp.CodeRequestReturnPuzzle, byte(s.config.PuzzleSize), byte(s.config.TargetBits)},
		puzzle,
	},
		[]byte{},
	)
	return puzzle, response, nil
}
