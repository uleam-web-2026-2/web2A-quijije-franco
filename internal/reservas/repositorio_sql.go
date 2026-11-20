package reservas

import (
	"context"
	"database/sql"
	"github.com/uleam-web-2026-2/servA-quijije-franco/internal/modelos"
)

// SQLRepositorio persiste Reserva en PostgreSQL. La máquina de estados no
// depende de esta capa: permanece en Servicio y estado.go.
type SQLRepositorio struct{ DB *sql.DB }

func (r SQLRepositorio) Listar() []modelos.Reserva {
	rows, err := r.DB.QueryContext(context.Background(), `SELECT reserva_id,oferta_id,inicio,fin,monto_acordado,comision_app,monto_proveedor,estado,estado_liquidacion,COALESCE(comprobante_liquidacion,'') FROM reservas ORDER BY reserva_id`)
	if err != nil {
		return []modelos.Reserva{}
	}
	defer rows.Close()
	var out []modelos.Reserva
	for rows.Next() {
		var x modelos.Reserva
		if rows.Scan(&x.ID, &x.OfertaID, &x.Inicio, &x.Fin, &x.MontoAcordado, &x.ComisionApp, &x.MontoProveedor, &x.Estado, &x.EstadoLiquidacion, &x.ComprobanteLiquidacion) == nil {
			out = append(out, x)
		}
	}
	return out
}

func (r SQLRepositorio) Obtener(id uint) (modelos.Reserva, bool) {
	var x modelos.Reserva
	err := r.DB.QueryRowContext(context.Background(), `SELECT reserva_id,oferta_id,inicio,fin,monto_acordado,comision_app,monto_proveedor,estado,estado_liquidacion,COALESCE(comprobante_liquidacion,'') FROM reservas WHERE reserva_id=$1`, id).Scan(&x.ID, &x.OfertaID, &x.Inicio, &x.Fin, &x.MontoAcordado, &x.ComisionApp, &x.MontoProveedor, &x.Estado, &x.EstadoLiquidacion, &x.ComprobanteLiquidacion)
	return x, err == nil
}

func (r SQLRepositorio) Guardar(x modelos.Reserva) modelos.Reserva {
	if x.ID == 0 {
		_ = r.DB.QueryRowContext(context.Background(), `INSERT INTO reservas(oferta_id,inicio,fin,monto_acordado,comision_app,monto_proveedor,estado,estado_liquidacion,comprobante_liquidacion) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING reserva_id`, x.OfertaID, x.Inicio, x.Fin, x.MontoAcordado, x.ComisionApp, x.MontoProveedor, x.Estado, x.EstadoLiquidacion, x.ComprobanteLiquidacion).Scan(&x.ID)
		return x
	}
	_, _ = r.DB.ExecContext(context.Background(), `UPDATE reservas SET estado=$2,comprobante_liquidacion=$3 WHERE reserva_id=$1`, x.ID, x.Estado, x.ComprobanteLiquidacion)
	return x
}
func (r SQLRepositorio) RegistrarPropietario(_, _ uint) {}
func (r SQLRepositorio) EsPropietario(reservaID, usuarioID uint) bool {
	var ok bool
	_ = r.DB.QueryRowContext(context.Background(), `SELECT EXISTS(SELECT 1 FROM reservas r JOIN ofertas o ON o.oferta_id=r.oferta_id JOIN solicitudes s ON s.solicitud_id=o.solicitud_id JOIN clientes c ON c.cliente_id=s.cliente_id WHERE r.reserva_id=$1 AND c.usuario_id=$2)`, reservaID, usuarioID).Scan(&ok)
	return ok
}
func (r SQLRepositorio) PagoDeOferta(ofertaID uint) (modelos.Pago, bool) {
	var p modelos.Pago
	e := r.DB.QueryRowContext(context.Background(), `SELECT pago_id,oferta_id,monto,url_comprobante,estado,fecha_pago FROM pagos WHERE oferta_id=$1`, ofertaID).Scan(&p.ID, &p.OfertaID, &p.Monto, &p.URLComprobante, &p.Estado, &p.FechaPago)
	return p, e == nil
}
func (r SQLRepositorio) GuardarPago(p modelos.Pago) modelos.Pago {
	if p.ID == 0 {
		_ = r.DB.QueryRowContext(context.Background(), `INSERT INTO pagos(oferta_id,monto,url_comprobante,estado,fecha_pago) VALUES($1,$2,$3,$4,$5) RETURNING pago_id`, p.OfertaID, p.Monto, p.URLComprobante, p.Estado, p.FechaPago).Scan(&p.ID)
	}
	return p
}
