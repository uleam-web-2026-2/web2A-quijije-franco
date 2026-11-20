# Semana 12 · Guía de entrega y evidencias (Quijije y Franco)

## Qué se adaptó

El trabajo se basa en el **repositorio original de Quijije y Franco**, de **subasta inversa de canchas deportivas**. No se copió el negocio de encuestas de Cantos. La entidad que se prueba es **Reserva**, asociada con Oferta y Pago, y su ruta de transición es **`PATCH /reservas/{id}/estado`** (no `PATCH /reservas/{id}`). Puerto original: **8080**. Los usuarios ficticios siguen siendo `admin@reservas.local` y `cliente@reservas.local`.

La guía del docente exige evidencia genuina: no son válidas salidas copiadas de otro proyecto, capturas retocadas ni respuestas HTTP inventadas. Por eso esta entrega incluye código, auditoría estática y un generador que registra **las respuestas reales al ejecutar el servidor**. No incluye archivos HTTP falsos.

## Antes de empezar: resolver imports en Windows

Si VS Code marca 14 errores en `main.go`, `reservas.go`, `tickets.go` o pruebas, primero ejecuta `.\ARREGLAR_IMPORTS_WINDOWS.cmd` desde la raíz del repositorio y sigue [SOLUCION_IMPORTS_VSCODE.md](SOLUCION_IMPORTS_VSCODE.md). No intentes generar evidencias hasta que `go build ./...` y `go test ./...` funcionen.

## 1. Una sola orden para las evidencias HTTP reales

En la raíz de este proyecto con **Go y Python 3** disponibles, ejecuta:

```bash
python scripts/semana12/generar_evidencias.py
```

Esto ejecuta `go build` y `go test`, inicia en un puerto libre la API **con el repositorio en memoria**, prueba siete inicios de sesión erróneos, reinicia para limpiar el contador, prepara una reserva con pago verificado y envía las seis peticiones. Después intenta inyección en el ID, ejecuta ambas peticiones previas CORS y hace un PATCH desde un origen no permitido con token válido contra **otra reserva**. Todos los resultados se guardan como archivos UTF-8 en `docs/taller12/`, con códigos y cabeceras devueltos por la API real. No hay que pegar salidas.

Archivos que se completan al ejecutar: `build_y_test.txt`, `rafaga.txt`, `logs_seis.txt`, `peticiones_seis_http.txt`, `peticiones.md`, `inyeccion.txt`, `cors_permitido.txt`, `cors_no_permitido.txt` y `cors_curl_ejecuta.txt`. Si hay un error de creación, registra la respuesta real y explica lo que falta en `problemas_ejecucion.md`. **No subas salidas inventadas** para rellenar los archivos faltantes.

El ataque del punto 1 va dirigido al ID, no al filtro `?estado=`, porque este proyecto no lo tiene. El GET normal de la reserva debe dar 200 y el intento `1 OR 1=1`, 400. El script no modifica reglas comerciales.

## 2. Reto PostgreSQL: crear usuario restringido sin subir contraseñas

El proyecto **no incluye Docker Compose**; tiene PostgreSQL opcional mediante `DATABASE_URL`. Para el reto, necesitas una instancia PostgreSQL accesible. En una terminal con `psql` y con una cuenta administradora:

```bash
psql -U postgres -d reservas_canchas -v ON_ERROR_STOP=1 -f internal/db/migrations/001_init.sql
psql -U postgres -d reservas_canchas -v ON_ERROR_STOP=1 -f scripts/semana12/crear_rol_api.sql
```

La segunda orden **pide interactivamente** una contraseña nueva para `reservas_api`; no se escribe en el repositorio. Si todavía no existe `reservas_canchas`, créala primero con `createdb -U postgres reservas_canchas`. Las diez tablas las crea el administrador y no el usuario de la API.

Configura temporalmente `DATABASE_URL` para que apunte al usuario **`reservas_api`**; en PowerShell, por ejemplo:

```powershell
$env:DATABASE_URL = "postgres://reservas_api:TU_CLAVE@localhost:5432/reservas_canchas?sslmode=disable"
python scripts/semana12/generar_evidencias.py --postgres
```

**Opcional, para preparar las dos Ofertas en una base local de prácticas:**

```bash
psql -U postgres -d reservas_canchas -v ON_ERROR_STOP=1 -f scripts/semana12/semilla_ofertas_demo.sql
```

El mensaje `OFERTAS S12 DISPONIBLES` muestra los IDs nuevos. Úsalos con `--oferta` y `--oferta-cors` en la prueba `--postgres`. La semilla incluye únicamente datos ficticios, no contraseñas reales, y debe ejecutarse solo en una base de desarrollo donde tengas permiso para crear registros. No borra datos existentes.

**Importante:** PostgreSQL aplica claves foráneas; antes de crear una reserva deben existir las Ofertas elegidas y sus datos relacionados (Cliente, Solicitud, Proveedor, Complejo, Cancha). Por defecto el script intenta `oferta_id` 1 y 2. Si en tu instancia son otras, ejecútalo así: `python scripts/semana12/generar_evidencias.py --postgres --oferta 15 --oferta-cors 16`, reemplazando los IDs por ofertas existentes que **todavía no tienen reservas**. Si no preparas esos datos, la prueba de seis peticiones con el usuario restringido queda honestamente pendiente.

En **Git Bash** (no en Windows PowerShell 5 con `>`), guarda las evidencias reales del reto:

```bash
psql -U reservas_api -d reservas_canchas -c 'BEGIN' -c 'DROP TABLE public.reservas' -c 'ROLLBACK' > docs/taller12/reto_drop.txt 2>&1
psql -U postgres -d reservas_canchas -c '\du' > docs/taller12/reto_usuarios.txt 2>&1
```

`reto_drop.txt` debe mostrar un error de falta de propiedad: si permite `DROP TABLE`, hay que revisar el rol. Las consultas se hacen dentro de una transacción y el `ROLLBACK` protege de borrados accidentales. Las contraseñas quedan **fuera de Git**. Si antes usabas migración automática, ahora la ejecución normal de la API **no migra**: se separó intencionalmente mediante `db.Open` y `db.Migrate`.

## 3. Publicación Git

Trabaja dentro del **repositorio Git real de Quijije y Franco**, no dentro del ZIP del otro equipo. Conserva la carpeta oculta `.git` de tu clon original, copia o reemplaza los archivos de este paquete sobre ese clon, y crea/cambia a `taller12` a partir de tu último taller:

```bash
git status
git switch -c taller12
# Si ya existe: git switch taller12
# Copia aquí los archivos preparados sobre tu clon, excepto .git
python scripts/semana12/generar_evidencias.py
git status
git add main.go main_test.go seguridad_semana12_test.go go.mod go.sum internal/db/postgres.go README.md docs/ scripts/
git commit -m "feat: adapta la seguridad de la semana 12 a reservas deportivas" -m "IA: ChatGPT para adaptar codigo Go, auditoria SQL, documentacion y automatizacion de pruebas; las respuestas HTTP y PostgreSQL se obtienen solo al ejecutarlas realmente."
git push -u origin taller12
```

El docente pide preferentemente **un commit por punto y otro para el reto**, y que se registre actividad de ambos integrantes. Si quieres cumplir estrictamente esa parte, separa los archivos/etapas en varios commits antes de hacer `git push`; la adaptación del ZIP no contiene historia Git ajena ni simula commits de otros estudiantes.

El script necesita poder descargar las dependencias Go si todavía no están instaladas. Si falla `go build` por red/versión de Go, queda grabado el error real en `build_y_test.txt` y no se presentarán pruebas HTTP inexistentes.
