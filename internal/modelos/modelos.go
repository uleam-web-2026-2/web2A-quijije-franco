package modelos

import "time"

const (
	RolCliente          = "cliente"
	RolProveedor        = "proveedor"
	RolAdmin            = "admin"
	ReservaPendiente    = "pendiente"
	ReservaConfirmada   = "confirmada"
	ReservaCompletada   = "completada"
	ReservaCancelada    = "cancelada"
	SolicitudAbierta    = "abierta"
	SolicitudAdjudicada = "adjudicada"
	SolicitudCancelada  = "cancelada"
)

type Usuario struct {
	ID       uint   `json:"usuario_id"`
	Nombre   string `json:"nombre"`
	Correo   string `json:"correo"`
	Password string `json:"-"`
	Rol      string `json:"rol"`
}
type Cliente struct {
	ID        uint   `json:"cliente_id"`
	UsuarioID uint   `json:"usuario_id"`
	Cedula    string `json:"cedula"`
	Telefono  string `json:"telefono"`
}
type Proveedor struct {
	ID            uint   `json:"proveedor_id"`
	UsuarioID     uint   `json:"usuario_id"`
	RUC           string `json:"ruc"`
	Telefono      string `json:"telefono"`
	Banco         string `json:"banco"`
	NumeroCuenta  string `json:"numero_cuenta"`
	TitularCuenta string `json:"titular_cuenta"`
}
type Complejo struct {
	ID          uint   `json:"complejo_id"`
	ProveedorID uint   `json:"proveedor_id"`
	Nombre      string `json:"nombre"`
	Direccion   string `json:"direccion"`
	Ciudad      string `json:"ciudad"`
}
type Cancha struct {
	ID               uint    `json:"cancha_id"`
	ComplejoID       uint    `json:"complejo_id"`
	Tipo             string  `json:"tipo"`
	Capacidad        int     `json:"capacidad"`
	PrecioReferencia float64 `json:"precio_referencia"`
}
type Solicitud struct {
	ID                uint            `json:"solicitud_id"`
	ClienteID         uint            `json:"cliente_id"`
	FechaEvento       time.Time       `json:"fecha_evento"`
	DuracionMinutos   int             `json:"duracion_minutos"`
	CantidadPersonas  int             `json:"cantidad_personas"`
	PresupuestoMaximo float64         `json:"presupuesto_maximo"`
	Estado            string          `json:"estado"`
	Items             []ItemSolicitud `json:"items,omitempty"`
	Ofertas           []Oferta        `json:"ofertas,omitempty"`
}
type ItemSolicitud struct {
	ID          uint   `json:"item_id"`
	SolicitudID uint   `json:"solicitud_id"`
	Tipo        string `json:"tipo"`
	Cantidad    int    `json:"cantidad"`
	Detalle     string `json:"detalle,omitempty"`
}
type Oferta struct {
	ID          uint      `json:"oferta_id"`
	SolicitudID uint      `json:"solicitud_id"`
	ComplejoID  uint      `json:"complejo_id"`
	CanchaID    uint      `json:"cancha_id"`
	MontoTotal  float64   `json:"monto_total"`
	Detalle     string    `json:"detalle,omitempty"`
	VenceEn     time.Time `json:"vence_en"`
	Estado      string    `json:"estado"`
}
type Reserva struct {
	ID                     uint      `json:"reserva_id"`
	OfertaID               uint      `json:"oferta_id"`
	Inicio                 time.Time `json:"inicio"`
	Fin                    time.Time `json:"fin"`
	MontoAcordado          float64   `json:"monto_acordado"`
	ComisionApp            float64   `json:"comision_app"`
	MontoProveedor         float64   `json:"monto_proveedor"`
	Estado                 string    `json:"estado"`
	EstadoLiquidacion      string    `json:"estado_liquidacion"`
	ComprobanteLiquidacion string    `json:"comprobante_liquidacion,omitempty"`
}
type Pago struct {
	ID             uint      `json:"pago_id"`
	OfertaID       uint      `json:"oferta_id"`
	Monto          float64   `json:"monto"`
	URLComprobante string    `json:"url_comprobante"`
	Estado         string    `json:"estado"`
	FechaPago      time.Time `json:"fecha_pago"`
}
