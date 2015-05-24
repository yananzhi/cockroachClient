// session.go
package cmd

import "github.com/cockroachdb/cockroach/client"

// Session is a connect to cockroach server.
type Session struct {
	// can not be nil
	HttpKV *client.KV

	// current txn kv, if nil indacate there has no transacton
	TxnKV *Txn
}

func NewSession() *Session {
	return &Session{
		HttpKV: GetHttpKV(),
	}
}

func (s *Session) runCli() {
	//zyn todo

}
