package main

import (
	"net/http"
	"log"
	"sync"
	"text/template"
	"path/filepath"
	"flag"
	"github.com/ruandao/goPractice/trace"
	"os"
	"fmt"
)

var avatars Avatar = TryAvatars{
	UseFileSystemAvatar,
	UseAuthAvatar,
	UseGravatarAvatar,
}

type templateHandler struct {
	once 	sync.Once
	filename	string
	templ	*template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	userId, err := r.Cookie("auth")
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err.Error()), http.StatusInternalServerError)
		return
	}
	data := map[string]interface{} {
		"Host": r.Host,
		"name": userId.Value,
		"userid": userId.Value,
	}

	t.templ.Execute(w, data)
}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse()

	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets/"))))
	http.Handle("/chat", MustAuth(&templateHandler{filename:"chat.html"}))
	http.Handle("/login", &templateHandler{filename:"login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name: "auth",
			Value: "",
			Path: "/",
			MaxAge: -1,
		})
		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	http.Handle("/avatars/",
		http.StripPrefix("/avatars/",
			http.FileServer(http.Dir("./avatars"))))
	http.Handle("/upload", &templateHandler{filename:"upload.html"})
	http.HandleFunc("/uploader", uploaderHandler)
	http.Handle("/room", r)
	go r.run()
	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}