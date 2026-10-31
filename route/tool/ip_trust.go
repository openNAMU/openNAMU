package tool

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/netip"
	"strings"

	"github.com/gin-gonic/gin"
)

// Invalid_remote_address is returned by Get_IP when the immediate TCP peer
// (c.Request.RemoteAddr) cannot be parsed into a valid address. It is
// intentionally a non-IP string so that it is never treated as a registered
// identity, matched against a CIDR or a regex ban, or resolved to the loopback
// address. It is also denied by the authorization layer (no "site_view"), so a
// malformed connection is fail-closed rather than given the ::1 identity.
const Invalid_remote_address = "invalid_remote_address"

// trusted_proxies_raw holds the raw trusted_proxies setting as read from the
// other table.
var trusted_proxies_raw string

// trusted_prefixes holds the parsed trusted proxy network list. It is empty
// when no trusted proxy is configured (the default and the safe fallback).
var trusted_prefixes []netip.Prefix

// trusted_proxies_error describes why trusted_proxies could not be parsed.
// When non-empty, forwarding headers are ignored and only the connection peer
// is used (fail-safe).
var trusted_proxies_error string

// Trusted_prefix_count returns the number of parsed trusted proxy networks.
// Zero means no trusted proxy is configured (forwarding headers are ignored).
func Trusted_prefix_count() int {
	return len(trusted_prefixes)
}

// Trusted_proxies_error returns a non-empty string when trusted_proxies
// could not be parsed (fail-safe mode: forwarding headers are ignored).
func Trusted_proxies_error() string {
	return trusted_proxies_error
}

// Validate_trusted_proxies validates a raw trusted_proxies setting string.
// It returns nil when the setting is empty or a valid comma/newline-separated
// list of IP addresses and/or CIDR networks (no catch-all, no malformed token).
// It is used by the settings API to reject unsafe configurations before they
// reach the database.
func Validate_trusted_proxies(raw string) error {
	_, err := parse_trusted_proxies(raw)
	return err
}

// Get_trusted_proxies loads the trusted_proxies setting from the DB and parses
// it into trusted_prefixes. It is called at startup, mirroring Get_ip_select.
// On a parse failure it fails safe: the empty trust list is kept and forwarding
// headers are ignored.
func Get_trusted_proxies(db *sql.DB) {
	raw := ""
	QueryRow_DB(db, `select data from other where name = "trusted_proxies"`, []any{&raw})
	raw = strings.TrimSpace(raw)
	trusted_proxies_raw = raw
	prefixes, err := parse_trusted_proxies(raw)
	if err != nil {
		trusted_prefixes = nil
		trusted_proxies_error = err.Error()
		log.Printf("[IP TRUST] trusted_proxies %q is invalid (%v); forwarding headers will be ignored", raw, err)
		return
	}
	trusted_prefixes = prefixes
	trusted_proxies_error = ""
}

// parse_trusted_proxies parses the raw trusted_proxies setting into a list of
// trusted network prefixes. Tokens are separated by commas and/or newlines (and
// other whitespace). Each token must be a single IP address or a CIDR network.
// Empty tokens are ignored. An empty or whitespace-only setting yields an empty
// list (no trusted proxy) rather than a catch-all.
func parse_trusted_proxies(raw string) ([]netip.Prefix, error) {
	tokens := strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == '\n' || r == '\r' || r == ' ' || r == '\t'
	})
	prefixes := make([]netip.Prefix, 0, len(tokens))
	for _, tok := range tokens {
		p, err := parse_trusted_token(tok)
		if err != nil {
			return nil, err
		}
		prefixes = append(prefixes, p)
	}
	return prefixes, nil
}

// parse_trusted_token parses a single token as an IP address or CIDR network.
// It rejects catch-all networks (0.0.0.0/0 and ::/0) and malformed tokens.
func parse_trusted_token(tok string) (netip.Prefix, error) {
	s := strings.TrimSpace(tok)
	if s == "" {
		return netip.Prefix{}, fmt.Errorf("empty token")
	}
	p, err := netip.ParsePrefix(s)
	if err == nil {
		// Reject the catch-all universe. Only 0.0.0.0/0 and ::/0 can be /0.
		if p.Bits() == 0 {
			return netip.Prefix{}, fmt.Errorf("%q is a catch-all network (0.0.0.0/0 or ::/0) and is rejected", s)
		}
		return p, nil
	}
	// A bare IP address (no CIDR) is acceptable; treat it as a single-host prefix.
	addr, aerr := netip.ParseAddr(s)
	if aerr != nil {
		return netip.Prefix{}, fmt.Errorf("invalid address or network %q: %w", s, aerr)
	}
	return netip.PrefixFrom(addr, addr.BitLen()), nil
}

// parse_remote_addr parses a TCP peer string ("host:port") into a netip.Addr.
// It returns (zero, false) when the peer is missing or unparseable.
func parse_remote_addr(remote string) (netip.Addr, bool) {
	if remote == "" {
		return netip.Addr{}, false
	}
	addr, err := netip.ParseAddr(remote)
	if err == nil {
		return addr, true
	}
	// netip.ParseAddr does not accept a port. Tolerate "host:port" forms.
	host, _, herr := net.SplitHostPort(remote)
	if herr == nil {
		addr, err = netip.ParseAddr(host)
		if err == nil {
			return addr, true
		}
	}
	return netip.Addr{}, false
}

// parse_single_ip parses a raw header value that must denote exactly one IP
// address (optionally in "host:port" form). It rejects comma-separated,
// multi-token, or malformed values so that an untrusted client cannot select
// which entry of a list is used.
func parse_single_ip(raw string) (netip.Addr, bool) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return netip.Addr{}, false
	}
	// Reject list-like values before parsing.
	if strings.ContainsAny(s, ",;") {
		return netip.Addr{}, false
	}
	addr, err := netip.ParseAddr(s)
	if err == nil {
		return addr, true
	}
	// netip.ParseAddr does not accept a port. Tolerate "host:port" forms.
	host, _, herr := net.SplitHostPort(s)
	if herr == nil {
		addr, err = netip.ParseAddr(host)
		if err == nil {
			return addr, true
		}
	}
	return netip.Addr{}, false
}

// read_single_ip_header reads a single forwarding header and returns the address
// it denotes, or false when the header is missing, ambiguous, or malformed.
func read_single_ip_header(c *gin.Context, name string) (netip.Addr, bool) {
	return parse_single_ip(c.GetHeader(name))
}

// is_trusted_peer reports whether peer is contained in the configured trusted
// proxy list. A malformed or empty configuration (or a parse error) yields
// false, so forwarding headers are ignored and only the connection peer is used.
func is_trusted_peer(peer netip.Addr) bool {
	if trusted_proxies_error != "" || len(trusted_prefixes) == 0 {
		return false
	}
	return peer_in_prefixes(peer, trusted_prefixes)
}

// peer_in_prefixes reports whether peer is contained in any prefix. IPv4-mapped
// IPv6 peers (e.g. ::ffff:a.b.c.d) are also matched against IPv4 prefixes by
// considering the unmapped IPv4 view.
func peer_in_prefixes(peer netip.Addr, prefixes []netip.Prefix) bool {
	if peer.Is4In6() {
		if v4 := peer.Unmap(); v4.IsValid() {
			for _, p := range prefixes {
				if p.Contains(v4) {
					return true
				}
			}
		}
	}
	for _, p := range prefixes {
		if p.Contains(peer) {
			return true
		}
	}
	return false
}

// selected_header_names returns the header names to consult for the given
// load_ip_select mode, in order. For "REMOTE_ADDR" it returns none; for "default"
// or "" it returns the default forwarding headers; otherwise the explicit header
// (with the HTTP_ prefix stripped and underscores turned into hyphens).
func selected_header_names(setting string) []string {
	s := strings.TrimSpace(setting)
	if s == "" || strings.EqualFold(s, "default") {
		return []string{"X-Real-IP", "CF-Connecting-IP"}
	}
	if strings.EqualFold(s, "REMOTE_ADDR") {
		return []string{}
	}
	if len(s) > len("HTTP_") && strings.EqualFold(s[:len("HTTP_")], "HTTP_") {
		s = s[len("HTTP_"):]
		s = strings.ReplaceAll(s, "_", "-")
	}
	return []string{s}
}
