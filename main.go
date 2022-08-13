package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	// "strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	ID       uint `gorm:"primaryKey;autoIncrement:true"`
	Title    string
	Assignee string
	Deadline string
	Done     bool
}

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Todo{})

	tmpl := template.Must(template.ParseFiles("template/index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}
			title := r.FormValue("title")
			assignee := r.FormValue("assignee")
			deadline := r.FormValue("deadline")
			db.Create(&Todo{Title: title, Assignee: assignee, Deadline: deadline, Done: false})
		}
		//Request not POST
		var todos []Todo
		db.Find(&todos)
		data := TodoPageData{
			PageTitle: "My TODO list",
			Todos:     todos,
		}
		tmpl.Execute(w, data)
	})

	http.HandleFunc("/update/", func(w http.ResponseWriter, r *http.Request) {
		// id := strings.TrimPrefix(r.URL.Path, "/update/")
		id := r.FormValue("task_id")
		var todo Todo
		db.First(&todo, id)
		todo.Title = r.FormValue("title")
		todo.Assignee = r.FormValue("assignee")
		todo.Deadline = r.FormValue("deadline")
		db.Save(&todo)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/done/", func(w http.ResponseWriter, r *http.Request) {
		// id := strings.TrimPrefix(r.URL.Path, "/done/")
		id := r.FormValue("task_id")
		var todo Todo
		db.First(&todo, id)
		todo.Done = true
		db.Save(&todo)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/delete/", func(w http.ResponseWriter, r *http.Request) {
		// id := strings.TrimPrefix(r.URL.Path, "/delete/")
		id := r.FormValue("task_id")
		db.Delete(&Todo{}, id)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.ListenAndServe(":"+port, nil)
}
