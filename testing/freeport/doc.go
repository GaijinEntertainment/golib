// Package freeport provides small helpers to allocate free TCP/UDP ports
// in tests. It is useful when starting ephemeral servers (HTTP, gRPC,
// WebSocket, custom TCP/UDP protocols) inside integration/unit tests, where
// ports must not collide between concurrently running tests.
//
// # Why
//
//   - Tests often need to bind to a local port but hard‑coding ports leads to
//     flakiness and conflicts on CI machines and during parallel runs.
//   - The package offers two styles: probing for a free port (TCP/UDP) and
//     reserving a port for the lifetime of a test (ReserveTCP/ReserveUDP).
//   - Port reservation is extremely useful for the cases when testing against
//     docker container and container must be stopped and started multiple times
//     during a test, and each time it must bind to the same port.
//   - Probing is inherently racy (another process could claim the port after
//     probing but before binding), but is often good enough for quick tests.
//   - Reservation is not perfect either (another process could still claim the
//     port), but is much safer when used in tests that all use this package.
//
// How to use
//
//  1. Quick probe (no reservation):
//
//     p, err := freeport.TCP()
//     // use p to bind a listener; another process could still claim it
//
//  2. Safer within tests (reservation tied to t):
//
//     p := freeport.ReserveTCP(t)
//     ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", p))
//     // close the listener yourself when done; the port reservation is
//     // released automatically at the end of the test
//
// # Notes
//
//   - Reservation is best‑effort and only prevents collisions among tests that
//     also use this package in the same process. It does not prevent other
//     external processes from binding the port.
//   - ReleaseTCP/ReleaseUDP return whether a release actually happened; this is
//     mostly useful for tests that juggle many ports and want to free them early.
package freeport
