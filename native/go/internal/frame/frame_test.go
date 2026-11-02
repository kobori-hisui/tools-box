package frame

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"testing"
)

func TestRoundTrip(t *testing.T) {
	var buf bytes.Buffer
	for _, msg := range []string{`{"a":1}`, ``, `{"emoji":"✓"}`} {
		if err := Write(&buf, []byte(msg)); err != nil {
			t.Fatal(err)
		}
	}
	for _, want := range []string{`{"a":1}`, ``, `{"emoji":"✓"}`} {
		got, err := Read(&buf)
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != want {
			t.Fatalf("got %q want %q", got, want)
		}
	}
	if _, err := Read(&buf); err != io.EOF {
		t.Fatalf("expected io.EOF at clean end, got %v", err)
	}
}

func TestWriteTooLarge(t *testing.T) {
	err := Write(io.Discard, make([]byte, MaxOutbound+1))
	if !errors.Is(err, ErrTooLarge) {
		t.Fatalf("expected ErrTooLarge, got %v", err)
	}
	if err := Write(io.Discard, make([]byte, MaxOutbound)); err != nil {
		t.Fatalf("exactly MaxOutbound must be allowed: %v", err)
	}
}

func TestReadTooLarge(t *testing.T) {
	var hdr [4]byte
	binary.NativeEndian.PutUint32(hdr[:], MaxInbound+1)
	if _, err := Read(bytes.NewReader(hdr[:])); !errors.Is(err, ErrTooLarge) {
		t.Fatalf("expected ErrTooLarge, got %v", err)
	}
}

func TestTruncatedPayload(t *testing.T) {
	var buf bytes.Buffer
	_ = Write(&buf, []byte("hello"))
	trunc := bytes.NewReader(buf.Bytes()[:6])
	if _, err := Read(trunc); !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("expected ErrUnexpectedEOF, got %v", err)
	}
}
