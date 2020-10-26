// +build ignore

#include <linux/bpf.h>
#include <linux/if_packet.h>

#include <bpf/bpf_helpers.h>

struct
{
	__uint(type, BPF_MAP_TYPE_ARRAY);
	__type(key, __u32);
	__type(value, __u64);
	__uint(max_entries, 1);
} packets SEC(".maps");

SEC("socket")
int count_packets(struct __sk_buff *skb)
{
	if (skb->pkt_type != PACKET_OUTGOING)
		return 0;

	__u32 index = 0;
	__u64 *value = bpf_map_lookup_elem(&packets, &index);
	if (value)
		__sync_fetch_and_add(value, 1);

	return 0;
}

char _license[] SEC("license") = "BSD";
