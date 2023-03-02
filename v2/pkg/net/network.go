package network

import (
	"log"
	"net"
	"strings"

	"github.com/jackpal/gateway"
)

// HostAddresses returns all available ip addresses from all interfaces
func HostAddresses() ([]string, error) {
	ifaceAddrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	hostAddrs := []string{}
	// run the loop backwards so 127.0.0.1 ends up at the bottom of the list
	for idx := len(ifaceAddrs) - 1; 0 <= idx; idx-- {
		addr := ifaceAddrs[idx].String()
		if strings.Contains(addr, "::") {
			continue
		}

		parts := strings.Split(addr, "/")
		ip := net.ParseIP(parts[0])
		if ip == nil {
			continue
		}

		// Filter out IPv6 address
		if ip.To4() == nil {
			continue
		}

		hostAddrs = append(hostAddrs, ip.String())
	}

	return hostAddrs, nil
}

// GetSourceIP returns the ip address used to access target:port
// If target is the empty string then the default gateway ip is used.
// If the port is the empty string, then "80" is used by default.
func GetSourceIP(target, port string) (string, error) {
	if target == "" {
		ip, err := gateway.DiscoverGateway()
		if err != nil {
			return "", err
		}
		target = ip.String()
	}

	if port == "" {
		port = "80"
	}

	conn, err := net.Dial("udp", target+":"+port)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String(), nil
}

func PreferNonPrivateIP(ips []string) (string, string) {
	if len(ips) == 0 {
		return "", ""
	}

	var (
		preferredAddress net.IP
		port             string
	)

	for _, addr := range ips {
		host, p, err := net.SplitHostPort(addr)
		if err != nil {
			log.Printf("Failed to parse peer address: %s", err.Error())
			continue
		}
		ip := net.ParseIP(host)
		if ip == nil {
			log.Printf("Failed to parse peer address: %s", host)
			continue
		}

		if !ip.IsPrivate() {
			preferredAddress = ip
			port = p
			break
		}
	}

	if preferredAddress == nil {
		host, p, err := net.SplitHostPort(ips[0])
		if err != nil {
			log.Printf("Failed to parse peer address: %s", err)
		} else {
			preferredAddress = net.ParseIP(host)
			port = p
		}
	}

	return preferredAddress.String(), port
}
