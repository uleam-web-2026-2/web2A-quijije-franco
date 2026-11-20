package httpapi

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/uleam-web-2026-2/servA-quijije-franco/internal/auth"
	"github.com/uleam-web-2026-2/servA-quijije-franco/internal/modelos"
	"github.com/uleam-web-2026-2/servA-quijije-franco/internal/reservas"
	"github.com/uleam-web-2026-2/servA-quijije-franco/internal/respuesta"
	"net/http"
	"strconv"
	"time"
)

type ReservasHandler struct{ Service reservas.Servicio }

func NewReservasHandler(s reservas.Servicio) *ReservasHandler { return &ReservasHandler{Service: s} }
func (h *ReservasHandler) Listar(w http.ResponseWriter, r *http.Request) {
	respuesta.Exito(w, 200, h.Service.Listar())
}
func (h *ReservasHandler) Obtener(w http.ResponseWriter, r *http.Request) {
	id, e := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if e != nil || id == 0 {
		respuesta.Error(w, 400, "id_invalido", "el id debe ser positivo")
		return
	}
	x, e := h.Service.Obtener(uint(id))
	if e != nil {
		respuesta.Error(w, 404, "no_encontrado", e.Error())
		return
	}
	respuesta.Exito(w, 200, x)
}
func (h *ReservasHandler) Crear(w http.ResponseWriter, r *http.Request) {
	var in struct {
		modelos.Reserva
		PagoEstado string `json:"pago_estado"`
		PagoURL    string `json:"pago_url"`
	}
	if json.NewDecoder(r.Body).Decode(&in) != nil || in.OfertaID == 0 || in.Inicio.IsZero() || in.Fin.IsZero() || !in.Fin.After(in.Inicio) || in.MontoAcordado <= 0 {
		respuesta.Error(w, 422, "datos_invalidos", "inicio, fin y monto deben ser válidos")
		return
	}
	c := auth.ClaimsFrom(r)
	estadoPago := in.PagoEstado
	if estadoPago == "" {
		estadoPago = "pendiente"
	}
	if estadoPago != "pendiente" && estadoPago != "verificado" && estadoPago != "rechazado" {
		respuesta.Error(w, 422, "datos_invalidos", "pago_estado inválido")
		return
	}
	if estadoPago == "verificado" && c.Rol != modelos.RolAdmin {
		respuesta.Error(w, 403, "rol_no_autorizado", "solo el administrador puede verificar un pago")
		return
	}
	x := in.Reserva
	x.Estado = modelos.ReservaPendiente
	x.EstadoLiquidacion = "pendiente"
	x.ComisionApp = x.MontoAcordado * 0.10
	x.MontoProveedor = x.MontoAcordado - x.ComisionApp
	x = h.Service.Repo.Guardar(x)
	h.Service.Repo.RegistrarPropietario(x.ID, c.UsuarioID)
	if h.Service.Pagos != nil {
		h.Service.Pagos.GuardarPago(modelos.Pago{OfertaID: x.OfertaID, Monto: x.MontoAcordado, URLComprobante: in.PagoURL, Estado: estadoPago, FechaPago: time.Now()})
	}
	respuesta.Exito(w, 201, x)
}
func (h *ReservasHandler) CambiarEstado(w http.ResponseWriter, r *http.Request) {
	id, e := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if e != nil {
		respuesta.Error(w, 400, "id_invalido", "id inválido")
		return
	}
	var in struct {
		Estado string `json:"estado"`
	}
	if json.NewDecoder(r.Body).Decode(&in) != nil || in.Estado == "" {
		respuesta.Error(w, 422, "datos_invalidos", "se requiere estado")
		return
	}
	c := auth.ClaimsFrom(r)
	x, e := h.Service.CambiarEstado(uint(id), in.Estado, time.Now(), reservas.Actor{UsuarioID: c.UsuarioID, Rol: c.Rol})
	if e == reservas.ErrNotFound {
		respuesta.Error(w, 404, "no_encontrado", e.Error())
		return
	}
	if e == reservas.ErrTransicion {
		respuesta.Error(w, 422, "transicion_prohibida", e.Error())
		return
	}
	if e == reservas.ErrPagoNoVerificado {
		respuesta.Error(w, 422, "pago_no_verificado", e.Error())
		return
	}
	if e == reservas.ErrForbidden {
		respuesta.Error(w, 403, "rol_no_autorizado", e.Error())
		return
	}
	if e != nil {
		respuesta.Error(w, 500, "error_interno", "no se pudo cambiar el estado")
		return
	}
	respuesta.Exito(w, 200, x)
}
