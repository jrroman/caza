#include "vmlinux.h"

// bpf headers
#include "bpf_endian.h"
#include "bpf_helpers.h"
#include "bpf_tracing.h"

#define AF_INET 2 /* IPv4 address family */

char _license[] SEC("license") = "Apache 2.0";

/**
 * For CO-RE relocatable eBPF programs, __attribute__((preserve_access_index))
 * preserves the offset of the specified fields in the original kernel struct.
 * So here we don't need to include "vmlinux.h". Instead we only need to define
 * the kernel struct and their fields the eBPF program actually requires.
 *
 * Also note that BTF-enabled programs like fentry, fexit, fmod_ret, tp_btf,
 * lsm, etc. declared using the BPF_PROG macro can read kernel memory without
 * needing to call bpf_probe_read*().
 */

/**
 * This struct creates a ring buffer data structure so we can store
 * the tcp data and access it through our main.go program via the ringbuffer.
 * documentation for eBPF maps found here:
 *
 * https://www.kernel.org/doc/html/latest/bpf/maps.html
 */
struct {
        __uint(type, BPF_MAP_TYPE_RINGBUF);
        __uint(max_entries, 1 << 24);
} events SEC(".maps");

/**
 * struct event is the structure that will be submitted to userspace over
 * a ring buffer. It will emit struct event's type info into the ELF produced
 * by bpf2go so we can generate a go type from it.
 *
 */
struct event {
        u16 sport;
        u16 dport;
        u32 saddr;
        u32 daddr;
};
struct event *unused_event __attribute__((unused));

/**
 * Entry into our BPF program occurs here where we trigger the program to fire
 * on tcp_close events. This will allow us to see src_addr:port dst_addr:port as
 * well as the round trip time from src to receiving ack from dst. We are currently
 * disregarding rtt and have the ability to add this functionality later.
 */
SEC("fentry/tcp_close")
int BPF_PROG(tcp_close, struct sock *sk) {
        // If the address is not in IPv4 addr family return
        if (sk->__sk_common.skc_family != AF_INET) {
                return 0;
        }

        struct event *tcp_info;
        tcp_info = bpf_ringbuf_reserve(&events, sizeof(struct event), 0);
        if (!tcp_info) {
                return 0;
        }

        tcp_info->saddr = sk->__sk_common.skc_rcv_saddr;
        tcp_info->daddr = sk->__sk_common.skc_daddr;
        tcp_info->sport = sk->__sk_common.skc_num;
        // ntohs converts u_short from TCP/IP network byte order to host byte order
        tcp_info->dport = bpf_ntohs(sk->__sk_common.skc_dport);

        // submit tcp_info "struct event" to the ring buffer populated with data
        bpf_ringbuf_submit(tcp_info, 0);

        return 0;
};
