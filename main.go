package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

var db *sql.DB

type Todo struct {
	ID   int    `json:"id"`
	Task string `json:"task"`
}

func main() {
	var err error
	db, err = initDB()
	if err != nil {
			log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/", listTodos)
	http.HandleFunc("/add", addTodo)
	http.HandleFunc("/delete", deleteTodo)

	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}

func listTodos(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, task FROM todos")
	if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
			var todo Todo
			if err := rows.Scan(&todo.ID, &todo.Task); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
			}
			todos = append(todos, todo)
	}

	tmpl := `
	<h1>Todo List</h1>
	<form action="/add" method="POST">
			<input type="text" name="task" required>
			<button type="submit">Add</button>
	</form>
	<ul>
			{{range .}}
			<li>{{.Task}} <a href="/delete?id={{.ID}}">Delete</a></li>
			{{end}}
	</ul>
	`
	t := template.Must(template.New("todos").Parse(tmpl))
	t.Execute(w, todos)
}

func addTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
			task := r.FormValue("task")
			_, err := db.Exec("INSERT INTO todos (task) VALUES (?)", task)
			if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	_, err := db.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}