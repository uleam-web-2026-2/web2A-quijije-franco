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

Si se define `DATABASE_URL`, el arranque conecta a PostgreSQL **sin ejecutar migraciones**. Las tablas se crean por separado con cuenta administradora, como indica [la guía del reto](docs/taller12/GUIA_DE_PRUEBAS.md). Para ejecutar las migraciones administrativas se usa `RUN_MIGRATIONS=true` con las credenciales correspondientes; no se debe usar el rol limitado de la API para crear tablas.

## Verificaciones

```powershell
go test ./...
go vet ./...
npx @redocly/cli lint docs/openapi.yaml --config redocly.yaml --format stylish
```

La colección está en [docs/hito2_bruno/](docs/hito2_bruno/), y las cuentas ficticias en [docs/hito2_hoja.md](docs/hito2_hoja.md).

Las rutas heredadas `/tickets` y `/explotar` no se registran en la API final: pertenecían a la plantilla de mesa de ayuda y no al negocio de reservas deportivas.

## Semana 12 · Seguridad de la API (Quijije y Franco)

Esta rama aplica al **proyecto de subasta inversa de canchas deportivas** el Taller 12. Las consultas de PostgreSQL de `Reserva` y `Pago` permanecen parametrizadas (`$1`, `$2`, …); el ID de ruta se valida numéricamente. Se agregó CORS para `http://localhost:5173`, un límite de **5 POST /login por minuto y por IP**, y se separó la migración de la conexión habitual de la API para habilitar una cuenta PostgreSQL limitada (`reservas_api`).

**Primero leer [la guía de la semana 12](docs/taller12/GUIA_DE_PRUEBAS.md)** para completar con resultados reales `docs/taller12/`, seguir las pruebas de Bruno y configurar el reto sin guardar contraseñas. El código conserva `PATCH /reservas/{id}/estado`, las reglas de pago verificado, comisión de 10 % y las seis respuestas del Hito 2.

## Si VS Code muestra imports subrayados

Abre esta carpeta como proyecto y ejecuta desde ella **`.\ARREGLAR_IMPORTS_WINDOWS.cmd`** en la terminal de Windows. Después reinicia el servidor de lenguaje Go desde `Ctrl+Shift+P`. Lee [SOLUCION_IMPORTS_VSCODE.md](docs/taller12/SOLUCION_IMPORTS_VSCODE.md) si alguna descarga falla.
