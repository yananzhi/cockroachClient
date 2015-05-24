/*
simulate a tpcc test, it should be ssi isolation,
there has two tables: table1, table2
table1 if for generator a unique orderID
table1
keyOrder　orderID


table2 save every order
table2
orderid   monery

there is a transaction to new a order
transaction1:
orderID = get(keyOrder)
if orderid = nil{
	// init the order id
	put(keyOrder, 1)
	commit transaction
	return
}

if get(orderID) != nil{
	fatal("txn error, it should be nil")
}

put(orderID, xxMoney)
orderID++
delete(keyOrder)
put(keyOrder, orderID)
commit transaction



there is another transaction to read the table1
get(keyOrder)

the transacton1 and transaction2 is execute parallel,
*/

package cmd

import (
	"fmt"
	"sync"
	"time"

	"github.com/cockroachdb/cockroach/client"
	"github.com/cockroachdb/cockroach/proto"
	"github.com/spf13/cobra"
)

// CmdTxn start to run a transactional kv cli interface
var CmdTpcc = &cobra.Command{
	Use:   "tpcc  --addr=localhost:8080",
	Short: "run a simple tpcc testt",
	Long: `
new a transactional kv client
`,
	Example: `  cli tpcc`,
	Run:     runTpcc,
}

const (
	orderIDKey = 100

	// it must great than orderIDKey
	initOrderID = 200
	// go routine number for run txn1
	txn1Parellel = 2

	// go routine number for run txn2
	txn2Parellel = 0
)

// txnNumAll generate unique txnNumber
type txnNumAll struct {
	sync.RWMutex
	num int64
}

func (t *txnNumAll) allcateNum() int64 {
	t.Lock()
	defer t.Unlock()
	t.num++
	return t.num
}

func runTpcc(cmd *cobra.Command, args []string) {
	fmt.Println("run Tpcc ")

	kv := GetHttpKV()

	// Init the orderID
	if err := PutInt64(orderIDKey, initOrderID, kv); err != nil {
		panic("init the orderID error")
	}

	allocate := &txnNumAll{
		num: 200,
	}

	for i := 0; i < txn1Parellel; i++ {
		go txn1Task(kv, int64(i+1), allocate)
	}

	for i := 0; i < txn2Parellel; i++ {
		go txn2Task(kv)
	}

	//	transaction1(kv)

	<-time.After(time.Second * 60)

}

//put(orderID, xxMoney)
//orderID++
//delete(keyOrder)
//put(keyOrder, orderID)
//commit transaction
func transaction1(kv *client.KV, routineID int64, txnID int64) {
	txn := newTxn(kv, &client.TransactionOptions{
		Name: "transaction1" + fmt.Sprintf(":%v.%v", routineID, txnID),
		//		Isolation: proto.SNAPSHOT, // use snapshot isolation
		Isolation: proto.SERIALIZABLE, // use SERIALIZABLE isolation to debug
	})
	commit := false

	// it encounter fail , rolllback the transaction
	defer func() {
		if !commit {
			// there may w-w conflic
			fmt.Printf("rollback transaction1\n")

			EndTransaction(txn, false)
		}
	}()

	// get the current orderID
	orderID, exist, err := GetInit64(orderIDKey, txn)
	if err != nil {
		return
	} else if !exist {
		// the orderID must be inited before the txn1 start
		panic("orderid not exist")
		//		if err = PutInt64(orderIDKey, initOrderID, txn); err == nil {
		//			fmt.Printf("init the order succeed\n")
		//			commit = true
		//			//commit the transaction
		//			if err, _ = EndTransaction(txn, true); err != nil {
		//				fmt.Printf("commit transaction err: %v\n", err)
		//				return
		//			}
		//			fmt.Printf("commit transaction1 succeed\n")
		//			return
		//		}

		//shoud not panic ,may a push transaciton error
		//		fmt.Printf("err=%v", err)
		//		panic("can init the order id")

	}

	//check the key orderID is not exist
	if _, exist, err := GetInit64(orderID, txn); err != nil {
		return
	} else if exist {
		fmt.Printf("routineID=%v ,  txnNumber=%v", routineID, txnID)
		fmt.Printf("orderid=%v\n", orderID)
		panic("it should not exist")
	}

	if err := PutInt64(orderID, 6000, txn); err != nil {
		fmt.Printf("put order id =%v err, err=%v\n", orderID, err)
		return
	}

	// update the orderID++
	orderID++
	//	if err = DeleteInt64(orderIDKey, txn); err != nil {
	//		fmt.Printf("delete orderIDKey err: %v\n", err)
	//		return

	//	}
	if err = PutInt64(orderIDKey, orderID, txn); err != nil {
		fmt.Printf("put orderIDKey err: %v\n", err)
		return
	}

	//commit the transaction
	if err, _ = EndTransaction(txn, true); err != nil {
		fmt.Printf("commit transaction err: %v\n", err)
		return
	}
	fmt.Printf("commit transaction1 succeed\n")
	commit = true

}

func transaction2(kv *client.KV) {
	txn := newTxn(kv, &client.TransactionOptions{
		Name:      "transaction2",
		Isolation: proto.SNAPSHOT, // use snapshot isolation
	})

	GetInit64(orderIDKey, txn)
	GetInit64(orderIDKey, txn)
	GetInit64(orderIDKey, txn)

}

func txn1Task(kv *client.KV, routineID int64, numAll *txnNumAll) {
	for {
		transaction1(kv, routineID, numAll.allcateNum())
	}

}

func txn2Task(kv *client.KV) {
	for {
		transaction2(kv)
	}
}
