# Semana 12 · Auditoría de inyección de SQL

Proyecto **Quijije y Franco: subasta inversa de canchas deportivas**.

Se aplicaron los 12 patrones de `patrones.txt` a archivos `.go` no-test del proyecto actualizado. `auditoria.txt` registra las **11 coincidencias reales**: 8 puntos relacionados con ejecución de SQL y 3 coincidencias que no generan SQL. No se encontró SQL construido mediante concatenación de datos de peticiones HTTP.

| Archivo y línea | ¿Arma SQL? | ¿Llega texto del cliente? | ¿Va aparte? | Qué se hizo |
|---|---|---|---|---|
| `internal/db/postgres.go:38` (`ExecContext`) | Sí, DDL de inicialización | No; carga el SQL **fijo** de `migrations/001_init.sql` | No aplica | Se separó la migración del arranque normal: solo con `RUN_MIGRATIONS=true` y cuenta administradora. |
| `internal/reservas/repositorio_sql.go:14` (`QueryContext`) | Sí, SELECT de reservas | No; listado sin filtros del cliente | No aplica | Se deja consulta constante, sin interpolar valores. |
| `internal/reservas/repositorio_sql.go:31` (`QueryRowContext`) | Sí, SELECT por ID | Sí, el `{id}` de `/reservas/{id}` | Sí: `$1`; antes se valida con `strconv.ParseUint` | Se mantiene parametrizado. |
| `internal/reservas/repositorio_sql.go:37` (`QueryRowContext`) | Sí, INSERT Reserva | Sí, los valores de POST `/reservas` | Sí: `$1` a `$9` | Se mantiene parametrizado. |
| `internal/reservas/repositorio_sql.go:40` (`ExecContext`) | Sí, UPDATE estado | Sí, `id`, `estado`, comprobante | Sí: `$1` a `$3` | Se mantiene parametrizado. |
| `internal/reservas/repositorio_sql.go:46` (`QueryRowContext`) | Sí, SELECT de propiedad | Sí, identificadores de reserva/usuario | Sí: `$1` y `$2` | Se mantiene parametrizado. |
| `internal/reservas/repositorio_sql.go:51` (`QueryRowContext`) | Sí, SELECT Pago por Oferta | Sí, ID de oferta recibido al crear reserva | Sí: `$1` | Se mantiene parametrizado. |
| `internal/reservas/repositorio_sql.go:56` (`QueryRowContext`) | Sí, INSERT Pago | Sí, JSON de la reserva | Sí: `$1` a `$5` | Se mantiene parametrizado. |
| `internal/auth/auth.go:30` (`fmt.Sprintf`) | No; formatea el ID aleatorio del token JWT | No | No aplica | Sin modificación: no es SQL. |
| `main.go:48` (`"+port`) | No; construye una dirección HTTP | No | No aplica | No es una sentencia SQL. |
| `main.go:83` (`r.Group`) | No; registra rutas de chi | No | No aplica | No es `GROUP BY` de SQL. |

**Por qué no hay un filtro `?estado=` en esta API:** el proyecto original solo publica `GET /reservas` y `GET /reservas/{id}`. Añadir un filtro sería cambiar el comportamiento del negocio fuera del alcance de este taller. El ataque equivalente se envía al segmento `{id}` con `1 OR 1=1`: el handler tiene que devolver **400** gracias a `strconv.ParseUint` y no permitir que ese texto llegue a `$1`.

**Evidencia HTTP real:** `inyeccion.txt` se genera desde el servidor con `python scripts/semana12/generar_evidencias.py`. Su ausencia significa que aún no se ejecutó esa verificación; no se escribió ninguna salida simulada.
