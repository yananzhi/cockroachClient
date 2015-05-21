package cmd

import (
	"net/url"

	"github.com/cockroachdb/cockroach/base"
	"github.com/cockroachdb/cockroach/client"
	"github.com/cockroachdb/cockroach/server/cli"
	"github.com/cockroachdb/cockroach/util"
)

// GetHttpKV return a client.KV wrap a httpSender
func GetHttpKV() (*client.KV, error) {

	urlString := cli.Context.RequestScheme() +
		"://root@" + util.EnsureHost(cli.Context.Addr) +
		"?certs=" + cli.Context.Certs

	u, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}
	ctx := &base.Context{}
	ctx.InitDefaults()
	ctx.Insecure = false
	httpsender, err := client.NewHTTPSender(u.Host, ctx)

	kv = client.NewKV(nil, httpsender)
	kv.User = u.User.Username()

	return kv, nil

}
