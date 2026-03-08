package main

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"net/mail"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func reply(w http.ResponseWriter, status int, ok bool, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{ok, msg})
}

var (
	validNama   = regexp.MustCompile(`^[a-zA-Z\s.,'-]+$`)
	validAlamat = regexp.MustCompile(`^[a-zA-Z0-9\s.,'-/]+$`)
	validPhone  = regexp.MustCompile(`^08[0-9]{8,11}$`)
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.ParseForm() != nil || r.Method != http.MethodPost {
		reply(w, 400, false, "bad request")
		return
	}
	email, pass := strings.TrimSpace(r.FormValue("email")), r.FormValue("password")
	if len(email) > 254 || len(pass) > 72 {
		reply(w, 400, false, "input terlalu panjang")
		return
	}
	if _, err := mail.ParseAddress(email); err != nil {
		reply(w, 400, false, "email tidak valid")
		return
	}
	if len(pass) < 8 {
		reply(w, 400, false, "password minimal 8 karakter")
		return
	}
	if email == "demo@nusantarakreatif.id" && pass == "Demo1234!" {
		reply(w, 200, true, "login berhasil")
		return
	}
	reply(w, 401, false, "kredensial salah")
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.ParseForm() != nil || r.Method != http.MethodPost {
		reply(w, 400, false, "bad request")
		return
	}

	n := strings.TrimSpace(r.FormValue("nama"))
	e := strings.TrimSpace(r.FormValue("email"))
	t := strings.TrimSpace(r.FormValue("telepon"))
	a := strings.TrimSpace(r.FormValue("alamat"))
	p := r.FormValue("password")
	c := r.FormValue("konfirmasi")

	_, errM := mail.ParseAddress(e)

	switch {
	case len(n) > 100 || len(e) > 254 || len(a) > 100 || len(p) > 72:
		reply(w, 400, false, "input melebihi batas")
	case n == "" || !validNama.MatchString(n):
		reply(w, 400, false, "nama tidak valid")
	case errM != nil:
		reply(w, 400, false, "email tidak valid")
	case a == "" || !validAlamat.MatchString(a):
		reply(w, 400, false, "alamat tidak valid")
	case !validPhone.MatchString(t):
		reply(w, 400, false, "telepon tidak valid")
	case len(p) < 8:
		reply(w, 400, false, "password < 8 karakter")
	case p != c:
		reply(w, 400, false, "pass!=confirm")
	default:
		h, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
		if err != nil {
			reply(w, 500, false, "server error")
			return
		}
		safeNama := html.EscapeString(n)
		safeEmail := html.EscapeString(e)
		safeAlamat := html.EscapeString(a)
		safeTelepon := html.EscapeString(t)

		_ = safeAlamat
		_ = safeTelepon
		fmt.Printf("[REG] %s | %s | %s | %s | %s\n", safeNama, safeEmail, safeTelepon, safeAlamat, h)
		reply(w, 200, true, "pendaftaran berhasil")
	}
}
