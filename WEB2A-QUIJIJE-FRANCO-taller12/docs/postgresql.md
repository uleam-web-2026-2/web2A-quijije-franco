# PostgreSQL en Windows 11

La aplicación usa `DATABASE_URL` y el driver Go `pgx`. No se ejecutan comandos destructivos: la migración usa `CREATE TABLE IF NOT EXISTS`.

```powershell
psql -U postgres -c "CREATE DATABASE reservas_canchas;"
$env:DATABASE_URL = "postgres://postgres:TU_CLAVE@localhost:5432/reservas_canchas"
$env:JWT_SECRET = "clave-local-de-demostracion"
$env:PORT = "8080"
go run .
```

Al iniciar con `DATABASE_URL`, [postgres.go](../internal/db/postgres.go) ejecuta [001_init.sql](../internal/db/migrations/001_init.sql). Para revisar las tablas:

```powershell
psql -U postgres -d reservas_canchas -c "\dt"
```

La migración crea: `usuarios`, `clientes`, `proveedores`, `complejos`, `canchas`, `solicitudes`, `item_solicitudes`, `ofertas`, `reservas` y `pagos`. Además crea una FK compuesta para impedir que una Oferta combine una cancha con un Complejo diferente. Con PostgreSQL, la petición 0 requiere una `oferta_id` existente porque la ficha exige la relación `Oferta 1:0..1 Reserva`.

Sin `DATABASE_URL`, el servidor usa memoria para la demostración y las pruebas sin modificar una base existente.
