// Copyright 2014 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License. See the AUTHORS file
// for names of contributors.
//
// Author: Spencer Kimball (spencer.kimball@gmail.com)

package cmd

import (
	"github.com/cockroachdb/cockroach/client"
	"github.com/cockroachdb/cockroach/proto"
)

// A txnSender proxies requests to the underlying KVSender,
// automatically beginning a transaction and then propagating txn
// changes to all commands. On receipt of TransactionRetryError, the
// transaction epoch is incremented and error passed to caller. On
// receipt of TransactionAbortedError, the transaction is re-created
// and error passed to caller.
//
// txnSender is not thread safe.
type TxnSender struct {
	wrapped client.KVSender
	txnEnd  bool // True if EndTransaction was invoked internally
	txn     *proto.Transaction
}

// newTxnSender returns a new instance of txnSender which wraps a
// KVSender and uses the supplied transaction options.
func newTxnSender(wrapped client.KVSender, opts *client.TransactionOptions) *TxnSender {
	return &TxnSender{
		wrapped: wrapped,
		txn: &proto.Transaction{
			Name:      opts.Name,
			Isolation: opts.Isolation,
		},
	}
}

// Send proxies requests to wrapped kv.KVSender instance, taking care
// to attach txn message to each request and update it on each
// response. In the event of a transaction abort, reset txn with a
// minimum priority.
func (ts *TxnSender) Send(call *client.Call) {
	// Send call through wrapped sender.
	call.Args.Header().Txn = ts.txn
	ts.wrapped.Send(call)
	ts.txn.Update(call.Reply.Header().Txn)

}

// Close is a noop for the txnSender. Note that the wrapped sender
// isn't closed because the txnSender is closed immediately upon
// transaction completion and the wrapped sender is reused.
func (ts *TxnSender) Close() {
}
