package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

// Todo のデータ構造
type Todo struct {
	ID   int
	Task string
	Done bool
}

var (
	tmpl = template.Must(template.ParseFiles("index.html"))
	db   *sql.DB
	mu   sync.Mutex
)

// initDB: DB 接続の初期化とテーブル作成
func initDB() {
	// 環境変数から DB 接続情報を取得
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// PostgreSQL 用の接続文字列を組み立てる
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("DB接続エラー:", err)
	}

	// DB への接続確認
	if err = db.Ping(); err != nil {
		log.Fatal("Pingエラー:", err)
	}

	// Todo テーブルの作成（存在しなければ）
	sqlStmt := `CREATE TABLE IF NOT EXISTS todos (
		id SERIAL PRIMARY KEY,
		task TEXT NOT NULL,
		done BOOLEAN NOT NULL DEFAULT FALSE
	);`
	if _, err = db.Exec(sqlStmt); err != nil {
		log.Fatal("テーブル作成エラー:", err)
	}
}

func main() {
	// DB の初期化
	initDB()
	defer db.Close()

	// ルーティング設定
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/add", addHandler)

	// サーバー起動
	log.Println("Server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

// indexHandler: DB から Todo を取得して一覧表示
func indexHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	rows, err := db.Query("SELECT id, task, done FROM todos")
	if err != nil {
		http.Error(w, "クエリエラー: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var t Todo
		if err := rows.Scan(&t.ID, &t.Task, &t.Done); err != nil {
			http.Error(w, "行スキャンエラー: "+err.Error(), http.StatusInternalServerError)
			return
		}
		todos = append(todos, t)
	}

	if err := tmpl.Execute(w, todos); err != nil {
		http.Error(w, "テンプレート実行エラー: "+err.Error(), http.StatusInternalServerError)
	}
}

// addHandler: フォームから送信された Todo を DB に追加
func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	task := r.FormValue("task")
	if task == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	_, err := db.Exec("INSERT INTO todos (task, done) VALUES ($1, $2)", task, false)
	if err != nil {
		http.Error(w, "挿入エラー: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
