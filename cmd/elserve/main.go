// cmd/elserve はExecution Ledger APIをlocalhostで提供する単体サーバー。
// フロントエンドのdev時（`npm run dev` + このサーバー）や、Wailsを介さずヘッドレスで
// 動かしたい場合に使う。Wailsアプリ本体（ルートのmain.go）も同じ internal/api.NewRouter を
// 使い回すため、挙動はズレない。
//
//	go run ./cmd/elserve -addr :8421 -db execution-ledger.db
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/chankei613/execution-ledger/internal/api"
	"github.com/chankei613/execution-ledger/internal/db"
)

func main() {
	addr := flag.String("addr", ":8421", "待ち受けアドレス")
	dbPath := flag.String("db", "execution-ledger.db", "SQLiteファイル")
	flag.Parse()

	conn, err := db.Init(*dbPath)
	if err != nil {
		log.Fatalf("db init failed: %v", err)
	}

	router := api.NewRouter(conn)
	log.Printf("execution-ledger backend listening on %s", *addr)
	if err := http.ListenAndServe(*addr, router); err != nil {
		log.Fatal(err)
	}
}
