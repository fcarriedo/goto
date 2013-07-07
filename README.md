# goto

Service that lets you create short mapping for easy sharing a la Google's
internal `http://go/somewhere`.

### Usage

Ideally when using it in production, you would like to map the the internal DNS
service of `go` to this service and make it listen (or proxy it) on port `80`.

In development, you might want to create an entry in your `hosts` file for
`127.0.0.1      go`.

For maximizing succinctness it runs on port `80` by default but can be changed
at startup with the `-p` param.

```
  $ goto -p 8080 # starts the service in port 80
```

### Reserved endpoints

Supposing that you have the mapping on your DNS or `hosts` as `go => [servers
IP address]`:

`http://go/mappings` - Will give you a list of the current existing mappings.

`http://go/mappings/{entry}` - Will give you the url that maps the given entry
if exists or *404 Not Found*.

### Roadmap

The goal of this project is to stay as simple and hackable as possible. Features
need to be very, **very** justified.

#### What it will NOT

  * On the few pages that it entails, it will **never** support crippled
    irresponsible browsers (yeah, yeah... IE).
  * Filter unwanted (possible malignus) URLs. Since the final destination gets
    obscured by the mapped URL, it expects the users to be responsible to where
    they redirect.

#### What it WILL

  * Very simple csv datastore

#### MAYBE

  * Simple statistics (maybe achiebable by integrating google analytics - will
    it work on an *only* internal network?)
  * Authenticated creation of mappings? - Since we will not be filtering
    unwanted URL, at least we might want to know who registered them (will this
    prevent abuse?).
    * Possible shortcommings: Support LDAP? Support DB Auth? What DB? - We can
      see it starts to get spooky. Maybe no support at all and count on mature
      responsible users.

### Customizing UI

The templates for the few existing pages live under `web/tmpl`. You can modify
them as needed.

### License

MIT
