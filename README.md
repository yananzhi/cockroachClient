# cockroachClient


##start a session with transaction kv cli
./cli txn --addr=localhost:8080

s [isolation]  [transactionName]
start a txn with given isolation(si or ssi) and transaction name
p [key] [value]
put a key/value
g [key]
get the value with give key
c 
commit the txn
r
rollback the txn
d [key]
delete the give key/value with given key
