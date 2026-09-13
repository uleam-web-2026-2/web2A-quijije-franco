package middleware

import (
	"log"
	"net/http"
	"time"

	// REMPLAZA con la ruta de tu módulo
	"github.com/uleam-web-2026-2/servA-quijije-franco/internal/respuesta" 
)

type grabador struct {
	http.ResponseWriter
	estado int
}

func (g *grabador) WriteHeader(codigo int) {
	g.estado = codigo
	g.ResponseWriter.WriteHeader(codigo)
}

// Middleware de Registro (Logging)
func Registro(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inicio := time.Now()
		g := &grabador{ResponseWriter: w, estado: http.StatusOK}
		
		next.ServeHTTP(g, r) // Procesa la petición

		// Se coloca DESPUÉS de next.ServeHTTP porque el estadoHTTP y el tiempo transcurrido
		// solo se conocen una vez que el handler ha respondido.
		log.Printf("%s %s -> %d (%s)", r.Method, r.URL.Path, g.estado, time.Since(inicio))
	})
}

// Middleware de Recuperación (Recovery de Pánicos)
func Recuperacion(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if p := recover(); p != nil {
				log.Printf("PÁNICO en %s %s: %v", r.Method, r.URL.Path, p)
				respuesta.Error(w, http.StatusInternalServerError, "error_interno", "ocurrió un error inesperado")
			}
		}()
		next.ServeHTTP(w, r)
	})
}