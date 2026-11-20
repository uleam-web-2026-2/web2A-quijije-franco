package httpapi

import (
	"github.com/go-chi/chi/v5"
	"github.com/uleam-web-2026-2/servA-quijije-franco/internal/auth"
	"github.com/uleam-web-2026-2/servA-quijije-franco/internal/reservas"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIntegracionReservaProtegida(t *testing.T) {
	a := auth.New("test")
	h := NewReservasHandler(reservas.Servicio{Repo: reservas.NuevaMemoria()})
	r := chi.NewRouter()
	r.Get("/reservas", a.Middleware(http.HandlerFunc(h.Listar)).ServeHTTP)
	q := httptest.NewRequest("GET", "/reservas", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, q)
	if w.Code != 401 {
		t.Fatalf("esperaba 401, obtuvo %d", w.Code)
	}
}
