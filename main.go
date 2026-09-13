package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/uleam-web-2026-2/servA-quijije-franco/internal/middleware"
	"github.com/uleam-web-2026-2/servA-quijije-franco/internal/respuesta"
	"github.com/uleam-web-2026-2/servA-quijije-franco/internal/tickets"
)

func main() {
	r := chi.NewRouter()

	// Middlewares globales (Fase 2)
	r.Use(middleware.Registro)
	r.Use(middleware.Recuperacion)

	// Sobrescribir respuestas por defecto de Chi para mantener formato JSON
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		respuesta.Error(w, http.StatusNotFound, "ruta_inexistente", "la ruta no existe")
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		respuesta.Error(w, http.StatusMethodNotAllowed, "metodo_no_permitido", "el método no está permitido en esta ruta")
	})

	// Inicialización del almacén de datos
	almacen := tickets.NuevoAlmacen()

	// Rutas de la API
	r.Get("/tickets", almacen.Listar)
	r.Post("/tickets", almacen.Crear)
	r.Get("/tickets/{id}", almacen.Obtener)
	r.Get("/explotar", tickets.Explotar)

	log.Println("mesa de ayuda en :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
