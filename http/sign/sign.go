package sign

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"net/http"
	"slices"
	"time"
)

type Version int

var (
	ErrInvalidSignature              = errors.New("invalid signature")
	ErrorUnsupportedSignatureVersion = errors.New("unsupported signature version")
)

const (
	VErr Version = iota
	_
	V2
	V3
	V4
)

type Digest struct {
	Ver Version
	ts  [8]byte
	h   hash.Hash
}

func NewDigest(v Version, timeStamp time.Time) *Digest {
	if v < V2 || v > V4 {
		panic("unsupported version")
	}

	s := &Digest{
		Ver: v,
		h:   sha256.New(),
		ts:  [8]byte{},
	}

	binary.BigEndian.PutUint64(s.ts[:], uint64(timeStamp.Unix())) //nolint:gosec

	_, _ = s.h.Write(s.ts[:])

	return s
}

func (s *Digest) AddBytes(data []byte) {
	if len(data) == 0 {
		return
	}

	_, _ = s.h.Write(data)
}

func (s *Digest) AddString(str string) {
	s.AddBytes([]byte(str))
}

func (s *Digest) AddRequest(r *http.Request) error {
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
		s.AddString(r.URL.Path)
		s.AddBytes(delimiter)
		s.AddString(r.URL.RawQuery)
		s.AddBytes(delimiter)
	}

	s.AddBytes(body)

	return nil
}

func (s *Digest) Sum(b []byte) []byte {
	return s.h.Sum(b)
}

func (s *Digest) Sign(signer Signer) (Signature, error) {
	signature, err := signer.Sign(s.Sum(nil))
	if err != nil {
		return nil, fmt.Errorf("sign failed: %w", err)
	}

	data := make([]byte, hdrSize+len(signature))

	data[0] = byte(s.Ver)
	copy(data[timeOffset:], s.ts[:])
	copy(data[signOffset:], signature)

	return data, nil
}

func getBody(r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	r.Body = io.NopCloser(bytes.NewReader(body))

	return body, nil
}

const (
	timeOffset = 1
	signOffset = 9
	hdrSize    = 1 + 8
)

type Signature []byte

func ParseSignature(s string) (Signature, error) {
	if len(s) < hdrSize || len(s)%2 != 1 {
		return nil, ErrInvalidSignature
	}

	v := Version(s[0] - '0')
	if v < V2 || v > V4 {
		return nil, ErrorUnsupportedSignatureVersion
	}

	data := make([]byte, (len(s)-1)>>1+1)

	data[0] = byte(v)

	_, err := hex.Decode(data[timeOffset:], []byte(s[timeOffset:]))
	if err != nil {
		return nil, fmt.Errorf("decode signature failed: %w", err)
	}

	return data, nil
}

// Equal compares two signatures. Signatures are equal if they have
// the same version, body and timestamps differ by no more than precision.
func (s Signature) Equal(other Signature, precision time.Duration) bool {
	if len(s) < hdrSize || len(other) < hdrSize || len(s) != len(other) {
		return false
	}

	if s[0] != other[0] {
		return false
	}

	if !slices.Equal(s[signOffset:], other[signOffset:]) {
		return false
	}

	dt := s.Time().Sub(other.Time())

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
	if len(s) < hdrSize {
		return time.Time{}
	}

	return time.Unix(int64(binary.BigEndian.Uint64(s[timeOffset:hdrSize])), 0) //nolint:gosec
}

func (s Signature) Data() []byte {
	if len(s) < hdrSize {
		return nil
	}

	return s[hdrSize:]
}

// IssuedAt checks if the signature was issued at the specified time.
// The signature is considered to be issued at the specified time if the difference
// between the signature time and the target time is less than the leeway.
// The signature with zero time is always considered to be issued at the specified time.
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

type Signer interface {
	Sign(digest []byte) (Signature, error)
}

type RSASigner rsa.PrivateKey

func (key *RSASigner) Sign(digest []byte) (Signature, error) {
	return rsa.SignPKCS1v15(nil, (*rsa.PrivateKey)(key), crypto.SHA256, digest) //nolint:wrapcheck
}
