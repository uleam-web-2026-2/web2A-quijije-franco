-- Ejecutar con psql -U reservas_api -d reservas_canchas -f scripts/semana12/reto_consultas.sql
-- BEGIN/ROLLBACK protegen la tabla incluso si el rol tuviera derechos de más.
BEGIN;
DROP TABLE public.reservas;
ROLLBACK;
