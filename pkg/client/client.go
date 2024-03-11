package client

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"pow/pkg/log"
	"pow/pkg/pow"
	"pow/pkg/tcp"
	"time"
)

const (
	BufSize = 512
)

type Client struct {
	conn net.Conn
}

func (c *Client) Close() {
	defer func() {
		if err := c.conn.Close(); err != nil {
			log.Errorf("error closing connection: %s", err)
			return
		}
	}()

}

func NewClient(host string, port int, connTimeout int) (*Client, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, fmt.Errorf("error connecting to server: %s", err)
	}

	err = conn.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(connTimeout)))
	if err != nil {
		return nil, fmt.Errorf("error connecting to server: %w", err)
	}

	return &Client{
		conn: conn,
	}, nil
}

// Authorize establishes authorization with the server
// Returns quote and error
// If any error occurs during the process, it returns an empty string and an error
func (c *Client) Authorize() (string, error) {
	var err error

	if _, err = c.conn.Write([]byte{tcp.CodeRequestPuzzle}); err != nil {
		return "", fmt.Errorf("error write request to get a puzzle: %w", err)
	}

	reader := bufio.NewReader(c.conn)
	requestData := make([]byte, BufSize)
	for {
		if _, err = reader.Read(requestData); err != nil {
			return "", fmt.Errorf("error reading error: %s", err)
		}

		code := requestData[tcp.IndexRequestCode]

		switch code {
		case tcp.CodeRequestReturnPuzzle:
			err = c.handlePuzzleResponse(requestData, err)
			if err != nil {
				return "", fmt.Errorf("error while handle puzzle response: %w", err)
			}
		case tcp.CodeRequestValidPoW:
			quoteByte := c.parseValidPoWResponse(requestData)
			return string(quoteByte), nil
		case tcp.CodeRequestInvalidPoW:
			return "", fmt.Errorf("failed to validate puzzle")
		default:
			log.Errorf("invalid code")
			return "", fmt.Errorf("invalid code from sever")
		}
	}
}

func (c *Client) handlePuzzleResponse(requestData []byte, err error) error {
	targetBits, puzzle := c.parsePuzzleResponse(requestData)
	log.Infof("puzzle: %v", puzzle)
	nonceBytes, hash := pow.ProofOfWork(puzzle, targetBits)
	if len(hash) == 0 {
		return fmt.Errorf("error during ProofOfWork puzzle: %w", err)
	}
	log.Infof("hash: %v, nonceBytes: %v", hash, nonceBytes)
	if _, err = c.conn.Write(bytes.Join([][]byte{{tcp.CodeRequestValidatePoW}, nonceBytes, hash}, []byte{})); err != nil {
		return fmt.Errorf("error writing CodeRequestValidatePoW: %s", err)
	}
	return nil
}

func (c *Client) parseValidPoWResponse(requestData []byte) []byte {
	sizeOfQuote := requestData[tcp.IndexSizeOfQuote]
	quoteByte := requestData[tcp.StartQuote : tcp.StartQuote+sizeOfQuote]
	return quoteByte
}

func (c *Client) parsePuzzleResponse(requestData []byte) (byte, []byte) {
	sizeOfPuzzle, targetBits := requestData[tcp.IndexRequestSizeOfPuzzle], requestData[tcp.IndexRequestTargetBits]
	puzzle := requestData[tcp.StartPuzzle : tcp.StartPuzzle+sizeOfPuzzle]
	return targetBits, puzzle
}
