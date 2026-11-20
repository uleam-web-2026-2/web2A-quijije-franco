package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFlujoDemostracionReservas(t *testing.T) {
	r := NewRouter()
	req := httptest.NewRequest(http.MethodGet, "/reservas", nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)
	if res.Code != 401 {
		t.Fatalf("sin token: %d", res.Code)
	}
	login := func(correo, password string) string {
		b, _ := json.Marshal(map[string]string{"correo": correo, "password": password})
		q := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(b))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, q)
		if w.Code != 200 {
			t.Fatalf("login %s: %d", correo, w.Code)
		}
		var out struct {
			Datos struct {
				Token string `json:"token"`
			} `json:"datos"`
		}
		_ = json.NewDecoder(w.Body).Decode(&out)
		return out.Datos.Token
	}
	admin := login("admin@reservas.local", "admin-2026")
	cliente := login("cliente@reservas.local", "cliente-2026")
	body := bytes.NewBufferString(`{"oferta_id":1,"inicio":"2030-10-15T19:00:00Z","fin":"2030-10-15T20:30:00Z","monto_acordado":35,"pago_estado":"verificado","pago_url":"demo://comprobante"}`)
	q := httptest.NewRequest(http.MethodPost, "/reservas", body)
	q.Header.Set("Authorization", "Bearer "+admin)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, q)
	if w.Code != 201 {
		t.Fatalf("crear: %d", w.Code)
	}
	var created struct {
		Datos struct {
			ID uint `json:"reserva_id"`
		} `json:"datos"`
	}
	_ = json.NewDecoder(w.Body).Decode(&created)
	change := func(token, state string) int {
		b := bytes.NewBufferString(`{"estado":"` + state + `"}`)
		q := httptest.NewRequest(http.MethodPatch, "/reservas/"+itoa(created.Datos.ID)+"/estado", b)
		q.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, q)
		return w.Code
	}
	if got := change(cliente, "confirmada"); got != 403 {
		t.Fatalf("rol: %d", got)
	}
	if got := change(admin, "confirmada"); got != 200 {
		t.Fatalf("admin: %d", got)
	}
	if got := change(admin, "confirmada"); got != 422 {
		t.Fatalf("prohibida: %d", got)
	}
}
func itoa(v uint) string {
	if v == 0 {
		return "0"
	}
	b := make([]byte, 0, 10)
	for v > 0 {
		b = append([]byte{byte('0' + v%10)}, b...)
		v /= 10
	}
	return string(b)
}
