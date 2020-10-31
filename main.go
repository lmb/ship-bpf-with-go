package main

import (
	"fmt"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go counter counter.c -- -I./include -nostdinc -O3

func main() {
	const SO_ATTACH_BPF = 50
	const loopback = 1

	err := unix.Setrlimit(unix.RLIMIT_MEMLOCK, &unix.Rlimit{
		Cur: unix.RLIM_INFINITY,
		Max: unix.RLIM_INFINITY,
	})
	if err != nil {
		fmt.Println("WARNING: Failed to adjust rlimit")
	}

	specs, err := newCounterSpecs()
	if err != nil {
		panic(err)
	}

	objs, err := specs.Load(nil)
	if err != nil {
		panic(err)
	}
	defer objs.Close()

	sock, err := openRawSock(loopback)
	if err != nil {
		panic(err)
	}
	defer syscall.Close(sock)

	if err := syscall.SetsockoptInt(sock, syscall.SOL_SOCKET, SO_ATTACH_BPF, objs.ProgramCountPackets.FD()); err != nil {
		panic(err)
	}

	for range time.Tick(time.Second) {
		var count uint64
		if err := objs.MapPackets.Lookup(uint32(0), &count); err != nil {
			panic(err)
		}

		fmt.Println("Saw", count, "packets")
	}
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
