package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

var (
	errInIsNil  = errors.New("in is nil")
	errOutIsNil = errors.New("out is nil")

	EOF                 = "...EOF"
	ConnectionWasClosed = "...Connection was closed by peer"
)

type TelnetClient interface {
	Connect() error
	Close() error
	Send() error
	Receive() error
}

type Telnet struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	log     io.Writer
	conn    net.Conn
	cancel  context.CancelFunc
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer, log io.Writer, cancel context.CancelFunc) TelnetClient {
	return &Telnet{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
		log:     log,
		cancel:  cancel,
	}
}

func (t *Telnet) Connect() error {
	if t.in == nil {
		return errInIsNil
	}
	if t.out == nil {
		return errOutIsNil
	}

	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return fmt.Errorf("can't connect: %w", err)
	}

	t.conn = conn
	return nil
}

func (t *Telnet) Close() error {
	return t.conn.Close()
}

func (t *Telnet) Send() error {
	defer t.cancel()
	_, err := io.Copy(t.conn, t.in)
	if err != nil {
		return fmt.Errorf("can't send: %w", err)
	}
	t.logMsg(EOF)
	return nil
}

func (t *Telnet) Receive() error {
	defer t.cancel()
	_, err := io.Copy(t.out, t.conn)
	if err != nil {
		return fmt.Errorf("can't receive: %w", err)
	}
	t.logMsg(ConnectionWasClosed)
	return nil
}

func (t *Telnet) logMsg(s string) {
	if t.log != nil {
		fmt.Fprintln(t.log, s)
	}
}
