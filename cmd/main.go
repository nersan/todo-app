package main

import (
	"html/template"
	"log"
	"net/http"
	"sync"
)

// Todoのデータ構造
type Todo struct {
	ID   int
	Task string
	Done bool
}

// グローバル変数：Todoリスト、次のID、排他制御用mutex
var (
	todos  = []Todo{}
	nextID = 1
	mu     sync.Mutex
)

// HTMLテンプレートの読み込み（エラーがあればプログラム起動時に終了）
var tmpl = template.Must(template.ParseFiles("./index.html"))

func main() {
	// ルーティング設定
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/add", addHandler)

	// サーバー起動（ポート8080）
	log.Println("Server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

// indexHandler: Todoリストの一覧を表示するハンドラ
func indexHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	// テンプレートにTodoリストを渡してレンダリング
	if err := tmpl.Execute(w, todos); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// addHandler: 新しいTodoを追加するハンドラ
func addHandler(w http.ResponseWriter, r *http.Request) {
	// POSTメソッド以外は一覧ページにリダイレクト
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// フォームからtaskの値を取得
	task := r.FormValue("task")
	if task != "" {
		mu.Lock()
		// 新しいTodoを追加し、IDをインクリメント
		todos = append(todos, Todo{
			ID:   nextID,
			Task: task,
			Done: false,
		})
		nextID++
		mu.Unlock()
	}
	// 登録後は一覧ページにリダイレクト
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
