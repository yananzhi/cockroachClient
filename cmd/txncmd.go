package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/cockroachdb/cockroach/client"
	"github.com/cockroachdb/cockroach/proto"
	"github.com/cockroachdb/cockroach/rpc"
	"github.com/cockroachdb/cockroach/storage/engine"

	"code.google.com/p/go-commander"
)

// A CmdInit command initializes a new Cockroach cluster.
var CmdTxn = &commander.Command{
	UsageLine: "txn",
	Short:     "init a transactional cli",
	Long: `use a transactional kv client
`,
	Run:  runTxnKV,
	Flag: *flag.CommandLine,
}

type cmd struct {
	name string
	args []string
}

// cmdDict maps from command name to function implementing the command.
// Use only upper case letters for commands. More than one letter is OK.
var cmdDict = map[string]func(c *cmd) error{
	"S": startCmd,
	"P": putCmd,
	"G": getCmd,
	"C": commitCmd,
}

//current httpsender
var httpsender client.HTTPSender

var kv *client.KV

var txnkv *client.KV

var txnsender *TxnSender

func InitHttpSender(addr string) *client.HTTPSender {
	return client.NewHTTPSender(addr, &http.Transport{
		TLSClientConfig: rpc.LoadInsecureTLSConfig().Config(),
	})
}

func startCmd(c *cmd) error {
	//for test
	//	fmt.Printf("startcmd: %v\n", c)

	fmt.Println("start a transaction")

	if txnkv != nil {
		fmt.Printf("already in transaction, txn=%v", txnkv.Sender())
		return nil
	}

	txnsender = newTxnSender(kv.Sender(), &client.TransactionOptions{
		Name:      "kv txn",
		Isolation: proto.SERIALIZABLE, //todo use input argment to set the isolation level
	})

	//	txnkv = client.NewKV(txnsender, nil)
	txnkv = client.NewKV(nil, txnsender)
	txnkv.User = kv.User
	txnkv.UserPriority = kv.UserPriority

	return nil
}

func checkTxnExist() error {
	if txnkv == nil {
		return fmt.Errorf("txn not exist error")
	}
	return nil
}

func genKey(userkey string) proto.Key {
	return engine.MakeKey(proto.Key([]byte("~")), proto.Key([]byte(userkey)))
}

func putCmd(c *cmd) error {

	if err := checkTxnExist(); err != nil {
		fmt.Printf("%v", err)
		return nil
	}

	if len(c.args) != 2 {
		fmt.Printf("error put args, args=%v\n", c.args)
		return nil
	}

	key := genKey(c.args[0])

	value := []byte(c.args[1])
	if err := txnkv.Put(key, value); err != nil {
		fmt.Printf("put error , error=%v\n", err)
	}

	fmt.Println("put succeed")

	return nil
}

func getCmd(c *cmd) error {

	if err := checkTxnExist(); err != nil {
		fmt.Printf("%v\n", err)
		return nil
	}

	if len(c.args) != 1 {
		fmt.Printf("error get args, args=%v\n", c.args)
		return nil
	}

	key := genKey(c.args[0])

	if found, value, _, err := txnkv.Get(key); err != nil && found {
		if !found {
			fmt.Printf("given key is not found\n")
			return nil
		}
		fmt.Printf("get error , error=%v\n", err)
	} else {
		fmt.Printf("key=%v, value=%v\n", key, string(value))
	}

	return nil
}

func endtransaction(txnkv *client.KV, commit bool) (error, reply proto.Response) {
	etArgs := &proto.EndTransactionRequest{Commit: commit}
	etReply := &proto.EndTransactionResponse{}
	if err := txnkv.Call(proto.EndTransaction, etArgs, etReply); err != nil {
		fmt.Printf("commit transaction reuqest error:%v", err)
	}
	return nil, etReply

}

func commitCmd(c *cmd) error {
	//for test
	fmt.Printf("commit command: %v\n", c)

	if err := checkTxnExist(); err != nil {
		fmt.Printf("%v\n", err)
		return nil
	}

	if len(c.args) != 0 {
		fmt.Printf("argments error, args=%v\n", c.args)
		return nil
	}

	if err, reply := endtransaction(txnkv, true); err != nil {
		fmt.Printf("commit transaction fail: error=%v", err)
		return nil

	} else {
		if reply.Header().Error != nil {
			fmt.Printf("commit transaction fail: error=%v", reply.Header().Error)
			fmt.Println("will auto rollback transaction")

			if err, reply := endtransaction(txnkv, false); err != nil || reply.Header().Error != nil {
				fmt.Printf("rollback transaction fail: error=%v, reply.error=%v", err, reply.Header().Error)
				return nil
			}

			fmt.Println("rollback transaction succeed")

		} else {
			fmt.Println("commit transaction succeed")
		}
	}

	txnkv = nil

	return nil
}

func initCmd(str string) (*cmd, error) {
	args := strings.Split(str, " ")

	if len(args) < 1 {
		return nil, fmt.Errorf("not enouf arguments")
	}

	c := &cmd{
		name: strings.ToUpper(args[0]),
		args: args[1:],
	}

	return c, nil
}

func runTxnKV(cmd *commander.Command, args []string) {
	fmt.Printf("txn kv client:\n")

	context := &Context{}
	InitFlags(context)
	httpsender := InitHttpSender(context.HTTP)
	//	kv = client.NewKV(httpsender, nil)
	kv = client.NewKV(nil, httpsender)
	kv.User = "root"
	kv.UserPriority = -1

	for {
		reader := bufio.NewReader(os.Stdin)
		strBytes, _, err := reader.ReadLine()

		if err == nil {
			str := string(strBytes)

			if str == "" {
				continue
			}

			if c, err := initCmd(string(strBytes)); err == nil {
				if fn, exist := cmdDict[c.name]; exist {
					fn(c)
				} else {
					fmt.Printf("cmd %v not exist\n", c)
				}

			} else {
				fmt.Printf("getcmd error:%v", err)
			}
		}

	}

}
