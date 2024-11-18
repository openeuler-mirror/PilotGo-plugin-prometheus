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

SEC("kprobe/tcp_sendmsg")
int BPF_KPROBE(tcp_sendmsg, struct sock *sk, size_t size) {
    u32 pid = bpf_get_current_pid_tgid() >> INT_LEN;
    // ttcode
    u8 comm[16], _comm[TARGET_NUM][16] = TARGET_PROC;
    (void)bpf_get_current_comm(&comm, sizeof(comm));
    if (strcmp(comm, _comm[0]) == 1 && strcmp(comm, _comm[1]) == 1 && strcmp(comm, _comm[2]) == 1) {
        return 0;
    }
    bpf_printk("(tcp_sendmsg) pid: %u", pid);

    struct tcp_metrics *metrics = bpf_map_lookup_elem(&tcp_link_map, &pid);
    // struct tcp_metrics *metrics = bpf_map_lookup_elem(&tcp_link_map, &sk);
    if (!metrics) {
        return 0;
    }

    metrics->family = _(sk->sk_family);
    if (metrics->role == LINK_ROLE_CLIENT) {
        if (metrics->family == AF_INET) {
            metrics->c_ip = _(sk->sk_rcv_saddr);
            metrics->s_ip = _(sk->sk_daddr);

        } else {
            BPF_CORE_READ_INTO(metrics->c_ip6, sk, sk_v6_rcv_saddr);
            BPF_CORE_READ_INTO(metrics->s_ip6, sk, sk_v6_daddr);
        }
        metrics->s_port = bpf_ntohs(_(sk->sk_dport));
        metrics->c_port = _(sk->sk_num);
    } else {
        if (metrics->family == AF_INET) {
            metrics->s_ip = _(sk->sk_rcv_saddr);
            metrics->c_ip = _(sk->sk_daddr);
        } else {
            BPF_CORE_READ_INTO(metrics->s_ip6, sk, sk_v6_rcv_saddr);
            BPF_CORE_READ_INTO(metrics->c_ip6, sk, sk_v6_daddr);
        }
        metrics->s_port = _(sk->sk_num);
        metrics->c_port = bpf_ntohs(_(sk->sk_dport));
    }
    metrics->pid = pid;
    (void)bpf_get_current_comm(metrics->comm, sizeof(metrics->comm));
    return 0;
}