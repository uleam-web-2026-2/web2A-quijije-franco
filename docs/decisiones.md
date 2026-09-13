## D3 Manejo de campos desconocidos en el JSON
**Decisión:** Rechazar peticiones que contengan campos no declarados usando `dec.DisallowUnknownFields()`.
- **Por qué:** Evita que clientes o programas externos intenten modificar campos del dominio que no deben alterar en la creación (por ejemplo, intentar enviar `"estado": "cerrado"` al crear el ticket).
- **Alternativa descartada:** Ignorar los campos desconocidos en silencio. Se descartó porque puede ocultar errores de tipeo del cliente sin notificarle.
- **Cómo lo comprobamos:** Enviando una petición POST con `{"titulo": "Ticket Prueba", "prioridad": "alta", "estado": "cerrado"}` y confirmando que responde con código `400 Bad Request` y `cuerpo_malformado`.