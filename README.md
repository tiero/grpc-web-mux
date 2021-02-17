# grpc-web-mux
Multiplex gRPC and gRPC Web on the same port, switching on HTTP Content-Type Header.


## 🚚 Available Transports

* 🤢 Insecure clear-text 
* 🔐 TLS via [CertMagic](https://github.com/caddyserver/certmagic)
* 🧅 Onion hidden service via [libtor](https://github.com/ipsn/go-libtor) 


## 📩 Install

```sh
$ go get github.com/tiero/grpc-web-mux@latest
```

## ℹ️ Usage

For in-depth documentation refer to [GoDoc examples](https://pkg.go.dev/github.com/tiero/grpc-web-mux/pkg/mux/#pkg-examples)