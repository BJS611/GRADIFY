package main

import "strings"

func validEmail(s string) bool {
	s = strings.TrimSpace(s)
	at := strings.Index(s, "@")
	if at <= 0 {
		return false
	}
	d := s[at+1:]
	return len(d) > 0 && strings.Index(d, ".") > 0
}

func validPhone(s string) bool {
	s = strings.TrimSpace(s)
	return strings.HasPrefix(s, "08") && len(s) >= 10 && len(s) <= 14
}

func validateLogin(email, pass string) (bool, string) {
	if !validEmail(email) {
		return false, "email tidak valid"
	}
	if len(pass) < 8 {
		return false, "password minimal 8 karakter"
	}
	return true, "ok"
}

func validateRegister(nama, email, phone, pass, confirm string) (bool, string) {
	if strings.TrimSpace(nama) == "" {
		return false, "nama tidak boleh kosong"
	}
	if !validEmail(email) {
		return false, "email tidak valid"
	}
	if !validPhone(phone) {
		return false, "nomor telepon tidak valid"
	}
	if len(pass) < 8 {
		return false, "password minimal 8 karakter"
	}
	if pass != confirm {
		return false, "konfirmasi password tidak cocok"
	}
	return true, "ok"
}
