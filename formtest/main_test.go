package main

import "testing"

func TestEmail(t *testing.T) {
	valid := []string{"budi@gmail.com", "a@b.co.id"}
	for _, e := range valid {
		if !validEmail(e) {
			t.Errorf("%q harus valid", e)
		}
	}
	invalid := []string{"", "budigmail.com", "budi@", "@gmail.com", "budi@gmailcom"}
	for _, e := range invalid {
		if validEmail(e) {
			t.Errorf("%q harus ditolak", e)
		}
	}
}

func TestPhone(t *testing.T) {
	valid := []string{"081234567890", "0812345678", "08123456789012"}
	for _, p := range valid {
		if !validPhone(p) {
			t.Errorf("%q harus valid", p)
		}
	}
	invalid := []string{"", "0812345", "021234567", "081234567890123456"}
	for _, p := range invalid {
		if validPhone(p) {
			t.Errorf("%q harus ditolak", p)
		}
	}
}

func TestLogin(t *testing.T) {
	if ok, _ := validateLogin("budi@gmail.com", "Pass1234"); !ok {
		t.Error("login valid gagal")
	}
	if ok, _ := validateLogin("budigmail.com", "Pass1234"); ok {
		t.Error("email salah harus ditolak")
	}
	if ok, _ := validateLogin("budi@gmail.com", "abc"); ok {
		t.Error("password pendek harus ditolak")
	}
}

func TestRegister(t *testing.T) {
	ok, _ := validateRegister("Budi", "budi@gmail.com", "081234567890", "Pass1234", "Pass1234")
	if !ok {
		t.Error("register valid gagal")
	}

	cases := []struct{ nama, email, phone, pass, confirm string }{
		{"", "budi@gmail.com", "081234567890", "Pass1234", "Pass1234"},
		{"Budi", "budigmail.com", "081234567890", "Pass1234", "Pass1234"},
		{"Budi", "budi@gmail.com", "0812345", "Pass1234", "Pass1234"},
		{"Budi", "budi@gmail.com", "081234567890", "short", "short"},
		{"Budi", "budi@gmail.com", "081234567890", "Pass1234", "Pass5678"},
	}
	for _, c := range cases {
		ok, _ := validateRegister(c.nama, c.email, c.phone, c.pass, c.confirm)
		if ok {
			t.Errorf("input tidak valid lolos: %+v", c)
		}
	}
}
