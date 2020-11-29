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
)

type TelnetClient interface {
	Connect() error
	Close() error
	Send() error
	Receive() error
	Done() <-chan struct{}
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	ctx, cancel := context.WithCancel(context.Background())
	return &Telnet{
		ctx:     ctx,
		cancel:  cancel,
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

type Telnet struct {
	ctx     context.Context
	cancel  context.CancelFunc
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
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
	defer t.cancel()
	return t.conn.Close()
}

func (t *Telnet) Send() error {
	_, err := io.Copy(t.conn, t.in)
	if err != nil {
		return fmt.Errorf("can't send: %w", err)
	}
	return nil
}

func (t *Telnet) Receive() error {
	_, err := io.Copy(t.out, t.conn)
	if err != nil {
		return fmt.Errorf("can't receive: %w", err)
	}
	return nil
}

func (t *Telnet) Done() <-chan struct{} {
	return t.ctx.Done()
}
