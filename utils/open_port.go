package utils

import (
	"errors"
	"fmt"
	"net"
	"slices"
	"syscall"
	"time"
)

var (
	ErrMustBePositiveTimeout = errors.New("timeout on FindOpenLocalPort must be > 0")
	ErrPrivilegedPort        = errors.New("cannot check privileged port")
	ErrNoOpenPorts           = errors.New("no ports open in the requested range")
)

const (
	lowestPort  uint16 = 1025
	highestPort uint16 = 65535
)

// FindOpenLocalPort will find an open tcp port among the given ports, waiting up to timeout before the next port is checked
// If there is exactly one port given, it will search between that port and 65535 for the first available port
// If there are exactly two ports given, it will search the closed range between [low, high]
// If there are more than 3 ports, it will search only those ports
// If there are none, it will search the entire non-privileged port range (1025 - 65535)
// If any ports are < 1025, an error will be thrown
// If there are no open ports, an error will be thrown
func FindOpenLocalPort(timeout time.Duration, ports ...uint16) (uint16, error) {
	if timeout <= 0 {
		return 0, ErrMustBePositiveTimeout
	}

	slices.Sort(ports)

	if slices.ContainsFunc(ports, func(port uint16) bool {
		return port < lowestPort
	}) {
		return 0, ErrPrivilegedPort
	}

	exact := false
	switch len(ports) {
	case 0:
		ports = []uint16{lowestPort, highestPort}
	case 1:
		ports = []uint16{ports[0], highestPort}
	case 2:
		// nothing to do here, but separate it from the default case
	default:
		exact = true
	}

	if exact {
		for _, port := range ports {
			open, err := isOpen(timeout, port)
			if err != nil || open {
				return port, err
			}
		}

		return 0, ErrNoOpenPorts
	}

	low, high := ports[0], ports[1]
	for portOffset := range high - low {
		port := low + portOffset
		open, err := isOpen(timeout, port)
		if err != nil || open {
			return port, err
		}
	}

	return 0, ErrNoOpenPorts
}

func isOpen(timeout time.Duration, port uint16) (bool, error) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), timeout)
	if errors.Is(err, syscall.ECONNREFUSED) {
		return true, nil
	}
	conn.Close()
	return false, err
}
