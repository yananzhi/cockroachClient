package cmd

import (
	"bufio"
	"log"

	"fmt"

	"os"
	"strings"

	"github.com/cockroachdb/cockroach/base"
	"github.com/cockroachdb/cockroach/client"
	"github.com/cockroachdb/cockroach/proto"

	"github.com/spf13/cobra"
)

// CmdTxn start to run a transactional kv cli interface
var CmdTxn = &cobra.Command{
	Use:   "txn  --addr=localhost:8080",
	Short: "new a txn client",
	Long: `
new a transactional kv client
`,
	Example: `  cli txn`,
	Run:     runTxnKV,
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
	"D": deleteCmd,
	"C": commitCmd,
	"R": rollbackCmd,
}

// global http kv
var kv *client.KV

// transactio kv, each new transaction will update it
var txnkv *Txn

// NewTestBaseContext creates a base context for testing.
// The certs file loader is overriden in individual main_test files.
func NewBaseContext() *base.Context {
	return &base.Context{
		Certs: "/home/zyn/gopath/src/github.com/cockroachdb/cockroach/certs",
	}
}

// startCmd start a transaction, has two parameter,  isolationtype: ssi si , default is si
// transactionName:
func startCmd(c *cmd) error {

	if txnkv != nil {
		fmt.Printf("already in transaction, txn=%v", txnkv)
		return nil
	}

	//default name
	txnName := "testtxn"
	isolation := proto.SNAPSHOT

	fmt.Printf("len c.args = %v\n", len(c.args))
	if len(c.args) >= 1 {
		iso := c.args[0]
		fmt.Println("iso=", iso)
		if strings.EqualFold(iso, "ssi") {
			fmt.Println("ssi")
			isolation = proto.SERIALIZABLE
		} else if strings.EqualFold(iso, "si") {
			fmt.Println("si")
			//do nothing
		} else {
			return fmt.Errorf("error transaction isolation type, must be si or ssi, input is %v", iso)
		}

		if len(c.args) == 2 {
			txnName = c.args[1]
		}

	}
	fmt.Printf("start a transaction, isolation=%v , tansactionName=%v\n", isolation, txnName)

	txnkv = newTxn(kv, &client.TransactionOptions{
		Name:      txnName,
		Isolation: isolation,
	})

	return nil
}

func checkTxnExist() error {
	if txnkv == nil {
		return fmt.Errorf("txn not exist error")
	}
	return nil
}

func genKey(userkey string) proto.Key {
	return proto.Key([]byte(userkey))
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

	if err := txnkv.Run(client.Put(key, value)); err != nil {
		fmt.Printf("put error , error=%v\n", err)
	}

	fmt.Println("put succeed")

	return nil
}

func deleteCmd(c *cmd) error {

	if err := checkTxnExist(); err != nil {
		fmt.Printf("%v\n", err)
		return nil
	}

	if len(c.args) != 1 {
		fmt.Printf("error get args, args=%v\n", c.args)
		return nil
	}

	key := genKey(c.args[0])

	call := client.Delete(key)

	err := txnkv.Run(call)
	if err == nil {

		fmt.Printf("delete key=%v succeed\n", key)
	} else {
		fmt.Printf("get error , error=%v\n", err)
	}

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

	call := client.Get(key)
	gr := call.Reply.(*proto.GetResponse)

	err := txnkv.Run(call)
	if err == nil {

		fmt.Printf("key=%v, value=%v\n", key, resToString(gr))
	} else {
		fmt.Printf("get error , error=%v\n", err)
	}

	return nil
}

// if the key not exist,  res.Value == nil
func resToString(res *proto.GetResponse) string {

	if res.Value != nil {
		if res.Value.Bytes != nil {
			return string(res.Value.Bytes)
		} else {
			fmt.Printf("res.value.bytes == nil\n")
			return "nil"
		}
	} else {

		return "nil"
	}

}

func commitCmd(c *cmd) error {

	if err := checkTxnExist(); err != nil {
		fmt.Printf("%v\n", err)
		return nil
	}

	if len(c.args) != 0 {
		fmt.Printf("argments error, args=%v\n", c.args)
		return nil
	}

	if err, reply := EndTransaction(txnkv, true); err != nil {
		fmt.Printf("commit transaction fail: error=%v", err)
		return nil

	} else {
		if reply.Header().Error != nil {
			fmt.Printf("commit transaction fail: error=%v\n", reply.Header().Error)
			fmt.Println("will auto rollback transaction")

			if err, reply := EndTransaction(txnkv, false); err != nil || reply.Header().Error != nil {
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

// rollback the transaction
func rollbackCmd(c *cmd) error {

	if err := checkTxnExist(); err != nil {
		fmt.Printf("%v\n", err)
		return nil
	}

	if len(c.args) != 0 {
		fmt.Printf("argments error, args=%v\n", c.args)
		return nil
	}

	if err, reply := EndTransaction(txnkv, false); err != nil {
		fmt.Printf("rollback transaction fail: error=%v", err)
		return nil

	} else {
		if reply.Header().Error != nil {
			fmt.Printf("rollback transaction fail: error=%v", reply.Header().Error)

		} else {
			fmt.Println("rollback transaction succeed")
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

func runTxnKV(cmd *cobra.Command, args []string) {
	fmt.Printf("txn kv client:\n")

	if httpKV, err := GetHttpKV(); err == nil {
		kv = httpKV
	} else {
		log.Fatalf("GetHttpKV err, %v", err)
	}

	// for loop read console input command and execute it
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
