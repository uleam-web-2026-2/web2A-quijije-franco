package httpapi

import (
	"encoding/json"
	"github.com/uleam-web-2026-2/servA-quijije-franco/internal/auth"
	"github.com/uleam-web-2026-2/servA-quijije-franco/internal/respuesta"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

type Account struct {
	ID                    uint
	Correo, Password, Rol string
}
type LoginHandler struct {
	Auth     *auth.Service
	Accounts map[string]Account
}

func (h LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Correo   string `json:"correo"`
		Password string `json:"password"`
	}
	if json.NewDecoder(r.Body).Decode(&in) != nil {
		respuesta.Error(w, 400, "cuerpo_malformado", "JSON inválido")
		return
	}
	a, ok := h.Accounts[strings.ToLower(strings.TrimSpace(in.Correo))]
	if !ok || bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(in.Password)) != nil {
		respuesta.Error(w, 401, "credenciales_invalidas", "correo o contraseña incorrectos")
		return
	}
	t, e := h.Auth.Token(a.ID, a.Rol)
	if e != nil {
		respuesta.Error(w, 500, "error_interno", "no se pudo crear token")
		return
	}
	respuesta.Exito(w, 200, map[string]string{"token": t, "rol": a.Rol})
}
