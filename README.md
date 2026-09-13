# Servidor 

Servidor HTTP desarrollado en Go utilizando el enrutador Chi para la gestión centralizada de tickets de soporte técnico universitario.

## Tabla de Rutas

| Método | Ruta | Éxito | Errores posibles |
| ------ | ---- | ----- | ---------------- |
| `GET` | `/tickets` | `200 OK` (lista de tickets) | - |
| `POST` | `/tickets` | `201 Created` (ticket creado) | `400` (`cuerpo_malformado`, `datos_invalidos`) |
| `GET` | `/tickets/{id}` | `200 OK` (ticket encontrado) | `400` (`id_invalido`), `404` (`no_encontrado`) |
| `GET` | `/explotar` | - | `500` (`error_interno`) |
| `*` | Ruta inexistente | - | `404` (`ruta_inexistente`) |
| `*` | Método no soportado | - | `405` (`metodo_no_permitido`) |

---

Revisión cruzada: Revisado por [Nombre de tu compañero]: sin hallazgos