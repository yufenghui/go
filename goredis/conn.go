package goredis

import (
	"bufio"
	"io"
	"log"
	"net"
	"sync/atomic"
)

type sizeWriter int64

func (s *sizeWriter) Write(p []byte) (int, error) {
	*s += sizeWriter(len(p))
	return len(p), nil
}

type Conn struct {
	c net.Conn

	respReader *RespReader
	respWriter *RespWriter

	totalReadSize sizeWriter
	totalWritSize sizeWriter

	closed int32
}

func ConnectWithSize(addr string, readSize int, writeSize int) (*Conn, error) {
	conn, err := net.Dial(getProto(addr), addr)
	if err != nil {
		return nil, err
	}

	return NewConnWithSize(conn, readSize, writeSize)
}

func NewConn(conn net.Conn) (*Conn, error) {
	return NewConnWithSize(conn, 1024, 1024)
}

func NewConnWithSize(conn net.Conn, readSize int, writeSize int) (*Conn, error) {
	c := new(Conn)

	c.c = conn

	br := bufio.NewReaderSize(io.TeeReader(c.c, &c.totalReadSize), readSize)
	bw := bufio.NewWriterSize(io.MultiWriter(c.c, &c.totalWritSize), writeSize)

	c.respReader = NewRespReader(br)
	c.respWriter = NewRespWriter(bw)

	atomic.StoreInt32(&c.closed, 0)

	return c, nil
}

func (c *Conn) Close() {
	if atomic.LoadInt32(&c.closed) == 1 {
		return
	}

	c.c.Close()

	atomic.StoreInt32(&c.closed, 1)
}

func (c *Conn) isClosed() bool {
	return atomic.LoadInt32(&c.closed) == 1
}

func (c *Client) newConn(addr string, pass string) (*Conn, error) {
	co, err := ConnectWithSize(addr, c.readBufferSize, c.writeBufferSize)
	if err != nil {
		return nil, err
	}

	if len(pass) > 0 {
		_, err = co.Do("AUTH", pass)
		if err != nil {
			co.Close()
			return nil, err
		}
	}

	return co, nil
}

// Send RESP command and receive the reply
func (c *Conn) Do(cmd string, args ...interface{}) (interface{}, error) {
	if err := c.Send(cmd, args...); err != nil {
		return nil, err
	}

	return c.Receive()
}

// Send RESP command
func (c *Conn) Send(cmd string, args ...interface{}) error {
	if err := c.respWriter.WriteCommand(cmd, args...); err != nil {
		c.Close()
		return err
	}

	return nil
}

// Receive RESP reply
func (c *Conn) Receive() (interface{}, error) {
	log.Printf("receive data: %s", c.respReader.br)

	if reply, err := c.respReader.Parse(); err != nil {
		c.Close()
		return nil, err
	} else {
		if e, ok := reply.(Error); ok {
			return reply, e
		} else {
			return reply, nil
		}
	}
}
