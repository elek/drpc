// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package drpcwire

import (
	"bytes"
	"io"
	"testing"

	"github.com/zeebo/assert"
)

func TestWriter(t *testing.T) {
	run := func(size int) func(t *testing.T) {
		return func(t *testing.T) {
			var exp []byte
			var got bytes.Buffer

			wr := NewWriter(&got, size)
			for i := 0; i < 1000; i++ {
				fr := RandFrame()
				exp = AppendFrame(exp, fr)
				assert.NoError(t, wr.WriteFrame(fr))
			}
			assert.NoError(t, wr.Flush())
			assert.That(t, bytes.Equal(exp, got.Bytes()))
		}
	}

	t.Run("Size 0B", run(0))
	t.Run("Size 1MB", run(1024*1024))
}

func TestReader(t *testing.T) {
	type testCase struct {
		Packet Packet
		Frames []Frame
	}

	p := func(kind Kind, id uint64, data string) Packet {
		return Packet{
			Data: []byte(data),
			ID:   ID{Stream: 1, Message: id},
			Kind: kind,
		}
	}

	f := func(kind Kind, id uint64, data string, done bool) Frame {
		return Frame{
			Data: []byte(data),
			ID:   ID{Stream: 1, Message: id},
			Kind: kind,
			Done: done,
		}
	}

	m := func(pkt Packet, frames ...Frame) testCase {
		return testCase{
			Packet: pkt,
			Frames: frames,
		}
	}

	cases := []testCase{
		m(p(Kind_Message, 1, "hello world"),
			f(Kind_Message, 1, "hello", false),
			f(Kind_Message, 1, " ", false),
			f(Kind_Message, 1, "world", true)),

		m(p(Kind_Cancel, 2, ""),
			f(Kind_Message, 1, "hello", false),
			f(Kind_Message, 1, " ", false),
			f(Kind_Cancel, 2, "", true)),
	}

	for _, tc := range cases {
		var buf []byte
		for _, fr := range tc.Frames {
			buf = AppendFrame(buf, fr)
		}

		rd := NewReader(bytes.NewReader(buf))
		pkt, err := rd.ReadPacket()
		assert.NoError(t, err)
		assert.DeepEqual(t, tc.Packet, pkt)
		_, err = rd.ReadPacket()
		assert.Equal(t, err, io.EOF)
	}
}
