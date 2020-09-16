# Caddy SessionID

> Caddy v2 Module that sets a session ID.

The Session ID would ideally be random, but to enable multiple edge requests, I think its easier to use a hash that is re-creatable, based on remoteIP, ...

## Usage

### With a Caddyfile

Usage with a Caddyfile is fairly straightforward. Simply add the `session_id` directive to a handler block, and the `{http.session_id}` template will be set.

Will also add a Cookie named `x-caddy-sessionid` for auto-session tracking.

If you wish to use the directive in a top level block, you must explicitly define the order.

```
{
  order session_id before header
}
```

### With JSON Config

In the JSON Config, all you need to do is add the `session_id` hander to your `handle[]` chain. The same template as documented above will be set. Note that you must set this before any handlers that you want to use the template in.

```json
{
  "handler": "session_id"
}
```

### Example

The following example Caddyfile sets the `x-session-id` header for all responses.

```
{
  order session_id before header
}

localhost {
  session_id <optional cookie domain>

  header * x-session-id "{http.session_id}"
  respond * "{http.session_id}" 200
}
```
