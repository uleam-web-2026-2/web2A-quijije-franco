# Decisiones del Hito 2

## Correcciones respecto al Hito 1

| Archivo | Cambio | Justificación |
|---|---|---|
| `internal/modelos/modelos.go` | Se añadieron los diez structs y `Usuario.Rol`. | La ficha no incluía rol en el struct, pero la rúbrica exige permisos; el rol pertenece a la cuenta y no agrega una entidad de negocio. |
| `internal/db/migrations/001_init.sql` | Se crearon las diez tablas con FK, unicidad y checks. | La ficha exige relaciones y estados persistidos; `IF NOT EXISTS` evita borrar datos. |
| `internal/reservas/estado.go` | Se centralizó la máquina de estados de Reserva. | La regla no depende de HTTP y una transición prohibida devuelve 422 en el handler. |
| `main.go` | Se reemplazó la demostración de tickets por Reserva y se dejaron de registrar `/tickets` y `/explotar`. | El negocio oficial es la subasta inversa deportiva, no mesa de ayuda. El código heredado queda sin rutas activas. |
| `docs/openapi.yaml` | Las entidades sin endpoint llevan `(plan)`. | La rúbrica pide contrato completo sin programar CRUD prematuro. |

## Semana 8 · estructura

La entidad principal es Reserva y está dividida en:

- `internal/httpapi/`: decodifica HTTP, valida forma, traduce errores y responde JSON.
- `internal/reservas/servicio.go`: caso de uso y coordinación del repositorio.
- `internal/reservas/estado.go`: regla pura de transiciones, sin `net/http` ni base de datos.
- `internal/reservas/repositorio.go`: interfaz y repositorio en memoria para pruebas/demostración.
- `internal/db/`: conexión opcional a PostgreSQL y migración idempotente.

Así el handler no consulta la base y la regla de negocio puede probarse sin servidor. Cambiar el repositorio no obliga a cambiar HTTP ni la máquina de estados.

## Decisiones de negocio

La comisión se calcula al crear la reserva como 10% del monto acordado y el proveedor recibe el 90%, conforme a la ficha. La reserva inicia `pendiente`, pasa a `confirmada` solo mediante el rol administrador y un Pago asociado en estado `verificado`. Pago inexistente, pendiente o rechazado produce HTTP 422. La cancelación se permite al administrador o al cliente propietario, respetando la hora de inicio; el completado se expone como operación interna del sistema cuando ya pasó `fin`.

`CompletarVencidas` contiene la operación de sistema para completar reservas confirmadas cuyo `fin` ya pasó. El planificador periódico que la invoque queda fuera del alcance del Hito 2; no se permite que un usuario fuerce manualmente la transición con un token normal.

## Actualización de la ficha para Hito 2

La ficha conserva las diez entidades y sus relaciones. Se documenta aquí, y debe mantenerse también en la sección 7 de la ficha, que Reserva es la entidad principal demostrada porque es la entidad con máquina de estados; Solicitud y Oferta permanecen como flujo de negocio y sus endpoints no implementados están marcados `(plan)`. Esta elección no elimina Solicitud, Oferta ni Pago.

El cambio de Reserva a `confirmada` no es solo una autorización de rol: el servicio consulta Pago por `oferta_id` y exige `estado = verificado`. Las rutas heredadas `/tickets` y `/explotar` se retiraron de la API de reservas para no mezclar el ejemplo de mesa de ayuda con el negocio deportivo.
