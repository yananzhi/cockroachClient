package cmd

import (
	"flag"
)

type Context struct {
	HTTP string
}

var (
	httpAddr = flag.String("http", "localhost:8080", "TCP network address to bind to for HTTP traffic")
)

// InitFlags sets the server.Context values to flag values.
// Keep in sync with "server/context.go". Values in Context should be
// settable here.
//func InitFlags(ctx *Context) {

//	flag.StringVar(&ctx.HTTP, "http", "127.0.0.1:8080", "host:port to bind for HTTP traffic; 0 to pick unused port")

//}
