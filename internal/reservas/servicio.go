package reservas

import (
	"errors"
	"github.com/uleam-web-2026-2/servA-quijije-franco/internal/modelos"
	"time"
)

var (
	ErrNotFound         = errors.New("reserva no encontrada")
	ErrPagoNoVerificado = errors.New("el pago asociado no está verificado")
	ErrForbidden        = errors.New("el usuario no tiene permiso")
)

type Actor struct {
	UsuarioID uint
	Rol       string
}
type Servicio struct {
	Repo  Repositorio
	Pagos PagoRepositorio
}

func (s Servicio) Listar() []modelos.Reserva { return s.Repo.Listar() }
func (s Servicio) Obtener(id uint) (modelos.Reserva, error) {
	r, ok := s.Repo.Obtener(id)
	if !ok {
		return r, ErrNotFound
	}
	return r, nil
}
func (s Servicio) CambiarEstado(id uint, destino string, ahora time.Time, actor Actor) (modelos.Reserva, error) {
	r, e := s.Obtener(id)
	if e != nil {
		return r, e
	}
	if destino == modelos.ReservaConfirmada {
		if actor.Rol != "admin" {
			return r, ErrForbidden
		}
		if s.Pagos == nil {
			return r, ErrPagoNoVerificado
		}
		p, ok := s.Pagos.PagoDeOferta(r.OfertaID)
		if !ok || p.Estado != "verificado" {
			return r, ErrPagoNoVerificado
		}
	}
	if destino == modelos.ReservaCancelada {
		if actor.Rol != "admin" && !s.Repo.EsPropietario(id, actor.UsuarioID) {
			return r, ErrForbidden
		}
	}
	if destino == modelos.ReservaCompletada && actor.Rol != "system" {
		return r, ErrForbidden
	}
	if e = Transicionar(r, destino, ahora); e != nil {
		return r, e
	}
	r.Estado = destino
	return s.Repo.Guardar(r), nil
}

func (s Servicio) CompletarVencidas(ahora time.Time) int {
	n := 0
	for _, r := range s.Repo.Listar() {
		if r.Estado == modelos.ReservaConfirmada && !ahora.Before(r.Fin) {
			if _, e := s.CambiarEstado(r.ID, modelos.ReservaCompletada, ahora, Actor{Rol: "system"}); e == nil {
				n++
			}
		}
	}
	return n
}
