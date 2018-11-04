package goredis

import (
	"bufio"
	"errors"
	"io"
	"strconv"
)

var (
	intBuffer [][]byte
	respTerm  = []byte("\r\n")
	nullBulk  = []byte("-1")
	nullArray = []byte("-1")
)

func init() {
	cnt := 10000
	intBuffer = make([][]byte, cnt)
	for i := 0; i < cnt; i++ {
		intBuffer[i] = []byte(strconv.Itoa(i))
	}
}

type Error string

func (err Error) Error() string {
	return string(err)
}

var (
	okReply   interface{} = "OK"
	pongReply interface{} = "PONG"
)

type RespReader struct {
	br *bufio.Reader
}

type RespWriter struct {
	bw *bufio.Writer
	// Scratch space for formatting integers and floats.
	numScratch [40]byte
}

func NewRespReader(br *bufio.Reader) *RespReader {
	r := &RespReader{br}
	return r
}

// Parse RESP
func (resp *RespReader) Parse() (interface{}, error) {
	line, err := readLine(resp.br)
	if err != nil {
		return nil, err
	}
	if len(line) == 0 {
		return nil, errors.New("short resp line")
	}

	switch line[0] {
	case '+':
		switch {
		case len(line) == 3 && line[1] == 'O' && line[2] == 'K':
			// Avoid allocation for frequent "+OK" response.
			return okReply, nil
		case len(line) == 5 && line[1] == 'P' && line[2] == 'O' && line[3] == 'N' && line[4] == 'G':
			// Avoid allocation in PING command benchmarks :)
			return pongReply, nil
		default:
			return string(line[1:]), nil
		}
	case '-':
		return Error(string(line[1:])), nil
	case ':':
		n, err := parseInt(line[1:])
		return n, err
	case '$':
		n, err := parseLen(line[1:])
		if n < 0 || err != nil {
			return nil, err
		}

		p := make([]byte, n)
		_, err = io.ReadFull(resp.br, p)
		if err != nil {
			return nil, err
		}

		if line, err := readLine(resp.br); err != nil {
			return nil, err
		} else if len(line) != 0 {
			return nil, errors.New("bad bulk string for format")
		}

		return p, nil
	case '*':
		n, err := parseLen(line[1:])
		if n < 0 || err != nil {
			return nil, err
		}

		r := make([]interface{}, n)
		for i := range r {
			r[i], err = resp.Parse()
			if err != nil {
				return nil, err
			}
		}

		return r, nil
	}

	return nil, errors.New("unexpected response line")
}

func readLine(br *bufio.Reader) ([]byte, error) {
	p, err := br.ReadSlice('\n')
	if err == bufio.ErrBufferFull {
		return nil, errors.New("long resp line")
	}
	if err != nil {
		return nil, err
	}

	i := len(p) - 2
	if i < 0 || p[i] != '\r' {
		return nil, errors.New("bad resp line terminator")
	}

	return p[:i], nil
}

// parseInt parses an integer reply.
func parseInt(p []byte) (int64, error) {
	if len(p) == 0 {
		return 0, errors.New("malformed integer")
	}

	var negate bool
	if p[0] == '-' {
		negate = true
		p = p[1:]
		if len(p) == 0 {
			return 0, errors.New("malformed integer")
		}
	}

	var n int64
	for _, b := range p {
		n *= 10
		if b < '0' || b > '9' {
			return 0, errors.New("illegal bytes in length")
		}

		n += int64(b - '0')
	}

	if negate {
		n = -n
	}

	return n, nil
}

// parseLen parses bulk string and array lengths.
func parseLen(p []byte) (int, error) {
	if len(p) == 0 {
		return -1, errors.New("malformed length")
	}

	if p[0] == '-' && len(p) == 2 && p[1] == '1' {
		// handle $-1 and $-1 null replies.
		return -1, nil
	}

	var n int
	for _, b := range p {
		n *= 10
		if b < '0' || b > '9' {
			return -1, errors.New("illegal bytes in length")
		}

		n += int(b - '0')
	}

	return n, nil
}

func NewRespWriter(bw *bufio.Writer) *RespWriter {
	r := &RespWriter{bw: bw}
	return r
}

// RESP command is array of bulk string
func (resp *RespWriter) WriteCommand(cmd string, args ...interface{}) error {

	return nil
}
