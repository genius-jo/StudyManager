package model

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type todoSqlHandler struct {
	db *sql.DB
}

func (tsh *todoSqlHandler) GetTodos(sessionId string) []*Todo {
	todos := []*Todo{}
	rows, err := tsh.db.Query("SELECT id, name, completed, createdAt FROM todos WHERE sessionId=?", sessionId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var todo Todo
		rows.Scan(&todo.Id, &todo.Name, &todo.Completed, &todo.Created)
		todos = append(todos, &todo)
	}
	return todos
}

func (tsh *todoSqlHandler) AddTodo(name string, sessionId string) *Todo {
	stmt, err := tsh.db.Prepare("INSERT INTO todos (sessionId, name, completed, createdAt) VALUES (?, ?, ?, datetime('now'))")
	if err != nil {
		panic(err)
	}
	rst, err := stmt.Exec(sessionId, name, false)
	if err != nil {
		panic(err)
	}
	id, _ := rst.LastInsertId()

	var todo Todo
	todo.Id = int(id)
	todo.Name = name
	todo.Completed = false
	todo.Created = time.Now()
	return &todo
}

func (tsh *todoSqlHandler) RemoveTodo(id int) bool {
	stmt, err := tsh.db.Prepare("DELETE FROM todos WHERE id=?")
	if err != nil {
		panic(err)
	}
	rst, err := stmt.Exec(id)
	if err != nil {
		panic(err)
	}
	cnt, _ := rst.RowsAffected()

	return cnt > 0
}

func (tsh *todoSqlHandler) CompleteTodo(id int, complete bool) bool {
	stmt, err := tsh.db.Prepare("UPDATE todos SET completed=? WHERE id=?")
	if err != nil {
		panic(err)
	}
	rst, err := stmt.Exec(complete, id)
	if err != nil {
		panic(err)
	}
	cnt, _ := rst.RowsAffected()

	return cnt > 0
}

func (tsh *todoSqlHandler) Close() {
	tsh.db.Close()
}

func newTodoSqlHandler(filepath string) TodoDBHandler {
	database, err := sql.Open("sqlite3", filepath)
	if err != nil {
		panic(err)
	}
	//sessionId 탐색을 빠르게 하기 위해 index사용
	statement, _ := database.Prepare(
		`CREATE TABLE IF NOT EXISTS todos (
			id        INTEGER  PRIMARY KEY AUTOINCREMENT,
			sessionId STRING,
			name 			TEXT,
			completed BOOLEAN,
			createdAt DATETIME
		);
		CREATE INDEX IF NOT EXISTS sessionIdIndexOnTodos ON todos (
			sessionId ASC
		);`)
	statement.Exec()
	return &todoSqlHandler{db: database}
}
