# Cómo solucionar los imports subrayados en VS Code (Taller 12)

Este proyecto es **Quijije y Franco — subasta inversa y reserva de canchas**. No renombres el módulo de Go ni cambies las rutas de `internal/` por las del taller de Cantos (OPINAEC).

## 1. Abre la carpeta correcta

En VS Code, **Archivo → Abrir carpeta** y selecciona la que contiene `main.go`, `go.mod`, `go.sum` e `internal/`, no la carpeta superior del ZIP. Abre **Terminal → Nueva terminal**. Confirma con `dir go.mod` (PowerShell) o `ls go.mod` (Git Bash).

## 2. Restaura dependencias y verifica

Desde la carpeta del proyecto, en **PowerShell**, ejecuta:

```powershell
.\ARREGLAR_IMPORTS_WINDOWS.cmd
```

O, manualmente (también sirve en Git Bash):

```bash
go version
go mod tidy
go mod download
go build ./...
go test ./...
```

Necesitas **Go 1.27.1 o posterior**, acceso a Internet la primera vez y conexión a `proxy.golang.org` o al proxy de Go que use tu universidad. El script detiene las pruebas al encontrar un error; no escribe evidencias HTTP ficticias.

## 3. Actualiza VS Code

1. `Ctrl+Shift+P` → **Go: Restart Language Server**.
2. Si no aparece ese comando, verifica que esté instalada y habilitada la extensión **Go** (Go Team at Google).
3. `Ctrl+Shift+P` → **Developer: Reload Window** si los subrayados siguen incluso después de que `go build ./...` funcione.

## 4. Si hay errores

- **`go: go.mod requires go >= 1.27.1`**: instala una versión adecuada de Go y reinicia la terminal de VS Code (no cambies arbitrariamente `go.mod`).
- **`missing go.sum entry`**, **`no required module provides package`**: ejecuta `go mod tidy` y `go mod download` desde el directorio que contiene `go.mod`.
- **`dial tcp ... lookup proxy.golang.org`**: falla DNS o Internet. Conecta el equipo y ejecuta de nuevo el script.
- **`module ... found, but does not contain package`**: comprueba la ruta del módulo en `go.mod`; aquí debe ser `github.com/uleam-web-2026-2/servA-quijije-franco`.
- **Errores distintos después de que descargue todo**: copia el texto exacto de `go build ./...`; será un problema real de código y no solo de caché de VS Code.

## 5. Completa las evidencias verdaderas del Taller 12

Después de que `go build ./...` y `go test ./...` pasen, ejecuta `python scripts/semana12/generar_evidencias.py` (sin `--postgres` para pruebas en memoria) y revisa los resultados dentro de `docs/taller12`. Los archivos `reto_drop.txt` y `reto_usuarios.txt` exigen una base PostgreSQL real y no deben inventarse.

**No subas a Git** contraseñas, `.env`, ni salidas inventadas. Si `go mod tidy` cambia `go.mod`/`go.sum`, esos archivos sí forman parte de la corrección y deben ir al commit.
