#include "tcp_netflow.h"

char g_linsence[] SEC("license") = "GPL";

struct {
    __uint(type, BPF_MAP_TYPE_RINGBUF);
    __uint(max_entries, 8 * 4096);
} tcp_output SEC(".maps");

#define __TCP_LINK_MAX (10 * 1024)

struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(key_size, sizeof(u32));
    __uint(value_size, sizeof(struct tcp_metrics));
    __uint(max_entries, 20 * 1024);
} tcp_link_map SEC(".maps");
