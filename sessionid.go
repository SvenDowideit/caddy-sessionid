package sessionid

import (
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
type SessionID struct{}

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

	// generate a new sessionid..
	uid := strings.ReplaceAll(uuid.New().String(), "-", "")
	repl.Set("http.session_id", uid)

	c, err := r.Cookie("x-caddy-sessionid")
	if err != nil {
		c = &http.Cookie{
			Name:  "x-caddy-sessionid",
			Value: uid,
			//Domain: "stackdomain..",
			//Path:  "/",
			//Expires: ,
		}
	}
	http.SetCookie(w, c)

	return next.ServeHTTP(w, r)
}

// UnmarshalCaddyfile - this is a no-op
func (m *SessionID) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
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
