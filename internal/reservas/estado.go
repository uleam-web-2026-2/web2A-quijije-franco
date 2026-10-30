package reservas

import (
	"errors"
	"github.com/uleam-web-2026-2/servA-quijije-franco/internal/modelos"
	"time"
)

var ErrTransicion = errors.New("transicion de estado no permitida")

func Transicionar(r modelos.Reserva, destino string, ahora time.Time) error {
	switch r.Estado {
	case modelos.ReservaPendiente:
		if destino == modelos.ReservaConfirmada || destino == modelos.ReservaCancelada {
			r.Estado = destino
			return nil
		}
	case modelos.ReservaConfirmada:
		if destino == modelos.ReservaCancelada && ahora.Before(r.Inicio) || destino == modelos.ReservaCompletada && (!ahora.Before(r.Fin)) {
			return nil
		}
	}
	return ErrTransicion
}
