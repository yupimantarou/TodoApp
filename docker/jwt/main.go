package main

import (
	"log"
	"net/http"
)

func main() {
	// データベース接続の初期化
	initDB()

	// JWTサーバー用のマルチハンドラ
	mux := http.NewServeMux()
	mux.HandleFunc("/signup/", registerHandler)
	mux.HandleFunc("/auth/", loginHandler)
	mux.HandleFunc("/protected/", protectedHandler)

	// サーバー起動
	log.Println("Starting server on port 3001")
	log.Fatal(http.ListenAndServe(":3001", mux))
}
