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

SEC("kretprobe/tcp_sendmsg")
int BPF_KRETPROBE(tcp_sendmsg_exit, int ret) {
    // ttcode
    u8 comm[16], _comm[TARGET_NUM][16] = TARGET_PROC;
    (void)bpf_get_current_comm(&comm, sizeof(comm));
    if (strcmp(comm, _comm[0]) == 1 && strcmp(comm, _comm[1]) == 1 && strcmp(comm, _comm[2]) == 1) {
        return 0;
    }

    u32 pid = bpf_get_current_pid_tgid() >> INT_LEN;
    if (ret < 0) {
        bpf_printk("(tcp_sendmsg_exit) pid: %u errorcode: %d", pid, ret);
        return 0;
    }

    struct tcp_metrics *metrics = bpf_map_lookup_elem(&tcp_link_map, &pid);
    if (!metrics) {
        // ttcode
        bpf_printk("(tcp_sendmsg_exit) %u %s not found tcp_metrics", pid, comm);
        return 0;
    }

    metrics->tx = (u64)ret;
    // __sync_fetch_and_add(&(metrics->tx), (u64)(ret));

    bpf_ringbuf_output(&tcp_output, metrics, sizeof(struct tcp_metrics), 0);

    // ttcode
    bpf_printk("(tcp_sendmsg_exit) pid: %u", pid);
    return 0;
}

SEC("kprobe/tcp_cleanup_rbuf")
int BPF_KPROBE(tcp_cleanup_rbuf, struct sock *sk, int copied) {
    // ttcode
    u8 comm[16], _comm[TARGET_NUM][16] = TARGET_PROC;
    (void)bpf_get_current_comm(&comm, sizeof(comm));
    if (strcmp(comm, _comm[0]) == 1 && strcmp(comm, _comm[1]) == 1 && strcmp(comm, _comm[2]) == 1) {
        return 0;
    }

    u32 pid = bpf_get_current_pid_tgid() >> INT_LEN;
    if (copied < 0) {
        bpf_printk("(tcp_cleanup_rbuf) pid: %u errorcode: %d", pid, copied);
        return 0;
    }

    // ttcode
    bpf_printk("(tcp_cleanup_rbuf) pid: %u", pid);

    struct tcp_metrics *metrics = bpf_map_lookup_elem(&tcp_link_map, &pid);
    // struct tcp_metrics *metrics = bpf_map_lookup_elem(&tcp_link_map, &sk);
    if (!metrics) {
        return 0;
    }

    if (copied <= 0) {
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

    metrics->rx = (u64)copied;
    // __sync_fetch_and_add(&(metrics->rx), (u64)(copied));

    bpf_ringbuf_output(&tcp_output, metrics, sizeof(struct tcp_metrics), 0);
    return 0;
}

// 确认客户端socket
SEC("kprobe/tcp_v4_connect")
int BPF_KPROBE(tcp_v4_connect, struct sock *sk, struct sockaddr *uaddr) {
    // ttcode
    u8 comm[16], _comm[TARGET_NUM][16] = TARGET_PROC;
    (void)bpf_get_current_comm(&comm, sizeof(comm));
    if (strcmp(comm, _comm[0]) == 1 && strcmp(comm, _comm[1]) == 1 && strcmp(comm, _comm[2]) == 1) {
        return 0;
    }

    u32 pid = bpf_get_current_pid_tgid() >> INT_LEN;
    // ttcode
    bpf_printk("(tcp_v4_connect) pid: %u", pid);

    struct tcp_metrics tcpmetrics = {0};
    tcpmetrics.role = LINK_ROLE_CLIENT;
    bpf_map_update_elem(&tcp_link_map, &pid, &tcpmetrics, BPF_ANY);
    return 0;
}