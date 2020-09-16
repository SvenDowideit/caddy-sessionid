package sessionid

import (
	"log"
	"net/http"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/google/uuid"
)

func init() {
	caddy.RegisterModule(SessionID{})
	httpcaddyfile.RegisterHandlerDirective("session_id", parseCaddyfile)
}

// SessionID implements an HTTP handler that writes a
// unique request ID to response headers.
type SessionID struct {
	CookieDomain string
}

// CaddyModule returns the Caddy module information.
func (SessionID) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.session_id",
		New: func() caddy.Module { return new(SessionID) },
	}
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (m SessionID) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	repl := r.Context().Value(caddy.ReplacerCtxKey).(*caddy.Replacer)

	c, err := r.Cookie("x-caddy-sessionid")
	if err != nil {
		// Get the domain to use from the session_id setting in the Caddyfile
		CookieDomain := m.CookieDomain
		log.Printf("[INFO] cookiedomain: %s, request: %s\n", CookieDomain, r.Host)
		if !strings.HasSuffix(r.Host, CookieDomain) {
			CookieDomain = r.Host
		}

		// generate a new sessionid..
		uid := strings.ReplaceAll(uuid.New().String(), "-", "")
		c = &http.Cookie{
			Name:   "x-caddy-sessionid",
			Value:  uid,
			Domain: CookieDomain, // Need to figure out how to share the same cookie, or to generate it so it can be used for a lookup
			Path:   "/",
			//Expires: ,
		}
	}
	http.SetCookie(w, c)
	repl.Set("http.session_id", c.Value)

	return next.ServeHTTP(w, r)
}

// UnmarshalCaddyfile - session_id <cookiedomain>
func (m *SessionID) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if !d.Args(&m.CookieDomain) {
			// not enough args
			return d.ArgErr()
		}
		if d.NextArg() {
			// too many args
			return d.ArgErr()
		}
	}
	log.Printf("[INFO] unmarshal session_id: %s\n", m.CookieDomain)

	return nil
}

func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	m := new(SessionID)
	err := m.UnmarshalCaddyfile(h.Dispenser)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// Interface guards
var (
	_ caddyhttp.MiddlewareHandler = (*SessionID)(nil)
	_ caddyfile.Unmarshaler       = (*SessionID)(nil)
)
