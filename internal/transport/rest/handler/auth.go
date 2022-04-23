package handler

import "net/http"

func signIn(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Sign In"))
}

func signUp(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Sign Up"))
}
