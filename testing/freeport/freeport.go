package freeport

import (
	"net"
	"sync"
	"testing"

	"dev.gaijin.team/go/golib/e"
)

// resolver manages port allocation and reservation for a specific protocol.
// It maintains a registry of ports that are currently reserved by tests
// and provides methods to get free ports or reserve ports exclusively.
type resolver struct {
	// reserved maps port numbers to the test that reserved them
	reserved map[int]*testing.T
	// reservedMu protects access to the reserved map
	reservedMu sync.Mutex `exhaustruct:"optional"`

	// resolverFn is the function used to find free ports for this protocol
	resolverFn func() (int, error)
}

func (r *resolver) getPort() (int, error) {
	r.reservedMu.Lock()
	defer r.reservedMu.Unlock()

	for {
		port, err := r.resolverFn()
		if err != nil {
			return 0, e.NewFrom("failed to get free port", err)
		}

		if _, ok := r.reserved[port]; !ok {
			return port, nil
		}
	}
}

func (r *resolver) reservePort(t *testing.T) (int, error) { //nolint:thelper
	if t == nil {
		panic("t is required to reserve a port")
	}

	t.Helper()

	r.reservedMu.Lock()
	defer r.reservedMu.Unlock()

	for {
		port, err := r.resolverFn()
		if err != nil {
			return 0, e.NewFrom("failed to reserve free port", err)
		}

		if _, ok := r.reserved[port]; !ok {
			r.reserved[port] = t
			t.Cleanup(func() {
				r.releasePort(t, port)
			})

			return port, nil
		}
	}
}

func (r *resolver) releasePort(t *testing.T, port int) bool { //nolint:thelper
	if t == nil {
		panic("t is required to release a port")
	}

	t.Helper()

	r.reservedMu.Lock()
	defer r.reservedMu.Unlock()

	if owner, ok := r.reserved[port]; ok && owner == t {
		delete(r.reserved, port)
		return true
	}

	return false
}

var (
	tcpResolv = &resolver{ //nolint:gochecknoglobals
		reserved: make(map[int]*testing.T),
		resolverFn: func() (int, error) {
			addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
			if err != nil {
				return 0, e.NewFrom("failed to resolve TCP address", err)
			}

			listener, err := net.ListenTCP("tcp", addr)
			if err != nil {
				return 0, e.NewFrom("failed to listen on TCP port", err)
			}
			defer func() { _ = listener.Close() }()

			return listener.Addr().(*net.TCPAddr).Port, nil //nolint:forcetypeassert
		},
	}
	udpResolv = &resolver{ //nolint:gochecknoglobals
		reserved: make(map[int]*testing.T),
		resolverFn: func() (int, error) {
			addr, err := net.ResolveUDPAddr("udp", "localhost:0")
			if err != nil {
				return 0, e.NewFrom("failed to resolve UDP address", err)
			}

			conn, err := net.ListenUDP("udp", addr)
			if err != nil {
				return 0, e.NewFrom("failed to listen on UDP port", err)
			}
			defer func() { _ = conn.Close() }()

			return conn.LocalAddr().(*net.UDPAddr).Port, nil //nolint:forcetypeassert
		},
	}
)

// TCP returns a free TCP port that is available for use.
// This function does not reserve the port, so there's a small chance
// another process could claim it before you use it.
func TCP() (int, error) {
	return tcpResolv.getPort()
}

// UDP returns a free UDP port that is available for use.
// This function does not reserve the port, so there's a small chance
// another process could claim it before you use it.
func UDP() (int, error) {
	return udpResolv.getPort()
}

// ReserveTCP reserves a TCP port for exclusive use within the test.
// The port is automatically released when the test completes.
//
// There is no guarantee the port will not be used by another process
// or code that does not use this package.
func ReserveTCP(t *testing.T) int {
	t.Helper()

	port, err := tcpResolv.reservePort(t)
	if err != nil {
		t.Fatalf("ReserveTCP: %v", err)
	}

	return port
}

// ReserveUDP reserves a UDP port for exclusive use within the test.
// The port is automatically released when the test completes.
//
// There is no guarantee the port will not be used by another process
// or code that does not use this package.
func ReserveUDP(t *testing.T) int {
	t.Helper()

	port, err := udpResolv.reservePort(t)
	if err != nil {
		t.Fatalf("ReserveUDP: %v", err)
	}

	return port
}

// ReleaseTCP manually releases a TCP port that was reserved for the test. This
// is typically not needed as ports are automatically released when the test
// completes, but can be useful for tests that reserve many ports and want to
// free them early. In case *testing.T does not own the port, or it already was
// released - this function is noop.
//
// Returns true if the port was released, false otherwise.
func ReleaseTCP(t *testing.T, port int) bool {
	t.Helper()
	return tcpResolv.releasePort(t, port)
}

// ReleaseUDP manually releases a UDP port that was reserved for the test. This
// is typically not needed as ports are automatically released when the test
// completes, but can be useful for tests that reserve many ports and want to
// free them early. In case *testing.T does not own the port, or it already was
// released - this function is noop.
//
// Returns true if the port was released, false otherwise.
func ReleaseUDP(t *testing.T, port int) bool {
	t.Helper()
	return udpResolv.releasePort(t, port)
}
