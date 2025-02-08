package limiter

import (
	"net"
	"net/http"
	"strings"
)

func (limiter *Limiter) GetToken(r *http.Request) string {
	return r.Header.Get("API_KEY")
}

func (limiter *Limiter) GetIP(r *http.Request) net.IP {
	return GetIP(r)
}

func (limiter *Limiter) CheckIfKeyIsIPAddress(ip string) bool {
	if net.ParseIP(ip) == nil {
		return false
	} else {
		return true
	}
}

func GetIP(r *http.Request) net.IP {
	remoteAddr := strings.TrimSpace(r.RemoteAddr)
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return net.ParseIP(remoteAddr)
	}

	return net.ParseIP(host)
}
