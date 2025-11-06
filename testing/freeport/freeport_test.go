package freeport_test

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"

	"dev.gaijin.team/go/golib/testing/freeport"
)

func TestGet(t *testing.T) {
	t.Parallel()

	t.Run("TCP", func(t *testing.T) {
		t.Parallel()

		p, err := freeport.TCP()

		require.NoError(t, err)
		require.Positive(t, p)

		// Verify we can listen on the port
		l, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: p, Zone: ""})
		require.NoError(t, err)
		require.NotNil(t, l)

		_ = l.Close()
	})

	t.Run("UDP", func(t *testing.T) {
		t.Parallel()

		p, err := freeport.UDP()

		require.NoError(t, err)
		require.Positive(t, p)

		// Verify we can bind UDP to the port
		c, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: p, Zone: ""})
		require.NoError(t, err)
		require.NotNil(t, c)

		_ = c.Close()
	})
}

func TestReserve(t *testing.T) {
	t.Parallel()

	t.Run("TCP", func(t *testing.T) {
		t.Parallel()

		p := freeport.ReserveTCP(t)
		require.Positive(t, p)

		t.Run("foreign-release", func(t2 *testing.T) { //nolint:paralleltest
			require.False(t2, freeport.ReleaseTCP(t2, p))
		})

		require.True(t, freeport.ReleaseTCP(t, p))
		require.False(t, freeport.ReleaseTCP(t, p))
	})

	t.Run("UDP", func(t *testing.T) {
		t.Parallel()

		p := freeport.ReserveUDP(t)
		require.Positive(t, p)

		t.Run("foreign-release", func(t2 *testing.T) { //nolint:paralleltest
			require.False(t2, freeport.ReleaseUDP(t2, p))
		})

		require.True(t, freeport.ReleaseUDP(t, p))
		require.False(t, freeport.ReleaseUDP(t, p))
	})

	t.Run("Panics", func(t *testing.T) {
		t.Parallel()

		t.Run("ReserveTCP_NilT", func(t *testing.T) {
			t.Parallel()
			require.Panics(t, func() { freeport.ReserveTCP(nil) })
		})

		t.Run("ReserveUDP_NilT", func(t *testing.T) {
			t.Parallel()
			require.Panics(t, func() { freeport.ReserveUDP(nil) })
		})

		t.Run("ReleaseTCP_NilT", func(t *testing.T) {
			t.Parallel()
			require.Panics(t, func() { freeport.ReleaseTCP(nil, 1) })
		})

		t.Run("ReleaseUDP_NilT", func(t *testing.T) {
			t.Parallel()
			require.Panics(t, func() { freeport.ReleaseUDP(nil, 1) })
		})
	})
}
