-- EJECUTAR SOLO COMO ADMINISTRADOR EN LA BASE reservas_canchas, con psql -f.
-- No almacena una contraseña. psql pedirá la contraseña de reservas_api.
DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname='reservas_api') THEN
    CREATE ROLE reservas_api LOGIN NOSUPERUSER NOCREATEDB NOCREATEROLE NOREPLICATION;
  END IF;
END $$;
ALTER ROLE reservas_api NOSUPERUSER NOCREATEDB NOCREATEROLE NOREPLICATION;
\password reservas_api
GRANT CONNECT ON DATABASE reservas_canchas TO reservas_api;
GRANT USAGE ON SCHEMA public TO reservas_api;
REVOKE ALL ON ALL TABLES IN SCHEMA public FROM reservas_api;
REVOKE ALL ON ALL SEQUENCES IN SCHEMA public FROM reservas_api;
-- Consultas reales desde internal/reservas/repositorio_sql.go:
GRANT SELECT, INSERT, UPDATE ON TABLE public.reservas TO reservas_api;
GRANT SELECT, INSERT ON TABLE public.pagos TO reservas_api;
GRANT SELECT ON TABLE public.ofertas, public.solicitudes, public.clientes TO reservas_api;
GRANT USAGE ON SEQUENCE public.reservas_reserva_id_seq, public.pagos_pago_id_seq TO reservas_api;
-- El rol no recibe CREATE ni ALTER ni DROP sobre las tablas del proyecto.
