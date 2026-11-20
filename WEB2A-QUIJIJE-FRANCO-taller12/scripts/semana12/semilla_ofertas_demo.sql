-- OPCIONAL: SOLO PARA UNA BASE LOCAL DE PRUEBAS, como reservas_canchas.
-- Crea dos Ofertas NUEVAS coherentes con el negocio de subasta inversa de canchas.
-- La columna password de los usuarios ficticios recibe un MARCADOR NO UTILIZABLE;
-- el login demo de main.go se autentica de manera independiente en memoria.
-- Ejecutar como administrador: psql -U postgres -d reservas_canchas -f scripts/semana12/semilla_ofertas_demo.sql
DO $$
DECLARE
  id_usuario_cliente BIGINT;
  id_usuario_proveedor BIGINT;
  id_cliente BIGINT;
  id_proveedor BIGINT;
  id_complejo BIGINT;
  id_cancha BIGINT;
  id_solicitud BIGINT;
  id_oferta1 BIGINT;
  id_oferta2 BIGINT;
BEGIN
  INSERT INTO usuarios(nombre, correo, password, rol)
    VALUES('Cliente demostración S12', 'cliente-semana12@demo.invalid', 'NO_ES_PASSWORD_DE_LOGIN', 'cliente')
    ON CONFLICT(correo) DO UPDATE SET nombre=EXCLUDED.nombre
    RETURNING usuario_id INTO id_usuario_cliente;
  INSERT INTO usuarios(nombre, correo, password, rol)
    VALUES('Proveedor demostración S12', 'proveedor-semana12@demo.invalid', 'NO_ES_PASSWORD_DE_LOGIN', 'proveedor')
    ON CONFLICT(correo) DO UPDATE SET nombre=EXCLUDED.nombre
    RETURNING usuario_id INTO id_usuario_proveedor;
  INSERT INTO clientes(usuario_id, cedula, telefono)
    VALUES(id_usuario_cliente, 'S12-CLIENTE-DEMO', '0990000000')
    ON CONFLICT(usuario_id) DO UPDATE SET telefono=EXCLUDED.telefono
    RETURNING cliente_id INTO id_cliente;
  INSERT INTO proveedores(usuario_id, ruc, telefono, banco, numero_cuenta, titular_cuenta)
    VALUES(id_usuario_proveedor, 'S12-PROVEEDOR-DEMO', '0990000000', 'DEMO', 'SIN-CUENTA-REAL', 'Proveedor Demo')
    ON CONFLICT(usuario_id) DO UPDATE SET telefono=EXCLUDED.telefono
    RETURNING proveedor_id INTO id_proveedor;
  INSERT INTO complejos(proveedor_id, nombre, direccion, ciudad)
    VALUES(id_proveedor, 'Complejo de prueba Semana 12', 'Local de pruebas', 'Manta')
    RETURNING complejo_id INTO id_complejo;
  INSERT INTO canchas(complejo_id, tipo, capacidad, precio_referencia)
    VALUES(id_complejo, 'fútbol 5', 10, 35)
    RETURNING cancha_id INTO id_cancha;
  INSERT INTO solicitudes(cliente_id, fecha_evento, duracion_minutos, cantidad_personas, presupuesto_maximo, estado)
    VALUES(id_cliente, '2030-10-15T19:00:00Z', 90, 10, 50, 'abierta')
    RETURNING solicitud_id INTO id_solicitud;
  INSERT INTO ofertas(solicitud_id, complejo_id, cancha_id, monto_total, detalle, vence_en, estado)
    VALUES(id_solicitud, id_complejo, id_cancha, 35, 'S12 prueba de seguridad - 1', '2030-10-14T23:00:00Z', 'aceptada')
    RETURNING oferta_id INTO id_oferta1;
  INSERT INTO ofertas(solicitud_id, complejo_id, cancha_id, monto_total, detalle, vence_en, estado)
    VALUES(id_solicitud, id_complejo, id_cancha, 35, 'S12 prueba de seguridad - 2', '2030-10-14T23:00:00Z', 'aceptada')
    RETURNING oferta_id INTO id_oferta2;
  RAISE NOTICE 'OFERTAS S12 DISPONIBLES: --oferta % --oferta-cors %', id_oferta1, id_oferta2;
END $$;
