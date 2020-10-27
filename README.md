How to ship BPF with your Go project
===

This repository shows you how to use [bpf2go](https://github.com/cilium/ebpf) to embed pre-compiled eBPF in your Go project for easy distribution.

```
$ go generate
$ go build
$ sudo ./ship-bpf-with-go
$ ping localhost # in another window
```

It's the basis of a lightning talk at the [2020 eBPF Summit](https://ebpf.io/summit-2020).
