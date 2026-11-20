package reservas

import (
	"github.com/uleam-web-2026-2/servA-quijije-franco/internal/modelos"
	"sync"
)

type Repositorio interface {
	Listar() []modelos.Reserva
	Obtener(uint) (modelos.Reserva, bool)
	Guardar(modelos.Reserva) modelos.Reserva
	RegistrarPropietario(reservaID, usuarioID uint)
	EsPropietario(reservaID, usuarioID uint) bool
}
type PagoRepositorio interface {
	PagoDeOferta(ofertaID uint) (modelos.Pago, bool)
	GuardarPago(modelos.Pago) modelos.Pago
}
type Memoria struct {
	mu             sync.RWMutex
	next, nextPago uint
	items          map[uint]modelos.Reserva
	owners         map[uint]uint
	payments       map[uint]modelos.Pago
}

func NuevaMemoria() *Memoria {
	return &Memoria{items: map[uint]modelos.Reserva{}, owners: map[uint]uint{}, payments: map[uint]modelos.Pago{}}
}
func (m *Memoria) Listar() []modelos.Reserva {
	m.mu.RLock()
	defer m.mu.RUnlock()
	x := make([]modelos.Reserva, 0, len(m.items))
	for _, r := range m.items {
		x = append(x, r)
	}
	return x
}
func (m *Memoria) Obtener(id uint) (modelos.Reserva, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	r, ok := m.items[id]
	return r, ok
}
func (m *Memoria) Guardar(r modelos.Reserva) modelos.Reserva {
	m.mu.Lock()
	defer m.mu.Unlock()
	if r.ID == 0 {
		m.next++
		r.ID = m.next
	}
	m.items[r.ID] = r
	return r
}
func (m *Memoria) RegistrarPropietario(reservaID, usuarioID uint) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.owners[reservaID] = usuarioID
}
func (m *Memoria) EsPropietario(reservaID, usuarioID uint) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.owners[reservaID] == usuarioID
}
func (m *Memoria) PagoDeOferta(ofertaID uint) (modelos.Pago, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.payments[ofertaID]
	return p, ok
}
func (m *Memoria) GuardarPago(p modelos.Pago) modelos.Pago {
	m.mu.Lock()
	defer m.mu.Unlock()
	if p.ID == 0 {
		m.nextPago++
		p.ID = m.nextPago
	}
	m.payments[p.OfertaID] = p
	return p
}
