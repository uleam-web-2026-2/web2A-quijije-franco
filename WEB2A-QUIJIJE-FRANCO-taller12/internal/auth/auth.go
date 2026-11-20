package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/uleam-web-2026-2/servA-quijije-franco/internal/respuesta"
	"net/http"
	"strings"
	"time"
)

type Claims struct {
	UsuarioID uint   `json:"usuario_id"`
	Rol       string `json:"rol"`
	jwt.RegisteredClaims
}
type Service struct{ Secret []byte }

func New(secret string) *Service {
	if secret == "" {
		secret = "hito2-secret-local"
	}
	return &Service{[]byte(secret)}
}
func (s *Service) Token(id uint, rol string) (string, error) {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	c := Claims{id, rol, jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)), IssuedAt: jwt.NewNumericDate(time.Now()), ID: fmt.Sprintf("%x", b)}}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(s.Secret)
}
func (s *Service) Parse(raw string) (*Claims, error) {
	t, e := jwt.ParseWithClaims(raw, &Claims{}, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("algoritmo no permitido")
		}
		return s.Secret, nil
	})
	if e != nil || !t.Valid {
		return nil, fmt.Errorf("token invalido")
	}
	c, ok := t.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("claims invalidos")
	}
	return c, nil
}

type key int

const claimsKey key = 1

func (s *Service) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			respuesta.Error(w, http.StatusUnauthorized, "no_autenticado", "se requiere un token Bearer")
			return
		}
		c, e := s.Parse(strings.TrimPrefix(h, "Bearer "))
		if e != nil {
			respuesta.Error(w, http.StatusUnauthorized, "token_invalido", "el token no es válido")
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), claimsKey, c)))
	})
}
func ClaimsFrom(r *http.Request) *Claims { c, _ := r.Context().Value(claimsKey).(*Claims); return c }
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := ClaimsFrom(r)
			for _, role := range roles {
				if c != nil && c.Rol == role {
					next.ServeHTTP(w, r)
					return
				}
			}
			respuesta.Error(w, http.StatusForbidden, "rol_no_autorizado", "el rol no puede ejecutar esta acción")
		})
	}
}
