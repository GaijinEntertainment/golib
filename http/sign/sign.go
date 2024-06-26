package sign

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"net/http"
	"slices"
	"time"
)

type Version int

const (
	VErr Version = iota
	_
	V2
	V3
	V4
)

type Signer struct {
	Ver Version
	ts  [8]byte
	h   hash.Hash
}

func NewSigner(v Version, timeStamp time.Time) *Signer {
	if v < V2 || v > V4 {
		panic("unsupported version")
	}

	s := &Signer{
		Ver: v,
		h:   sha256.New(),
	}

	binary.BigEndian.PutUint64(s.ts[:], uint64(timeStamp.Unix()))

	s.h.Write(s.ts[:])

	return s
}

func (s *Signer) AddBytes(data []byte) {
	_, _ = s.h.Write(data)
}

func (s *Signer) AddString(str string) {
	_, _ = s.h.Write([]byte(str))
}

func (s *Signer) AddRequest(r *http.Request) error {
	body, err := getBody(r)
	if err != nil {
		return fmt.Errorf("read request body failed: %w", err)
	}

	if s.Ver == V4 {
		var delimiter = []byte{0}

		s.AddString(r.Method)
		s.AddBytes(delimiter)
		s.AddString(r.Host)
		s.AddBytes(delimiter)
		s.AddString(r.RequestURI)
		s.AddBytes(delimiter)
	}

	s.AddBytes(body)

	return nil
}

func (s *Signer) Sign() Signature {
	data := make([]byte, s.h.Size()+8+1)

	data[0] = byte(s.Ver)
	copy(data[1:], s.ts[:])

	return s.h.Sum(data)
}

func getBody(r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	r.Body = io.NopCloser(bytes.NewReader(body))

	return body, nil
}

type Signature []byte

func ParseSignature(s string) (Signature, error) {
	if len(s) < 1+8<<1 || len(s)%2 != 1 {
		return nil, fmt.Errorf("invalid signature")
	}

	v := Version(s[0] - '0')
	if v < V2 || v > V4 {
		return nil, fmt.Errorf("unsupported version: %d", v)
	}

	data := make([]byte, (len(s)-1)>>1+1)

	data[0] = byte(v)

	_, err := hex.Decode(data[1:], []byte(s[1:]))
	if err != nil {
		return nil, fmt.Errorf("decode signature failed: %w", err)
	}

	return data, nil
}

// Equal compares two signatures. Signatures are equal if they have
// the same version, body and timestamps differ by no more than precision.
func (s Signature) Equal(other Signature, precision time.Duration) bool {
	if len(s) < 9 || len(other) < 9 || len(s) != len(other) {
		return false
	}

	if s[0] != other[0] {
		return false
	}

	if !slices.Equal(s[9:], other[9:]) {
		return false
	}

	t1 := time.Unix(int64(binary.BigEndian.Uint64(s[1:9])), 0)
	t2 := time.Unix(int64(binary.BigEndian.Uint64(other[1:9])), 0)

	dt := t1.Sub(t2)

	if dt < -precision || dt > precision {
		return false
	}

	return true
}

func (s Signature) Ver() Version {
	if len(s) == 0 {
		return VErr
	}

	v := Version(s[0])
	if v < V2 || v > V4 {
		return VErr
	}

	return v
}

func (s Signature) Time() time.Time {
	if len(s) < 9 {
		return time.Time{}
	}

	return time.Unix(int64(binary.BigEndian.Uint64(s[1:9])), 0)
}

func (s Signature) Data() []byte {
	if len(s) < 9 {
		return nil
	}

	return s[9:]
}

// IssuedAt checks if the signature was issued at the specified time.
// The signature is considered to be issued at the specified time if the difference
// between the signature time and the target time is less than the leeway.
// The signature with zero time is always considered to be issued at the right time.
func (s Signature) IssuedAt(t time.Time, leeway time.Duration) bool {
	ts := s.Time()
	if ts.Unix() == 0 {
		return true
	}

	dt := ts.Sub(t)

	return dt >= -leeway && dt <= leeway
}

func (s Signature) HexString() string {
	if len(s) == 0 {
		return ""
	}

	data := make([]byte, len(s)*2-1) // (len(s)-1)*2 + 1

	data[0] = '0' + s[0]
	hex.Encode(data[1:], s[1:])

	return string(data)
}
