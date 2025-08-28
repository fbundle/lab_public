package rpc

import (
	"context"
	"encoding/binary"
	"io"
)

// MessageIO - breaking stream of bytes into messages, possibly include encryption-decryption
type MessageIO interface {
	Write(ctx context.Context, w io.Writer, b []byte) (err error)
	Read(ctx context.Context, r io.Reader) (b []byte, err error)
}

func NewMessageIO() MessageIO {
	return lengthPrefixMessageIO{
		putUint: func(b []byte, v uint64) {
			binary.BigEndian.PutUint32(b, uint32(v))
		},
		getUint: func(b []byte) (v uint64) {
			return uint64(binary.BigEndian.Uint32(b))
		},
		uintSize: 4,
	}
}

type lengthPrefixMessageIO struct {
	putUint  func(b []byte, v uint64)
	getUint  func(b []byte) (v uint64)
	uintSize int
}

func (msgIO lengthPrefixMessageIO) Write(ctx context.Context, w io.Writer, b []byte) error {
	sizeBuf := make([]byte, msgIO.uintSize)
	msgIO.putUint(sizeBuf, uint64(len(b)))

	if err := writeFull(ctx, w, sizeBuf); err != nil {
		return err
	}
	return writeFull(ctx, w, b)
}

func (msgIO lengthPrefixMessageIO) Read(ctx context.Context, r io.Reader) (b []byte, err error) {
	sizeBuf := make([]byte, msgIO.uintSize)
	if err := readFull(ctx, r, sizeBuf); err != nil {
		return nil, err
	}

	size := msgIO.getUint(sizeBuf)
	b = make([]byte, size)
	return b, readFull(ctx, r, b)
}

func writeFull(ctx context.Context, w io.Writer, b []byte) error {
	offset := 0
	for offset < len(b) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		n, err := w.Write(b[offset:])
		if err != nil {
			return err
		}
		if n == 0 {
			return io.ErrShortWrite
		}
		offset += n
	}
	return nil
}

func readFull(ctx context.Context, r io.Reader, b []byte) error {
	offset := 0
	for offset < len(b) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		n, err := r.Read(b[offset:])
		if err != nil {
			return err
		}
		if n == 0 {
			return io.ErrUnexpectedEOF
		}
		offset += n
	}
	return nil
}
