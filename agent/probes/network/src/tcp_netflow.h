#ifndef __TCP_NETFLOW_H
#define __TCP_NETFLOW_H

#include "vmlinux.h"
#include "bpf_endian.h"
#include "bpf_helpers.h"
#include "bpf_tracing.h"
#include "bpf_core_read.h"

#define TARGET_NUM 3
#define TARGET_PROC {"toagent", "tcpclient", "tcpserver"}

#define LINK_ROLE_SERVER 0
#define LINK_ROLE_CLIENT 1

#define INT_LEN                 32
#define THOUSAND                1000
#define PATH_NUM                20
#define IP_LEN                  4
#define IP_STR_LEN              128
#define IP6_LEN                 16
#define IP6_STR_LEN             128

#define sk_dontcopy_begin       __sk_common.skc_dontcopy_begin
#define sk_dontcopy_end         __sk_common.skc_dontcopy_end
#define sk_hash                 __sk_common.skc_hash
#define sk_portpair             __sk_common.skc_portpair
#define sk_num                  __sk_common.skc_num
#define sk_dport                __sk_common.skc_dport
#define sk_addrpair             __sk_common.skc_addrpair
#define sk_daddr                __sk_common.skc_daddr
#define sk_rcv_saddr            __sk_common.skc_rcv_saddr
#define sk_family               __sk_common.skc_family
#define sk_state                __sk_common.skc_state
#define sk_reuse                __sk_common.skc_reuse
#define sk_reuseport            __sk_common.skc_reuseport
#define sk_ipv6only             __sk_common.skc_ipv6only
#define sk_net_refcnt           __sk_common.skc_net_refcnt
#define sk_bound_dev_if         __sk_common.skc_bound_dev_if
#define sk_bind_node            __sk_common.skc_bind_node
#define sk_prot                 __sk_common.skc_prot
#define sk_net                  __sk_common.skc_net
#define sk_v6_daddr             __sk_common.skc_v6_daddr
#define sk_v6_rcv_saddr __sk_common.skc_v6_rcv_saddr
#define sk_cookie               __sk_common.skc_cookie
#define sk_incoming_cpu         __sk_common.skc_incoming_cpu
#define sk_flags                __sk_common.skc_flags
#define sk_rxhash               __sk_common.skc_rxhash

#ifndef AF_INET
#define AF_INET     2   /* Internet IP Protocol */
#endif
#ifndef AF_INET6
#define AF_INET6    10  /* IP version 6 */
#endif

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