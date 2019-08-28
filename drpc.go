// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

// drpc is a light replacement for gprc.
package drpc

import (
	"context"
	"io"

	"github.com/zeebo/errs"
)

var (
	Error         = errs.Class("drpc")
	InternalError = errs.Class("internal error")
	ProtocolError = errs.Class("protocol error")
)

type Transport interface {
	io.Reader
	io.Writer
	io.Closer
}

type Message interface {
	Reset()
	String() string
	ProtoMessage()
}

type Conn interface {
	Close() error
	Transport() Transport

	Invoke(ctx context.Context, rpc string, in, out Message) error
	NewStream(ctx context.Context, rpc string) (Stream, error)
}

type Stream interface {
	Context() context.Context

	MsgSend(msg Message) error
	MsgRecv(msg Message) error

	CloseSend() error
	Close() error
}

type Handler = func(srv interface{}, ctx context.Context, in1, in2 interface{}) (out Message, err error)

type Description interface {
	NumMethods() int
	Method(n int) (rpc string, handler Handler, method interface{}, ok bool)
}
