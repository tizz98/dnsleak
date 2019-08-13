package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"

	"github.com/caddyserver/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
)

var urlRe = regexp.MustCompile(`^(?P<num>\d)\.(?P<id>\w+)\.dnsleak\.app\.$`)

// offset by one
const (
	numMatchIdx = 1 + iota
	idMatchIdx
)

type urlMatch struct {
	num int
	id  string
}

func (m *urlMatch) RedisKey() string {
	return fmt.Sprintf("%d.%s", m.num, m.id)
}

func extractFromUrl(url string) *urlMatch {
	result := urlRe.FindStringSubmatch(url)
	if len(result) == 0 {
		return nil
	}

	num, err := strconv.Atoi(result[numMatchIdx])
	if err != nil {
		panic("unable to convert to int: " + err.Error())
	}

	return &urlMatch{
		num: num,
		id:  result[idMatchIdx],
	}
}

func init() {
	caddy.RegisterPlugin("leak", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

type LeakHandler struct {
	Next plugin.Handler

	db *DB
}

func ParseIP(s string) (string, error) {
	ip, _, err := net.SplitHostPort(s)
	if err == nil {
		return ip, nil
	}

	ip2 := net.ParseIP(s)
	if ip2 == nil {
		return "", errors.New("invalid IP")
	}

	return ip2.String(), nil
}

func MustParseIP(ip string) string {
	parsed, err := ParseIP(ip)
	if err != nil {
		panic("unable to parse ip: " + err.Error())
	}
	return parsed
}

func (h LeakHandler) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	for _, q := range r.Question {
		extracted := extractFromUrl(q.Name)
		h.db.UpdateForUrlMatch(extracted, MustParseIP(w.RemoteAddr().String()))
	}

	return plugin.NextOrFailure(h.Name(), h.Next, ctx, w, r)
}

func (h LeakHandler) Name() string { return "leak" }

func setup(c *caddy.Controller) error {
	c.Next() // 'leak'

	if c.NextArg() {
		return plugin.Error("leak", c.ArgErr())
	}

	db := NewDB(os.Getenv("DNS_LEAK_REDIS_URI"))

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return LeakHandler{Next: next, db: db}
	})

	return nil
}
