package sign_test

import (
	"encoding/binary"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"golib/http/sign"
)

func TestEqual(t *testing.T) {
	t.Parallel()

	type args struct {
		a         sign.Signature
		b         sign.Signature
		precision time.Duration
	}
	tests := []struct {
		name  string
		args  args
		equal bool
	}{
		{
			name: "empty",
			args: args{
				a:         nil,
				b:         nil,
				precision: 0,
			},
			equal: false,
		},
		{
			name: "short a",
			args: args{
				a:         bs(sign.V4, 1000, "a")[:8],
				b:         bs(sign.V4, 1000, "a"),
				precision: 0,
			},
			equal: false,
		},
		{
			name: "short b",
			args: args{
				a:         bs(sign.V4, 1000, "a"),
				b:         bs(sign.V4, 1000, "a")[:8],
				precision: 0,
			},
			equal: false,
		},
		{
			name: "version mismatch",
			args: args{
				a:         bs(sign.V3, 1000, "abc"),
				b:         bs(sign.V4, 1000, "abc"),
				precision: 0,
			},
			equal: false,
		},
		{
			name: "length mismatch",
			args: args{
				a:         bs(sign.V3, 1000, "abc"),
				b:         bs(sign.V3, 1000, "abcd"),
				precision: 0,
			},
			equal: false,
		},
		{
			name: "body mismatch",
			args: args{
				a:         bs(sign.V3, 1000, "abc"),
				b:         bs(sign.V3, 1000, "abd"),
				precision: 0,
			},
			equal: false,
		},
		{
			name: "equal",
			args: args{
				a:         bs(sign.V4, 1000, "abc"),
				b:         bs(sign.V4, 1000, "abc"),
				precision: 0,
			},
			equal: true,
		},
		{
			name: "equal_with_precision",
			args: args{
				a:         bs(sign.V4, 1000, "abc"),
				b:         bs(sign.V4, 1000, "abc"),
				precision: time.Minute,
			},
			equal: true,
		},
		{
			name: "equal_with_precision_1",
			args: args{
				a:         bs(sign.V4, 1010, "abc"),
				b:         bs(sign.V4, 1000, "abc"),
				precision: time.Minute,
			},
			equal: true,
		},
		{
			name: "equal_with_precision_2",
			args: args{
				a:         bs(sign.V4, 1000, "abc"),
				b:         bs(sign.V4, 1010, "abc"),
				precision: time.Minute,
			},
			equal: true,
		},
		{
			name: "out_by_precision_1",
			args: args{
				a:         bs(sign.V4, 1061, "abc"),
				b:         bs(sign.V4, 1000, "abc"),
				precision: time.Minute,
			},
			equal: false,
		},
		{
			name: "out_by_precision_2",
			args: args{
				a:         bs(sign.V4, 1000, "abc"),
				b:         bs(sign.V4, 1061, "abc"),
				precision: time.Minute,
			},
			equal: false,
		},
		{
			name: "ok_with_big_precision",
			args: args{
				a:         bs(sign.V4, 1000, "abc"),
				b:         bs(sign.V4, 1061, "abc"),
				precision: 2 * time.Minute,
			},
			equal: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.args.a.Equal(tt.args.b, tt.args.precision); got != tt.equal {
				t.Errorf("Equal() = %v, want %v", got, tt.equal)
			}
		})
	}
}

func bs(v sign.Version, ts uint64, body string) sign.Signature {
	return append(binary.BigEndian.AppendUint64([]byte{byte(v)}, ts), []byte(body)...)
}

func TestParseSignature(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		s       string
		wantErr bool
	}{
		{
			name:    "empty",
			s:       "",
			wantErr: true,
		},
		{
			name:    "short",
			s:       "2AA",
			wantErr: true,
		},
		{
			name:    "invalid version",
			s:       "511111111111111110000",
			wantErr: true,
		},
		{
			name:    "uneven data",
			s:       "5111111111111111100000",
			wantErr: true,
		},
		{
			name:    "invalid data",
			s:       "411111111111111110000Z0",
			wantErr: true,
		},
		{
			name:    "ok",
			s:       "411111111111111110123456789ABCDEF",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := sign.ParseSignature(tt.s)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("ParseSignature() error = %v, wantErr %v", err, tt.wantErr)
				}

				return
			}

			if tt.s != got.HexString() {
				t.Errorf("ParseSignature() got = %v, want %v", got, tt.s)
			}
		})
	}

	s, err := sign.ParseSignature("300000000000000200203040506")

	assert.NoError(t, err)
	assert.Equal(t, sign.V3, s.Ver())
	assert.Equal(t, time.Unix(32, 0), s.Time())
	assert.Equal(t, []byte{2, 3, 4, 5, 6}, s.Data())
}
