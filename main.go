package main

import (
	"context"
	"fmt"
	"log"
	"time"

	badger "github.com/dgraph-io/badger/v2"
	"github.com/dgraph-io/badger/v2/pb"
)

func main() {
	db, err := badger.Open(badger.DefaultOptions("./p"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()

	handleSubscription := func(kvl *pb.KVList) error {
		for _, kv := range kvl.GetKv() {
			fmt.Println(string(kv.GetKey()), string(kv.GetValue()))
		}
		return nil
	}

	go func() {
		if err := db.Subscribe(ctx, handleSubscription, []byte("f")); err != nil {
			fmt.Println(err)
		}
	}()

	txn := db.NewTransaction(true)
	txn.Set([]byte("foo"), []byte("bar"))
	if err := txn.Commit(); err != nil {
		panic(err)
	}

	time.Sleep(time.Second * 1)
}
