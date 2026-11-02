// Package frame implements native-messaging framing: a 32-bit length in the
// machine's native byte order followed by that many bytes of UTF-8 JSON.
package frame

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	// MaxInbound is the largest message the browser may send to the host.
	MaxInbound = 64 << 20
	// MaxOutbound is the largest message the host may send to the browser
	// (1 MB). Larger messages make the browser drop the connection.
	MaxOutbound = 1 << 20
)

// ErrTooLarge is returned when a frame exceeds the configured limit.
var ErrTooLarge = errors.New("frame: message too large")

// Read reads one browser->host frame. It returns io.EOF when stdin is closed
// cleanly between frames.
func Read(r io.Reader) ([]byte, error) { return ReadMax(r, MaxInbound) }

// Write writes one host->browser frame, enforcing the 1 MB limit.
func Write(w io.Writer, payload []byte) error { return WriteMax(w, payload, MaxOutbound) }

// ReadMax reads one frame of at most max bytes.
func ReadMax(r io.Reader, max int) ([]byte, error) {
	var hdr [4]byte
	if _, err := io.ReadFull(r, hdr[:]); err != nil {
		return nil, err
	}
	n := binary.NativeEndian.Uint32(hdr[:])
	if int64(n) > int64(max) {
		return nil, fmt.Errorf("%w: inbound frame of %d bytes (max %d)", ErrTooLarge, n, max)
	}
	buf := make([]byte, n)
	if _, err := io.ReadFull(r, buf); err != nil {
		if errors.Is(err, io.EOF) {
			err = io.ErrUnexpectedEOF
		}
		return nil, err
	}
	return buf, nil
}

// WriteMax writes one frame of at most max bytes in a single Write call so
// concurrent writers (guarded by a mutex) never interleave header and body.
func WriteMax(w io.Writer, payload []byte, max int) error {
	if len(payload) > max {
		return fmt.Errorf("%w: outbound frame of %d bytes (max %d)", ErrTooLarge, len(payload), max)
	}
	buf := make([]byte, 4, 4+len(payload))
	binary.NativeEndian.PutUint32(buf, uint32(len(payload)))
	buf = append(buf, payload...)
	_, err := w.Write(buf)
	return err
}
