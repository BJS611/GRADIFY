package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func post(t *testing.T, h http.HandlerFunc, v url.Values) Response {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(v.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	h(w, req)
	var r Response
	json.NewDecoder(w.Result().Body).Decode(&r)
	return r
}

func TestLogin(t *testing.T) {
	if r := post(t, loginHandler, url.Values{"email": {"demo@nusantarakreatif.id"}, "password": {"Demo1234!"}}); !r.Success {
		t.Error("login valid gagal:", r.Message)
	}
	if r := post(t, loginHandler, url.Values{"email": {"demogmail.com"}, "password": {"Demo1234!"}}); r.Success {
		t.Error("email tanpa @ harus ditolak")
	}
	if r := post(t, loginHandler, url.Values{"email": {"demo@nusantarakreatif.id"}, "password": {""}}); r.Success {
		t.Error("password kosong harus ditolak")
	}
}

func regForm(nama, email, telepon, alamat, pass string) url.Values {
	return url.Values{
		"nama": {nama}, "email": {email}, "telepon": {telepon},
		"alamat": {alamat}, "password": {pass}, "konfirmasi": {pass},
	}
}

func TestRegisterStrictAllowlist(t *testing.T) {
	// 1. Success cases
	valid := []url.Values{
		regForm("Budi Santoso", "budi@gmail.com", "081234567890", "Jl. Sudirman No 1", "Pass1234!"),
		regForm("O'Connor", "oconnor@test.co.id", "0811223344", "Gg. Mawar, RT 01/02", "Pass1234!"),
		regForm("Dr. Budi, M.Sc", "dr.budi+123@univ.ac.id", "089912341234", "Kota Bandung-Selatan", "Pass1234!"),
	}

	for i, c := range valid {
		if r := post(t, registerHandler, c); !r.Success {
			t.Errorf("Valid test case %d gagal: %s", i, r.Message)
		}
	}

	// 2. Reject Empty or Invalid Format
	invalid := []url.Values{
		regForm("", "budi@gmail.com", "081234567890", "Jl. Jend. Sudirman", "Pass1234!"),
		regForm("Budi", "budigmail.com", "081234567890", "Jl. Jend. Sudirman", "Pass1234!"),
		regForm("Budi", "budi@gmail.com", "081234567890", "", "Pass1234!"),
		regForm("Budi", "budi@gmail.com", "0812345", "Jl. Jenderal Sudirman", "Pass1234!"),
		regForm("Budi@123", "budi@gmail.com", "081234567890", "Jl. Jend. Sudirman", "Pass1234!"),
		regForm("Budi", "budi@gmail.com", "081234A56789", "Jl. Jenderal Sudirman", "Pass1234!"),
		regForm("Budi", "budi@gmail.com", "081234567890", "Jalan<script>", "Pass1234!"), // Alamat ada karakter < > yang dilarang
	}
	for i, c := range invalid {
		if r := post(t, registerHandler, c); r.Success {
			t.Errorf("Invalid format test case %d lolos: %v", i, c)
		}
	}
}

func TestRegisterBoundaryLimits(t *testing.T) {
	longName := strings.Repeat("A", 101)
	longEmail := strings.Repeat("a", 245) + "@contoh.com" // > 254 chars
	longAlamat := strings.Repeat("B", 101)                // > 100 chars
	longPass := strings.Repeat("P", 73)

	limits := []struct {
		desc string
		form url.Values
	}{
		{"Nama > 100", regForm(longName, "test@test.com", "08123456789", "Jakarta", "Pass1234!")},
		{"Email > 254", regForm("Budi", longEmail, "08123456789", "Jakarta", "Pass1234!")},
		{"Alamat > 100", regForm("Budi", "test@test.com", "08123456789", longAlamat, "Pass1234!")},
		{"Password > 72", regForm("Budi", "test@test.com", "08123456789", "Jakarta", longPass)},
	}

	for _, p := range limits {
		if r := post(t, registerHandler, p.form); r.Success {
			t.Errorf("[%s] Payload melewati boundary limit tapi lolos", p.desc)
		}
	}
}

func TestSecurityPayloads(t *testing.T) {
	payloads := []struct {
		label string
		form  url.Values
	}{
		{"sql injection email", regForm("Budi", "' OR 1=1--@mail.com", "081234567890", "Jalan Sudi", "Pass1234!")},
		{"sql injection nama", regForm("'; DROP TABLE users;--", "budi@gmail.com", "081234567890", "Jalan Sudi", "Pass1234!")},
		{"xss nama", regForm("<script>alert('xss')</script>", "budi@gmail.com", "081234567890", "Jalan Sudi", "Pass1234!")},
		{"xss alamat", regForm("Budi", "budi@gmail.com", "081234567890", "<img onerror=alert(1) src=x>", "Pass1234!")},
		{"sql injection telepon", regForm("Budi", "budi@gmail.com", "1' OR '1'='1", "Jalan Sudi", "Pass1234!")},
	}
	for _, p := range payloads {
		r := post(t, registerHandler, p.form)
		if r.Success {
			t.Errorf("[%s] security payload berbahaya lolos validasi!", p.label)
		}
	}
}
