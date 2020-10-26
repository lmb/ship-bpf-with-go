package main

import (
	"syscall"
)

func main() {
	const SO_ATTACH_BPF = 50
	const loopback = 1
}

func openRawSock(index int) (int, error) {
	const ETH_P_ALL uint16 = 0x300
	sock, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW|syscall.SOCK_NONBLOCK|syscall.SOCK_CLOEXEC, int(ETH_P_ALL))
	if err != nil {
		return 0, err
	}
	sll := syscall.SockaddrLinklayer{}
	sll.Protocol = ETH_P_ALL
	sll.Ifindex = index
	if err := syscall.Bind(sock, &sll); err != nil {
		return 0, err
	}
	return sock, nil
}
