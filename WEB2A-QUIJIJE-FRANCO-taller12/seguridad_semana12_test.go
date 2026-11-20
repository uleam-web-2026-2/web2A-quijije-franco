package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Ejercicios del Taller 12: preflight, IP remota y rate limit en /login.
func TestSemana12CORSOrigenPermitidoYOtro(t *testing.T) {
	r := NewRouter()
	for _, tc := range []struct {
		origen    string
		permitido bool
	}{
		{"http://localhost:5173", true},
		{"http://malo.example", false},
	} {
		t.Run(tc.origen, func(t *testing.T) {
			q := httptest.NewRequest(http.MethodOptions, "/reservas/1/estado", nil)
			q.Header.Set("Origin", tc.origen)
			q.Header.Set("Access-Control-Request-Method", http.MethodPatch)
			q.Header.Set("Access-Control-Request-Headers", "authorization,content-type")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, q)
			if tiene := w.Header().Get("Access-Control-Allow-Origin") == tc.origen; tiene != tc.permitido {
				t.Fatalf("preflight de %s: permiso=%v, cabeceras=%v", tc.origen, tiene, w.Header())
			}
		})
	}
}

func TestSemana12RafagaSoloLogin(t *testing.T) {
	r := NewRouter()
	for intento := 1; intento <= 7; intento++ {
		q := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"correo":"admin@reservas.local","password":"incorrecta"}`))
		q.Header.Set("Content-Type", "application/json")
		q.RemoteAddr = "192.0.2.123:5000"
		w := httptest.NewRecorder()
		r.ServeHTTP(w, q)
		want := http.StatusUnauthorized
		if intento >= 6 {
			want = http.StatusTooManyRequests
		}
		if w.Code != want {
			t.Fatalf("intento %d: HTTP %d; esperaba %d. Cuerpo: %s", intento, w.Code, want, w.Body.String())
		}
		if intento >= 6 && w.Header().Get("Retry-After") == "" {
			t.Fatal("el 429 debe incluir Retry-After")
		}
	}
	// El limite no afecta las rutas ajenas al POST /login.
	for n := 0; n < 7; n++ {
		q := httptest.NewRequest(http.MethodGet, "/reservas", nil)
		q.RemoteAddr = fmt.Sprintf("192.0.2.123:%d", 6000+n)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, q)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("GET sin token tras rafaga: HTTP %d", w.Code)
		}
	}
}

func TestSemana12IDConTextoMalicioso(t *testing.T) {
	r := NewRouter()
	login := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"correo":"admin@reservas.local","password":"admin-2026"}`))
	login.Header.Set("Content-Type", "application/json")
	wLogin := httptest.NewRecorder()
	r.ServeHTTP(wLogin, login)
	if wLogin.Code != http.StatusOK {
		t.Fatalf("inicio de sesión: %d", wLogin.Code)
	}
	var out struct {
		Datos struct {
			Token string `json:"token"`
		} `json:"datos"`
	}
	if err := json.NewDecoder(wLogin.Body).Decode(&out); err != nil || out.Datos.Token == "" {
		t.Fatalf("no se pudo extraer token: %v", err)
	}
	q := httptest.NewRequest(http.MethodGet, "/reservas/1%20OR%201=1", nil)
	q.Header.Set("Authorization", "Bearer "+out.Datos.Token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, q)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("id con texto malicioso: HTTP %d, esperaba 400", w.Code)
	}
}
