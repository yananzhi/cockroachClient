package cmd

import (
	"net/url"
	"strconv"

	"github.com/cockroachdb/cockroach/base"
	"github.com/cockroachdb/cockroach/client"
	"github.com/cockroachdb/cockroach/proto"
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

// Runnable support Txn and client.KV, they both import Run function
type Runnable interface {
	Run(calls ...client.Call) (err error)
}

const (
	Base    = 10
	BitSize = 64
)

// PutInt64 put a int64 key value to cockroach server,
// return nil indicate put succeed and vice versa
func PutInt64(key, value int64, kv Runnable) error {
	return kv.Run(client.Put(proto.Key([]byte(strconv.FormatInt(key, Base))), []byte(strconv.FormatInt(value, Base))))
}

// GetInit64 get a int64 key, value is valid when the exist is true
func GetInit64(key int64, kv Runnable) (value int64, exist bool, err error) {
	call := client.Get(proto.Key([]byte(strconv.FormatInt(key, Base))))
	gr := call.Reply.(*proto.GetResponse)

	if err = kv.Run(call); err == nil {
		v := gr.Value
		if v != nil {
			value, err = strconv.ParseInt(string([]byte(v.Bytes)), Base, BitSize)
			return value, true, err
		} else {
			return 0, false, nil
		}
	} else {
		return 0, false, err
	}
}

// DeleteInt64 delete a int64 key, is a key is not exist, delete is return with no error.
func DeleteInt64(key int64, kv Runnable) error {
	return kv.Run(client.Delete(proto.Key(proto.Key([]byte(strconv.FormatInt(Base, Base))))))
}

// EndTransaction commit or rollback a transaction.
func EndTransaction(txnkv *Txn, commit bool) (error, proto.Response) {
	etArgs := &proto.EndTransactionRequest{Commit: commit}
	etReply := &proto.EndTransactionResponse{}

	call := client.Call{Args: etArgs, Reply: etReply}
	if err := txnkv.Run(call); err != nil {
		//		fmt.Printf("commit transaction reuqest error:%v", err)
		return err, nil
	}
	return nil, etReply

}
