# goto

Service that lets you create short mapping for easy sharing a la Google's
internal `http://go/somewhere`.

### Usage

Ideally when using it in production, you would like to map the the internal DNS
service of `go` to this service and make it listen (or proxy it) on port `80`.

In development, you might want to create an entry in your `hosts` file for
`127.0.0.1  go`.

For maximizing succinctness it runs on port `80` by default but can be changed
at startup with the `-p` param.

```
  goto -p 8080 # starts the service in port 80
```

### Customizing UI

The templates for the few existing pages live under `web/tmpl`. You can modify
them as needed.

### License

MIT
