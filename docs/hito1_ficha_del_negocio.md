# Hito 1 · Ficha del negocio

**Pareja:** John Franco · Jose Quijije
**Paralelo:** Aplicación para el Servidor Web A / Aplicaciones Web II A
**Negocio en una línea:** (qué vende, a quién)


## 1. Negocio de referencia

El negocio de referencia es una plataforma digital de reservas y organización de partidos y eventos deportivos. Opera como un intermediario tecnológico que permite a los usuarios reservar canchas sintéticas, o solicitar servicios complementarios (arbitraje, indumentaria, bebidas). El modelo de ingresos se basa en el cobro de una comisión por cada transacción o reserva completada exitosamente a través de la plataforma, tanto al usuario final como al complejo deportivo por la gestión de la demanda. Plataformas similares en mercados internacionales han demostrado una alta tracción, reportando ingresos mensuales recurrentes superiores a los $50,000 dólares gracias a la automatización del proceso de reservas y a la fidelización de complejos deportivos que buscan maximizar el uso de sus instalaciones durante horas de baja demanda.

**Enlace:** (https://sportskey.com/)

## 2. Caso de contraste

Un caso de contraste común son las aplicaciones locales y regionales que digitalizaron el alquiler de canchas ofreciendo un catálogo estático y cerrado, operando bajo un modelo de horario fijo sin flexibilidad ni servicios adicionales. Estas iniciativas se estancaron o enfrentaron resistencia debido a una diferencia profunda en el modelo de costos y mercado: ignoraron por completo los ingresos clave derivados del consumo complementario (bebidas, arbitraje, eventos) y no ofrecieron dinamismo frente a la alta variabilidad de la demanda local. Al ser un sistema rígido donde el complejo ponía un precio fijo sin tomar en cuenta la negociación ni la personalización del servicio, los usuarios preferían seguir contactando a los dueños por WhatsApp o llamadas telefónicas para regatear precios y armar paquetes a su medida, dejando la aplicación tradicional totalmente obsoleta.

**Fuente:** (https://www.icancha.com/)

## 3. Adaptación al Ecuador

1. Fuerte cultura del regateo y negociación: En Ecuador, los clientes están acostumbrados a negociar directamente el precio de los servicios o eventos grupales.
Efecto: Los usuarios intentarían lanzar ofertas con precios absurdamente bajos que no cubren los costos operativos de los complejos deportivos.

2. Alta preferencia por transferencias bancarias directas: Aunque existen pasarelas de pago, los complejos y usuarios prefieren el pago por transferencia (P2P) para evitar comisiones bancarias altas.
Efecto: El sistema no puede automatizar el cobro instantáneo mediante tarjeta sin fricción, requiriendo un estado intermedio de validación manual del comprobante de pago.

3. Alto índice de "no-shows" (clientes que reservan y no se presentan): Existe poca penalización cultural por faltar a una reserva verbal o digital no pagada.
Efecto: Exigencia obligatoria de un abono parcial o validación estricta de la oferta aceptada antes de bloquear el horario en el complejo.

**Qué cambió en el modelo por estas restricciones:** 
Qué cambió en el modelo por estas restricciones: Se incorporó la entidad CatalogoBase con un atributo de precio_minimo, y se implementó una regla de validación estricta en el servidor backend que rechaza automáticamente cualquier Solicitud cuyo presupuesto_cliente esté por debajo del costo operativo base establecido por los proveedores.

## 4. Modelo de datos

### Entidad: Usuario

| Atributo   | Tipo          | Obligatorio | Ejemplo              |
|------------|---------------|-------------|----------------------|
| usuario_id | número entero | si          | 1                    |
| nombre     | texto         | si          | 'Carlos Perez'       |
| correo     | texto         | si          | 'carlos@gmail.com'   |
| password   | texto         | si          | '$2a$10$Xk3...'      |

### Entidad: Cliente

| Atributo   | Tipo          | Obligatorio | Ejemplo        |
|------------|---------------|-------------|----------------|
| cliente_id | número entero | si          | 1              |
| usuario_id | número entero | si          | 1              |
| cedula     | texto         | si          | '1311682155'   |
| telefono   | texto         | si          | '0983660453'   |

### Entidad: Proveedor

| Atributo       | Tipo          | Obligatorio | Ejemplo           |
|----------------|---------------|-------------|-------------------|
| proveedor_id   | número entero | si          | 1                 |
| usuario_id     | número entero | si          | 2                 |
| ruc            | texto         | si          | '1391234567001'   |
| telefono       | texto         | si          | '0963560769'      |
| banco          | texto         | si          | 'Banco Pichincha' |
| numero_cuenta  | texto         | si          | '2200123456'      |
| titular_cuenta | texto         | si          | 'Juan Mendoza'    |

### Entidad: Complejo

| Atributo     | Tipo          | Obligatorio | Ejemplo                  |
|--------------|---------------|-------------|--------------------------|
| complejo_id  | número entero | si          | 1                        |
| proveedor_id | número entero | si          | 1                        |
| nombre       | texto         | si          | 'Canchas El Golazo'      |
| direccion    | texto         | si          | 'Av. 4 de Noviembre'     |
| ciudad       | texto         | si          | 'Manta'                  |

### Entidad: Cancha

| Atributo          | Tipo          | Obligatorio | Ejemplo  |
|-------------------|---------------|-------------|----------|
| cancha_id         | número entero | si          | 1        |
| complejo_id       | número entero | si          | 1        |
| tipo              | texto         | si          | 'futbol' |
| capacidad         | número entero | si          | 14       |
| precio_referencia | decimal       | si          | 25.00    |

### Entidad: Solicitud

| Atributo           | Tipo          | Obligatorio | Ejemplo            |
|--------------------|---------------|-------------|--------------------|
| solicitud_id       | número entero | si          | 1                  |
| cliente_id         | número entero | si          | 1                  |
| fecha_evento       | fecha y hora  | si          | '2026-10-15 19:00' |
| duracion_minutos   | número entero | si          | 90                 |
| cantidad_personas  | número entero | si          | 12                 |
| presupuesto_maximo | decimal       | si          | 40.00              |
| estado             | texto         | si          | 'abierta'          |

### Entidad: Item_Solicitud

| Atributo     | Tipo          | Obligatorio | Ejemplo        |
|--------------|---------------|-------------|----------------|
| item_id      | número entero | si          | 1              |
| solicitud_id | número entero | si          | 1              |
| tipo         | texto         | si          | 'bebidas'      |
| cantidad     | número entero | si          | 12             |
| detalle      | texto         | no          | 'Agua sin gas' |

### Entidad: Oferta

| Atributo     | Tipo          | Obligatorio | Ejemplo            |
|--------------|---------------|-------------|--------------------|
| oferta_id    | número entero | si          | 1                  |
| solicitud_id | número entero | si          | 1                  |
| complejo_id  | número entero | si          | 1                  |
| cancha_id    | número entero | si          | 1                  |
| monto_total  | decimal       | si          | 35.00              |
| detalle      | texto         | no          | 'Incluye bebidas'  |
| vence_en     | fecha y hora  | si          | '2026-10-10 23:59' |
| estado       | texto         | si          | 'pendiente'        |

### Entidad: Reserva

| Atributo           | Tipo          | Obligatorio | Ejemplo            |
|--------------------|---------------|-------------|--------------------|
| reserva_id         | número entero | si          | 1                  |
| oferta_id          | número entero | si          | 1                  |
| inicio             | fecha y hora  | si          | '2026-10-15 19:00' |
| fin                | fecha y hora  | si          | '2026-10-15 20:30' |
| monto_acordado     | decimal       | si          | 35.00              |
| comision_app       | decimal       | si          | 3.50               |
| estado_comision    | texto         | si          | 'pendiente'        |
| estado             | texto         | si          | 'confirmada'       |

### Entidad: Pago

| Atributo        | Tipo          | Obligatorio | Ejemplo                       |
|-----------------|---------------|-------------|-------------------------------|
| pago_id         | número entero | si          | 1                             |
| oferta_id       | número entero | si          | 1                             |
| monto           | decimal       | si          | 35.00                         |
| url_comprobante | texto         | si          | 'https://.../comprobante.jpg' |
| estado          | texto         | si          | 'pendiente'                   |
| fecha_pago      | fecha y hora  | si          | '2026-10-11 10:30'            |


### Relaciones

| Entidades            | Cardinalidad | Frase                                                  |
|----------------------|--------------|--------------------------------------------------------|
| Usuario - Cliente    | 1:0..1       | Un usuario puede tener un perfil de cliente            |
| Usuario - Proveedor  | 1:0..1       | Un usuario puede tener un perfil de proveedor          |
| Proveedor - Complejo | 1:N          | Un proveedor administra uno o varios complejos         |
| Complejo - Cancha    | 1:N          | Un complejo tiene una o varias canchas                 |
| Cliente - Solicitud  | 1:N          | Un cliente crea cero o varias solicitudes              |
| Solicitud - Item_Solicitud | 1:N    | Una solicitud incluye cero o varios ítems adicionales  |
| Solicitud - Oferta   | 1:N          | Una solicitud recibe cero o varias ofertas             |
| Complejo - Oferta    | 1:N          | Un complejo envía cero o varias ofertas                |
| Cancha - Oferta      | 1:N          | Una cancha se ofrece en cero o varias ofertas          |
| Oferta - Reserva     | 1:0..1       | Una oferta aceptada origina una sola reserva           |
| Oferta - Pago        | 1:0..1       | Una oferta aceptada se paga con un solo pago           |


### Structs en Go

```go
//package models

//import "time"

type Usuario struct {
    ID       uint   `gorm:"primaryKey" json:"usuario_id"`
    Nombre   string `gorm:"not null" json:"nombre"`
    Correo   string `gorm:"not null;uniqueIndex" json:"correo"`
    Password string `gorm:"not null" json:"-"` // guarda el hash, nunca se devuelve en el JSON
}

type Cliente struct {
    ID        uint   `gorm:"primaryKey" json:"cliente_id"`
    UsuarioID uint   `gorm:"not null;uniqueIndex" json:"usuario_id"`
    Cedula    string `gorm:"not null;uniqueIndex" json:"cedula"`
    Telefono  string `gorm:"not null" json:"telefono"`
}

type Proveedor struct {
    ID            uint   `gorm:"primaryKey" json:"proveedor_id"`
    UsuarioID     uint   `gorm:"not null;uniqueIndex" json:"usuario_id"`
    RUC           string `gorm:"not null;uniqueIndex" json:"ruc"`
    Telefono      string `gorm:"not null" json:"telefono"`
    Banco         string `gorm:"not null" json:"banco"`
    NumeroCuenta  string `gorm:"not null" json:"numero_cuenta"`
    TitularCuenta string `gorm:"not null" json:"titular_cuenta"`
}

type Complejo struct {
    ID          uint   `gorm:"primaryKey" json:"complejo_id"`
    ProveedorID uint   `gorm:"not null;index" json:"proveedor_id"`
    Nombre      string `gorm:"not null" json:"nombre"`
    Direccion   string `gorm:"not null" json:"direccion"`
    Ciudad      string `gorm:"not null" json:"ciudad"`
}

type Cancha struct {
    ID               uint    `gorm:"primaryKey" json:"cancha_id"`
    ComplejoID       uint    `gorm:"not null;index" json:"complejo_id"`
    Tipo             string  `gorm:"not null" json:"tipo"` // futbol|basquet|padel
    Capacidad        int     `gorm:"not null" json:"capacidad"`
    PrecioReferencia float64 `gorm:"type:numeric(10,2);not null" json:"precio_referencia"`
}

type Solicitud struct {
    ID                uint      `gorm:"primaryKey" json:"solicitud_id"`
    ClienteID         uint      `gorm:"not null;index" json:"cliente_id"`
    FechaEvento       time.Time `gorm:"not null" json:"fecha_evento"`
    DuracionMinutos   int       `gorm:"not null" json:"duracion_minutos"`
    CantidadPersonas  int       `gorm:"not null" json:"cantidad_personas"`
    PresupuestoMaximo float64   `gorm:"type:numeric(10,2);not null" json:"presupuesto_maximo"`
    Estado            string    `gorm:"not null;default:abierta" json:"estado"` // abierta|adjudicada|cancelada
}

type ItemSolicitud struct {
    ID          uint   `gorm:"primaryKey" json:"item_id"`
    SolicitudID uint   `gorm:"not null;index" json:"solicitud_id"`
    Tipo        string `gorm:"not null" json:"tipo"` // arbitraje|bebidas|indumentaria|otro
    Cantidad    int    `gorm:"not null" json:"cantidad"`
    Detalle     string `json:"detalle"`
}

type Oferta struct {
    ID          uint      `gorm:"primaryKey" json:"oferta_id"`
    SolicitudID uint      `gorm:"not null;index" json:"solicitud_id"`
    ComplejoID  uint      `gorm:"not null;index" json:"complejo_id"`
    CanchaID    uint      `gorm:"not null;index" json:"cancha_id"`
    MontoTotal  float64   `gorm:"type:numeric(10,2);not null" json:"monto_total"`
    Detalle     string    `json:"detalle"`
    VenceEn     time.Time `gorm:"not null" json:"vence_en"`
    Estado      string    `gorm:"not null;default:pendiente" json:"estado"` // pendiente|aceptada|rechazada|vencida
}

type Reserva struct {
    ID                    uint      `gorm:"primaryKey" json:"reserva_id"`
    OfertaID              uint      `gorm:"not null;uniqueIndex" json:"oferta_id"`
    Inicio                time.Time `gorm:"not null" json:"inicio"`
    Fin                   time.Time `gorm:"not null" json:"fin"`
    MontoAcordado         float64   `gorm:"type:numeric(10,2);not null" json:"monto_acordado"`
    ComisionApp           float64   `gorm:"type:numeric(10,2);not null" json:"comision_app"`
    MontoProveedor        float64   `gorm:"type:numeric(10,2);not null" json:"monto_proveedor"`
    Estado                string    `gorm:"not null;default:pendiente" json:"estado"` // pendiente|confirmada|cancelada|completada
    EstadoLiquidacion     string    `gorm:"not null;default:pendiente" json:"estado_liquidacion"` // pendiente|liquidada
    ComprobanteLiquidacion string   `json:"comprobante_liquidacion"`
}

type Pago struct {
    ID             uint      `gorm:"primaryKey" json:"pago_id"`
    OfertaID       uint      `gorm:"not null;uniqueIndex" json:"oferta_id"`
    Monto          float64   `gorm:"type:numeric(10,2);not null" json:"monto"`
    URLComprobante string    `gorm:"not null" json:"url_comprobante"`
    Estado         string    `gorm:"not null;default:pendiente" json:"estado"` // pendiente|verificado|rechazado
    FechaPago      time.Time `gorm:"not null" json:"fecha_pago"`
}
```


Fechas como `time.Time` y no como texto:campos como `FechaEvento`, `Inicio`, `Fin`, `VenceEn` y `FechaPago` se usan para comparar y ordenar (por ejemplo, saber si una oferta ya venció o si dos reservas se cruzan). Con `time.Time` GORM los guarda como `timestamp` en PostgreSQL y esas comparaciones se hacen directo en la consulta.

### Diagrama del modelo completo

erDiagram
    USUARIO ||--o| CLIENTE : tiene_perfil
    USUARIO ||--o| PROVEEDOR : tiene_perfil
    PROVEEDOR ||--|{ COMPLEJO : administra
    COMPLEJO ||--|{ CANCHA : tiene
    CLIENTE ||--o{ SOLICITUD : crea
    SOLICITUD ||--o{ ITEM_SOLICITUD : incluye
    SOLICITUD ||--o{ OFERTA : recibe
    COMPLEJO ||--o{ OFERTA : envia
    CANCHA ||--o{ OFERTA : se_ofrece_en
    OFERTA ||--o| RESERVA : origina
    OFERTA ||--o| PAGO : se_paga_con

    USUARIO {
        int id_usuario PK
        string nombre
        string email UK
        string telefono
        string password_hash
    }

    CLIENTE {
        int id_cliente PK
        int id_usuario FK, UK
    }

    PROVEEDOR {
        int id_proveedor PK
        int id_usuario FK, UK
        string nombre_comercial
    }

    COMPLEJO {
        int id_complejo PK
        int id_proveedor FK
        string nombre
        string direccion
        string ciudad
    }

    CANCHA {
        int id_cancha PK
        int id_complejo FK
        string tipo "futbol|basquet|padel|evento"
        string nombre
        int capacidad
    }

    SOLICITUD {
        int id_solicitud PK
        int id_cliente FK
        string tipo_servicio
        datetime fecha_evento
        int duracion_minutos
        int cantidad_personas
        decimal presupuesto_maximo
        string estado "abierta|adjudicada|cancelada"
        text descripcion
    }

    ITEM_SOLICITUD {
        int id_item_solicitud PK
        int id_solicitud FK
        string nombre
        int cantidad
        text detalle
    }

    OFERTA {
        int id_oferta PK
        int id_solicitud FK
        int id_complejo FK
        int id_cancha FK
        decimal monto_total
        text detalle
        datetime fecha_oferta
        string estado "pendiente|aceptada|rechazada|vencida"
    }

    RESERVA {
        int id_reserva PK
        int id_oferta FK, UK
        datetime inicio
        datetime fin
        string estado "pendiente|confirmada|cancelada|completada"
    }

    PAGO {
        int id_pago PK
        int id_oferta FK, UK
        decimal monto
        string metodo "tarjeta|transferencia|efectivo"
        datetime fecha_pago
        string estado "pendiente|aprobado|rechazado|reembolsado"
        string referencia_transaccion
    }
``

Decisiones discutibles del modelo y por qué las tomamos
Perfiles separados: CLIENTE y PROVEEDOR
Podría haberse usado un solo campo rol en USUARIO. Se separaron porque cada perfil representa una función de negocio distinta y permite que una misma cuenta tenga perfil de cliente, proveedor o ambos, sin duplicar los datos de acceso.

COMPLEJO y CANCHA como entidades distintas
Un complejo puede tener varias canchas o espacios con características diferentes. Separarlas permite ofertar una cancha concreta —por ejemplo, una de fútbol 5 o una de pádel— dentro del mismo complejo.

La oferta referencia tanto a COMPLEJO como a CANCHA
Técnicamente, la cancha ya permite conocer el complejo. Mantener ambos identificadores en OFERTA facilita consultas frecuentes y conserva explícitamente el complejo que realizó la propuesta. Debe validarse que la cancha pertenezca al complejo indicado.

Los ítems adicionales pertenecen a la solicitud
ITEM_SOLICITUD registra lo que el cliente necesita, como bebidas o indumentaria. No se creó una tabla de catálogo de productos para mantener el modelo esencial; si luego se requiere inventario, precios o disponibilidad por proveedor, convendría agregarla.

Reserva y pago se relacionan directamente con la oferta
La regla establece que una oferta aceptada genera, como máximo, una reserva y un pago. Por eso ambos tienen id_oferta único. Una alternativa más habitual sería asociar PAGO a RESERVA; sería preferible si se admitirán pagos parciales, cuotas, anticipos o reembolsos múltiples.

Una sola cancha por oferta
Se asumió que cada oferta reserva un único espacio. Si un campeonato necesita varias canchas al mismo tiempo, habría que reemplazar id_cancha en OFERTA por una tabla intermedia, por ejemplo OFERTA_CANCHA.

Un pago por oferta
Esta es una simplificación para el flujo inicial. El índice único en PAGO.id_oferta evita más de un pago. Si la plataforma manejará señas y saldos, debe eliminarse esa restricción y modelar varios pagos por reserva.

No se modeló disponibilidad de canchas inicialmente
El modelo se concentra en solicitudes, ofertas, reservas y pagos. Para evitar que un proveedor oferte una cancha ya ocupada, una siguiente versión debería incorporar horarios bloqueados o una agenda de disponibilidad asociada a CANCHA.

## 5. Máquina de estados

**Entidad con estados:** Reserva

| Estado | Qué significa |
|--------|---------------|
| pendiente (inicial) | La reserva se creó al aceptar la oferta, pero el cliente aún no ha pagado o el pago no se ha verificado |
| confirmada | El administrador verificó el comprobante del pago; el horario queda asegurado para el cliente |
| completada | El evento ya se realizó; desde aquí se puede liquidar al proveedor |
| cancelada | La reserva se anuló antes de realizarse; no se liquida nada |

| De         | A          |  Quién la hace | Condición |
|--------|---|------------|-----------              |
| pendiente  | confirmada | Administrador           | El pago asociado a la oferta pasó a `verificado` |
| pendiente  | cancelada  | Cliente o administrador | El cliente desiste, o el comprobante fue `rechazado` |
| confirmada | completada | Sistema | Ya pasó la hora `fin` de la reserva |
| confirmada | cancelada  | Cliente o administrador | Aún no llega la hora `inicio`; si hubo pago verificado, se reembolsa |

**Transición prohibida y por qué:** completada → cancelada. Una vez completada, el servicio ya se prestó y la comisión de la app quedó calculada para liquidar al proveedor. Si se pudiera cancelar después, habría que revertir una liquidación que quizá ya se transfirió, y el cliente podría pedir reembolso de un servicio que sí recibió. Tampoco se permite cancelada → confirmada ni cancelada → completada: una reserva cancelada es definitiva, y si el cliente quiere retomarla debe crear una nueva solicitud.

### Diagrama de estados

stateDiagram-v2
    [*] --> pendiente
    
    pendiente --> confirmada: Administrador verifica pago
    pendiente --> cancelada: Cliente o administrador rechaza pago
    
    confirmada --> completada: Sistema (pasó hora fin)
    confirmada --> cancelada: Cliente o administrador (antes de hora inicio)
    
    completada --> [*]
    cancelada --> [*]
    
    note right of pendiente
        Oferta aceptada
        Pago pendiente o no verificado
    end note
    
    note right of confirmada
        Pago verificado
        Horario asegurado
    end note
    
    note right of completada
        Evento realizado
        Listo para liquidar proveedor
    end note
    
    note right of cancelada
        Reserva anulada
        No se liquida nada
        Definitiva
    end note



## 6. Roles y permisos

| Acción | Cliente | Proveedor | Administrador |
|--------|---------|-----------|---------------|
| Registrarse y editar su perfil | solo los suyos | solo los suyos | solo los suyos |
| Crear solicitud | sí | no | no |
| Ver solicitudes | solo los suyos | todos (solo las abiertas) | todos |
| Cancelar solicitud | solo los suyos | no | todos |
| Registrar complejos y canchas | no | solo los suyos | no |
| Enviar oferta | no | solo los suyos | no |
| Ver ofertas recibidas | solo los suyos | solo los suyos | todos |
| Aceptar o rechazar oferta | solo los suyos | no | no |
| Subir comprobante de pago | solo los suyos | no | no |
| Verificar o rechazar pago | no | no | todos |
| Ver reservas | solo los suyos | solo los suyos | todos |
| Cancelar reserva | solo los suyos (si está pendiente o confirmada) | no | todos |
| Registrar liquidación al proveedor | no | no | todos |
| Ver comisión y monto a recibir | no | solo los suyos | todos |

## 7. Mapa de endpoints por rol

## 7. Mapa de endpoints por rol

| Endpoint | Rol que lo llama | Pantalla que lo consume | Qué devuelve | Qué valida | Código si falla |
|----------|------------------|-------------------------|--------------|------------|-----------------|
| `POST /solicitudes` | Cliente | Nueva solicitud | Solicitud creada | Rol cliente; fecha futura; montos > 0 | 401, 403, 422 |
| `GET /solicitudes` | Cliente, proveedor, admin | Mis solicitudes / Solicitudes abiertas | Lista según el rol | Token válido | 401 |
| `GET /solicitudes/{id}` | Cliente, proveedor, admin | Detalle de solicitud | Solicitud con ítems y ofertas | Que exista y le corresponda | 401, 403, 404 |
| `PUT /solicitudes/{id}` | Cliente (dueño) | Editar solicitud | Solicitud actualizada | Que sea suya y esté abierta | 403, 404, 409, 422 |
| `DELETE /solicitudes/{id}` | Cliente (dueño), admin | Botón Cancelar | 204 (pasa a `cancelada`) | Que sea suya y no esté adjudicada | 403, 404, 409 |
| `POST /solicitudes/{id}/ofertas` | Proveedor | Enviar oferta | Oferta `pendiente` | Solicitud abierta; cancha suya; monto > 0 | 403, 404, 409, 422 |

### Matriz pantalla × endpoint

<!-- Marquen con X la celda donde la pantalla usa ese endpoint. Sirve para ver dos errores:
     una pantalla sin datos y un endpoint que no usa nadie. -->

| Pantalla             | [GET /recurso] | [POST /solicitudes] | [PATCH /solicitudes/{id}/estado] | [GET /solicitudes/{id}] |
|----------------------|----------------|-----------------|------------------------------|---------------------|
|Nueva solicitud       |                |       x         |                              |                     |
|Mis solicitudes       |       x        |                 |                              |                     |
|Detalle de solicitud  |                |                 |               x              |            x        |

**Endpoints que ya están funcionando y en qué archivo:** 

## 8. Declaración de IA

Se usó Claude (Anthropic) como asistente de apoyo en este hito. Las decisiones de diseño fueron nuestras, y la IA ayudó a revisar, redactar y dar formato.

| Sección | Para qué se usó |
|---------|-----------------|
| Diagrama ER y tablas de entidades | Revisar el diagrama (sintaxis y consistencia) y pasar las entidades al formato de tablas |
| Relaciones | Redactar la tabla de cardinalidades y sus frases |
| Structs en Go | Generar los structs con etiquetas de GORM y JSON, y explicar la decisión de tipos |
| Máquina de estados | Proponer estados y transiciones de `Reserva` y el diagrama en Mermaid |
| Roles y permisos | Armar la tabla de permisos por rol |
| Mapa de endpoints y matriz pantalla × endpoint | Redactar las tablas y verificar que no hubiera pantallas sin datos ni endpoints sin uso |

**Decisiones que tomamos nosotros:** el modelo de subasta inversa, tablas separadas de Usuario y Proveedor, un solo pago por oferta, y que la comisión se le cobre a los proveedores y pase primero por la cuenta de la app.

**Cómo lo verificamos:** [completa: por ejemplo, "revisamos cada sección entre los dos y probamos los endpoints en el backend"].s