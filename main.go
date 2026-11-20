package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/uleam-web-2026-2/servA-quijije-franco/internal/auth"
	database "github.com/uleam-web-2026-2/servA-quijije-franco/internal/db"
	"github.com/uleam-web-2026-2/servA-quijije-franco/internal/httpapi"
	"github.com/uleam-web-2026-2/servA-quijije-franco/internal/middleware"
	"github.com/uleam-web-2026-2/servA-quijije-franco/internal/reservas"
	"github.com/uleam-web-2026-2/servA-quijije-franco/internal/respuesta"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	var repo reservas.Repositorio = reservas.NuevaMemoria()
	if url := os.Getenv("DATABASE_URL"); url != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		databaseConn, err := database.OpenAndMigrate(ctx, url)
		cancel()
		if err != nil {
			log.Fatalf("no se pudo preparar PostgreSQL: %v", err)
		}
		defer databaseConn.Close()
		repo = reservas.SQLRepositorio{DB: databaseConn}
		log.Println("PostgreSQL conectado; migración idempotente aplicada")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("API de reservas deportivas en :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, NewRouterWithRepo(repo)))
}
func NewRouter() *chi.Mux {
	return NewRouterWithRepo(reservas.NuevaMemoria())
}
func NewRouterWithRepo(repo reservas.Repositorio) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Registro, middleware.Recuperacion)
	r.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		respuesta.Error(w, 404, "ruta_inexistente", "la ruta no existe")
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, _ *http.Request) {
		respuesta.Error(w, 405, "metodo_no_permitido", "el método no está permitido")
	})
	a := auth.New(os.Getenv("JWT_SECRET"))
	var pagos reservas.PagoRepositorio
	if p, ok := repo.(reservas.PagoRepositorio); ok {
		pagos = p
	}
	h := httpapi.NewReservasHandler(reservas.Servicio{Repo: repo, Pagos: pagos})
	ad, _ := bcrypt.GenerateFromPassword([]byte("admin-2026"), bcrypt.MinCost)
	cl, _ := bcrypt.GenerateFromPassword([]byte("cliente-2026"), bcrypt.MinCost)
	login := httpapi.LoginHandler{Auth: a, Accounts: map[string]httpapi.Account{"admin@reservas.local": {ID: 1, Password: string(ad), Rol: "admin"}, "cliente@reservas.local": {ID: 2, Password: string(cl), Rol: "cliente"}}}
	r.Post("/login", login.Login)
	r.Group(func(p chi.Router) {
		p.Use(a.Middleware)
		p.Get("/reservas", h.Listar)
		p.Get("/reservas/{id}", h.Obtener)
		p.Post("/reservas", h.Crear)
		p.With(auth.RequireRole("admin", "cliente")).Patch("/reservas/{id}/estado", h.CambiarEstado)
	})
	return r
}
