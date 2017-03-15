package main

import (
	"net/http"
	"strings"
	"fmt"
	"log"
	"crypto/md5"
	"io"
)

type ChatUser interface {
	UniqueID() string
	AvatarURL() string
}

type chatUser struct {
	uniqueID string
}

func (u chatUser) AvatarURL() string {
	return ""
}
func (u chatUser) UniqueID() string {
	return u.uniqueID
}

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth")
	if  err == http.ErrNoCookie || cookie.Value == "" {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.next.ServeHTTP(w, r)
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next:handler}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	provider := segs[3]
	switch action {
	case "login":
		log.Print("TODO for action login with Provider", provider)
	case "callback":
		authCookieValue := r.URL.Query().Get("cookie")
		m := md5.New()
		io.WriteString(m, authCookieValue)
		userID := fmt.Sprintf("%x", m.Sum(nil))


		http.SetCookie(w, &http.Cookie{
			Name:"auth",
			Value:userID,
			Path: "/",
		})

		chatUser := &chatUser{uniqueID:userID}
		avatarURL, err := avatars.GetAvatarURL(chatUser)
		if err != nil {
			log.Println("Error when trying to GetAvatarURL", "-", err)
		}
		http.SetCookie(w, &http.Cookie{
			Name:"avatar_url",
			Value:avatarURL,
			Path:"/",
		})
		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action %s not supported", action)
	}
}

