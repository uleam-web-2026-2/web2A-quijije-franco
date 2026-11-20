package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
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
		databaseConn, err := database.Open(ctx, url)
		if err != nil {
			cancel()
			log.Fatalf("no se pudo preparar PostgreSQL: %v", err)
		}
		defer databaseConn.Close()
		if os.Getenv("RUN_MIGRATIONS") == "true" {
			if err := database.Migrate(ctx, databaseConn); err != nil {
				log.Fatalf("migración administrativa fallida: %v", err)
			}
			log.Println("migración administrativa completada; no se inicia la API")
			return
		}
		repo = reservas.SQLRepositorio{DB: databaseConn}
		cancel()
		log.Println("PostgreSQL conectado sin migración automática (cuenta de API)")
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
	// La petición previa OPTIONS no trae JWT; CORS debe instalarse antes de las rutas.
	// Solo el front local previsto puede leer las respuestas desde un navegador.
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
		MaxAge:         300,
	}))
	// Deriva la IP de RemoteAddr, no de una cabecera arbitraria del cliente.
	r.Use(chimw.ClientIPFromRemoteAddr)
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
	// Cinco intentos por minuto y por IP, solo para el inicio de sesión.
	r.With(httprate.LimitBy(5, time.Minute, porIP)).Post("/login", login.Login)
	r.Group(func(p chi.Router) {
		p.Use(a.Middleware)
		p.Get("/reservas", h.Listar)
		p.Get("/reservas/{id}", h.Obtener)
		p.Post("/reservas", h.Crear)
		p.With(auth.RequireRole("admin", "cliente")).Patch("/reservas/{id}/estado", h.CambiarEstado)
	})
	return r
}

// porIP agrupa los intentos por IP de la conexión de red.
// No confiar en X-Forwarded-For hasta configurar un proxy conocido.
func porIP(r *http.Request) (string, error) {
	ip := chimw.GetClientIP(r.Context())
	return httprate.CanonicalizeIP(ip), nil
}
