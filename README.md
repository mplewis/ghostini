# Ghostini

_Ghost + Gemini_

[![Go Reference](https://pkg.go.dev/badge/github.com/mplewis/ghostini.svg)](https://pkg.go.dev/github.com/mplewis/ghostini)
[![Docker Hub](https://flat.badgen.net/badge/icon/Docker%20Hub?icon=docker&label)](https://hub.docker.com/r/mplewis/ghostini)

Turn your [Ghost blog](https://ghost.org/) into a
[Gemini site](https://gemini.circumlunar.space/).

See mine in action: [Ghost blog](https://kesdev.com),
[Gemini site](gemini://kesdev.com).

# Usage

Ghostini uses [Figyr](https://github.com/mplewis/figyr) for configuration, so
you can configure it with environment variables, command-line flags, or a config
file.

```
Options:
    --ghost-site          required        The base URL of your Ghost website
    --content-key         required        The Content API key for your Ghost website
    --domains             required        The domains for which to serve Gemini content
    --gemini-certs-path   optional        The path to the certificates and keys for your Gemini domains
    --host                optional        The host on which to listen
    --port                default: 1965   The port on which to listen
```

- Set `ghost-site` to your Ghost site's base URL, e.g. `https://kesdev.com`.
- Set `content-key` to a [Content Key](https://ghost.org/docs/content-api/)
  configured in your Ghost admin. This grants Ghostini access to your posts.
- Set `domains` to a comma-separated list of domains to handle with TLS, e.g.
  `localhost,kesdev.com,example.com`.
- Put your cert and key files into `gemini-certs-path`. They must be named by
  domain: for example, certs for `kesdev.com` must be named `kesdev.com.crt` and
  `kesdev.com.key`.
- If certificates for a domain are missing, self-signed certificates will be
  implicitly generated and used.

# TODO

- Use zerolog
- Add `DEBUG` logging
- Add linting
