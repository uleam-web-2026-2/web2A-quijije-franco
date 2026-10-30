# API de reservas deportivas · Hito 2

Backend Go para una plataforma de subasta inversa de canchas deportivas. El flujo principal es `Reserva`: inicia pendiente, el administrador puede confirmarla y las transiciones inválidas responden `422`.

## Ejecución

```powershell
$bytes = New-Object byte[] 32
[Security.Cryptography.RandomNumberGenerator]::Create().GetBytes($bytes)
$env:JWT_SECRET = [Convert]::ToBase64String($bytes)
go run .
```

`JWT_SECRET` es una variable de entorno que Go utiliza para firmar y validar tokens JWT; no es la contraseña de PostgreSQL.

Si se define `DATABASE_URL`, el arranque conecta a PostgreSQL y aplica de forma idempotente [001_init.sql](internal/db/migrations/001_init.sql), sin eliminar datos.

## Verificaciones

```powershell
go test ./...
go vet ./...
npx @redocly/cli lint docs/openapi.yaml --config redocly.yaml --format stylish
```

La colección está en [docs/hito2_bruno/](docs/hito2_bruno/), y las cuentas ficticias en [docs/hito2_hoja.md](docs/hito2_hoja.md).

Las rutas heredadas `/tickets` y `/explotar` no se registran en la API final: pertenecían a la plantilla de mesa de ayuda y no al negocio de reservas deportivas.
