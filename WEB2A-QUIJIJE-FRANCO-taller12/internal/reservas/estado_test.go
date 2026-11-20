package reservas

import (
	"github.com/uleam-web-2026-2/servA-quijije-franco/internal/modelos"
	"testing"
	"time"
)

func TestTransicionProhibidaCompletadaACancelada(t *testing.T) {
	r := modelos.Reserva{Estado: modelos.ReservaCompletada, Inicio: time.Now().Add(-time.Hour)}
	if Transicionar(r, modelos.ReservaCancelada, time.Now()) != ErrTransicion {
		t.Fatal("una reserva completada no puede cancelarse")
	}
}
func TestTransicionConfirmadaACompletada(t *testing.T) {
	r := modelos.Reserva{Estado: modelos.ReservaConfirmada, Fin: time.Now().Add(-time.Minute)}
	if Transicionar(r, modelos.ReservaCompletada, time.Now()) != nil {
		t.Fatal("el evento terminado debe completarse")
	}
}

func TestConfirmacionRequierePagoVerificado(t *testing.T) {
	casos := []struct {
		nombre, estado string
		existe         bool
		debeFallar     bool
	}{
		{"inexistente", "", false, true}, {"pendiente", "pendiente", true, true}, {"rechazado", "rechazado", true, true}, {"verificado", "verificado", true, false},
	}
	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			m := NuevaMemoria()
			s := Servicio{Repo: m, Pagos: m}
			r := m.Guardar(modelos.Reserva{OfertaID: 10, Estado: modelos.ReservaPendiente, Inicio: time.Now().Add(time.Hour), Fin: time.Now().Add(2 * time.Hour)})
			if tc.existe {
				m.GuardarPago(modelos.Pago{OfertaID: r.OfertaID, Estado: tc.estado})
			}
			_, err := s.CambiarEstado(r.ID, modelos.ReservaConfirmada, time.Now(), Actor{UsuarioID: 1, Rol: modelos.RolAdmin})
			if tc.debeFallar && err != ErrPagoNoVerificado {
				t.Fatalf("esperaba pago no verificado, obtuvo %v", err)
			}
			if !tc.debeFallar && err != nil {
				t.Fatalf("el pago verificado debe permitir confirmar: %v", err)
			}
		})
	}
}

func TestCancelacionClientePropietarioYRestriccionTemporal(t *testing.T) {
	m := NuevaMemoria()
	s := Servicio{Repo: m, Pagos: m}
	r := m.Guardar(modelos.Reserva{OfertaID: 20, Estado: modelos.ReservaPendiente, Inicio: time.Now().Add(time.Hour), Fin: time.Now().Add(2 * time.Hour)})
	m.RegistrarPropietario(r.ID, 7)
	if _, err := s.CambiarEstado(r.ID, modelos.ReservaCancelada, time.Now(), Actor{UsuarioID: 7, Rol: modelos.RolCliente}); err != nil {
		t.Fatalf("el cliente propietario debe cancelar: %v", err)
	}
	r = m.Guardar(modelos.Reserva{OfertaID: 21, Estado: modelos.ReservaConfirmada, Inicio: time.Now().Add(-time.Minute), Fin: time.Now().Add(time.Hour)})
	m.RegistrarPropietario(r.ID, 7)
	if _, err := s.CambiarEstado(r.ID, modelos.ReservaCancelada, time.Now(), Actor{UsuarioID: 7, Rol: modelos.RolCliente}); err != ErrTransicion {
		t.Fatalf("no debe cancelar después del inicio: %v", err)
	}
}
