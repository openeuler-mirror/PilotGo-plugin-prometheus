#ifndef __TCP_NETFLOW_H
#define __TCP_NETFLOW_H

#include "vmlinux.h"
#include "bpf_endian.h"
#include "bpf_helpers.h"
#include "bpf_tracing.h"
#include "bpf_core_read.h"

#define INT_LEN                 32
#define THOUSAND                1000
#define PATH_NUM                20
#define IP_LEN                  4
#define IP_STR_LEN              128
#define IP6_LEN                 16
#define IP6_STR_LEN             128

struct tcp_metrics {
    __u32 pid;     // process id
    union {
        __u32 c_ip;
        unsigned char c_ip6[IP6_LEN];
    };
    union {
        __u32 s_ip;
        unsigned char s_ip6[IP6_LEN];
    };
    __u16 s_port;   // server port
    __u16 c_port;   // client port
    __u16 family;
    __u16 role;     // role: client:1/server:0
    u8 comm[TASK_COMM_LEN];

    __u16 opt_family;
    union {
        __u32 opt_c_ip;
        unsigned char opt_c_ip6[IP6_LEN];
    };

    __u64 rx;               // FROM tcp_cleanup_rbuf
    __u64 tx;               // FROM tcp_sendmsg
};
const struct tcp_metrics *unused __attribute__((unused));

#endif /* __TCP_NETFLOW_H */