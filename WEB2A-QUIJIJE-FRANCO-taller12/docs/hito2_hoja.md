# Hito 2 · Hoja de demostración

## 1 · Usuarios de ejemplo

| Usuario | Rol | Contraseña |
|---|---|---|
| admin@reservas.local | administrador | admin-2026 |
| cliente@reservas.local | cliente | cliente-2026 |

## 2 · Las seis peticiones

1. `GET /reservas` sin token → `401`.
2. `POST /login` con admin → `200` y token JWT.
3. `GET /reservas` con token → `200`.
4. `PATCH /reservas/{id}/estado` con cliente → `403`.
5. La misma petición con admin y `confirmada` → `200`.
6. La misma reserva, transición prohibida → `422`.

La petición 0 (`POST /reservas`) crea el registro de partida. La demostración usa la entidad Reserva y su flujo de subasta inversa queda representado por `oferta_id`; los endpoints de Solicitud–Oferta no implementados están marcados `(plan)` en el contrato.

La petición 0 obtiene ambos JWT mediante un script de preparación y crea el Pago en estado `verificado`; por eso la petición 5 puede confirmar la Reserva respetando la regla del negocio.

## 3 · Matriz ruta × rol

| Ruta | Sin sesión | Cliente | Admin |
|---|---|---|---|
| POST /login | Sí | Sí | Sí |
| GET /reservas | 401 | 200 | 200 |
| POST /reservas | 401 | 201 | 201 |
| PATCH /reservas/{id}/estado | 401 | 403 | 200/422 |

## 4 · Recuadro de cumplimiento

- [x] Structs y migración de las diez entidades.
- [ ] Validador Redocly: ejecutar en un entorno con Node/npm.
- [x] Seis peticiones y petición 0 documentadas.
- [x] `go test ./...` y prueba de integración.
- [x] Capas de Reserva explicadas en `docs/decisiones.md`.
- [ ] Capturas de respaldo: pendientes de generarse desde la demostración real.
