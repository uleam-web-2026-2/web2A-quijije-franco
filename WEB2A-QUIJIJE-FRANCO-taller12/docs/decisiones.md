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

---

## Semana 12 · Seguridad de la API de reservas deportivas

### D1 · Lo que encontramos y cerramos

La auditoría sobre la base de código de Quijije y Franco registró **11 coincidencias** de los patrones indicados por el docente. De ellas, **8 son llamadas de ejecución SQL** (7 en el repositorio de reservas/pagos y 1 al ejecutar el script de migración administrativa); las otras 3 son `fmt.Sprintf` para JWT, la concatenación de la dirección HTTP y un `Group` de rutas. En las 7 consultas operativas se utilizan consultas fijas y parámetros `$1`, `$2`, etc. Ninguna arma SQL con el texto del cliente mediante `+` o `Sprintf`.

El dato `id` de la URL se convierte con `strconv.ParseUint` antes de la búsqueda; el servicio devuelve `400` si no es numérico. Aquí **no se agregó un filtro `estado`** que no existía en el negocio. Prueba prevista: `GET /reservas/1%20OR%201%3D1` con un JWT válido, junto con el GET de una reserva existente. Su resultado real queda pendiente de guardarse en `docs/taller12/inyeccion.txt`.

### D2 · El origen que puede leer nuestras respuestas

Se configuró `github.com/go-chi/cors` con origen permitido **`http://localhost:5173`**, lugar previsto del frontend local de reservas deportivas, sin `*`. El middleware se instala inmediatamente después de crear el router: la petición previa OPTIONS no debe pasar por validación de JWT. Los métodos autorizados incluyen `GET`, `POST`, `PUT`, `PATCH`, `DELETE` (la API actual usa GET/POST/PATCH), y las cabeceras `Authorization` y `Content-Type`.

Un navegador no debe poder leer la respuesta CORS desde `http://malo.example` porque no recibe `Access-Control-Allow-Origin`. **CORS no autoriza la operación en el servidor:** una petición PATCH hecha con una herramienta como curl desde otro origen puede ejecutarse si su JWT, rol, pago y transición son válidos. La prueba debe usar una nueva reserva pendiente, porque una ya confirmada devolvería `422`. Cuando se publique el frontend se cambiará el origen autorizado por su URL real.

### D3 · El límite del login

`POST /login` tiene un máximo de **5 intentos por minuto y dirección IP**, implementado con `github.com/go-chi/httprate`. Esta ventana frena intentos repetitivos de contraseña y no afecta consultas ni modificaciones de reservas. Se esperan cinco respuestas de credenciales incorrectas `401`; la sexta y la séptima deben ser `429` con `Retry-After`. **No se afirman resultados HTTP aún no ejecutados**: la evidencia real debe guardarse en `docs/taller12/rafaga.txt` y `docs/taller12/logs_seis.txt`.

### D4 · Detrás de un proxy

`ClientIPFromRemoteAddr` obtiene la IP de la conexión TCP y no confía sin validación en `X-Forwarded-For`. Si en producción la API recibe todo el tráfico desde un proxy, varios usuarios podrían compartir el mismo contador y bloquearse entre sí. Para distinguirlos habrá que configurar primero un proxy de confianza que controle las cabeceras de IP; no sería seguro usar una cabecera arbitraria enviada por el cliente.

## Semana 12 · Reto: mínimo privilegio en PostgreSQL

### D5 · Con qué usuario se conecta la API

Se propone el rol operativo **`reservas_api`** sin `SUPERUSER`, `CREATEDB`, `CREATEROLE` ni propiedad de las tablas. El script `scripts/semana12/crear_rol_api.sql` concede solo las operaciones utilizadas por `internal/reservas/repositorio_sql.go`:

| Objeto | Permisos | Por qué |
|---|---|---|
| `reservas` | SELECT, INSERT, UPDATE | Listar, obtener, crear y cambiar estado |
| `pagos` | SELECT, INSERT | Comprobar y registrar un pago |
| `ofertas`, `solicitudes`, `clientes` | SELECT | Verificar pertenencia de una reserva al cliente |
| `reservas_reserva_id_seq`, `pagos_pago_id_seq` | USAGE | Obtener los identificadores `BIGSERIAL` al insertar |
| `usuarios`, `proveedores`, `complejos`, `canchas`, `item_solicitudes` | Ningún permiso de tabla para esta API | No hay consultas ni mutaciones de esas tablas en las rutas implementadas |

El intento de `DROP TABLE public.reservas` debe fallar por falta de propiedad; **el resultado aún no está verificado en una instancia PostgreSQL**. Una vez ejecutado se registrará en `reto_drop.txt`; `reto_usuarios.txt` mostrará los roles reales.

### D6 · Quién crea y cambia las tablas

Se separaron `db.Open` (abrir y comprobar conexión) y `db.Migrate` (aplicar `001_init.sql`). En el arranque normal no se ejecutan migraciones; `RUN_MIGRATIONS=true` solo activa migraciones administrativas y sale **sin iniciar el servidor**. Debe usarse temporalmente una conexión con un usuario administrador de la base y restaurar después `DATABASE_URL` con el usuario `reservas_api`. Si se altera un modelo y se arranca con el usuario limitado, la estructura no se cambia automáticamente: es necesario ejecutar la migración administrativa por separado.

### D7 · Lo que esto no frena

Incluso sin poder hacer `DROP TABLE`, un atacante que lograra ejecutar una instrucción SQL con los privilegios de `reservas_api` podría consultar reservas o pagos y modificar estados de las reservas, siempre dentro de los privilegios otorgados. Los privilegios mínimos reducen impacto; no reemplazan sentencias parametrizadas, control de acceso, autenticación y validación de transiciones.

**Estado de las pruebas:** código, permisos y análisis preparados a partir del repositorio. No se atribuyen resultados de ejecución ni capturas de otro proyecto. Los archivos `.txt` que requieren PostgreSQL o un servidor real solo se incluirán tras ejecutar las pruebas indicadas en `docs/taller12/GUIA_DE_PRUEBAS.md`.
