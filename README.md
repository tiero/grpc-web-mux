# grpc-web-mux
Multiplex gRPC and gRPC Web on the same port, switching on HTTP Content-Type Header.


## Available TCP Listener

* Insecure clear-text TCP listener
* TLS via [CertMagic](https://github.com/caddyserver/certmagic)
* Onion hidden service via [libtor](https://github.com/ipsn/go-libtor) 



## Usage

#### Insecure

```go
package main

import (
	"log"

	mux "github.com/tiero/grpc-web-mux"
	"google.golang.org/grpc"
)

func main() {

  myGrpcServer := grpc.NewServer()
  
  //Register your gRPC handler

	insecureMux, err := mux.NewMuxWithInsecure(
		myGrpcServer,
		mux.InsecureOptions{Address: ":8080"},
	)
	if err != nil {
		log.Panic(err)
	}

	insecureMux.Serve()
}

```

#### TLS

```go
package main

import (
	"log"

	mux "github.com/tiero/grpc-web-mux"
	"google.golang.org/grpc"
)

func main() {

	myGrpcServer := grpc.NewServer()

  //Register your gRPC handler

	tlsMux, err := mux.NewMuxWithTLS(
		myGrpcServer,
		mux.TLSOptions{
			Address: ":9945",
			Domain:  "mydomain.com",
		},
	)
	if err != nil {
		log.Panic(err)
	}

	tlsMux.Serve()
}

```

#### Onion service

```go
package main

import (
	"log"

	mux "github.com/tiero/grpc-web-mux"
	"google.golang.org/grpc"
)

func main() {

  myGrpcServer := grpc.NewServer()

  //Register your gRPC handler

	onionMux, err := mux.NewMuxWithOnion(
		myGrpcServer,
		mux.OnionOptions{
			Port:       80,
			PrivateKey: "myBase64SerializedPrivateKey",
		},
	)
	if err != nil {
		log.Panic(err)
	}

	onionMux.Serve()
}

```

