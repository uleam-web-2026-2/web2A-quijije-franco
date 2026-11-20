#!/usr/bin/env python3
"""Pruebas HTTP REALES de Taller 12, sin respuestas inventadas ni librerías Python externas.
Ejecutar desde la raíz del repositorio: python scripts/semana12/generar_evidencias.py
Por defecto prueba el repositorio en memoria. --postgres usa DATABASE_URL del entorno.
"""
import argparse
import json
import os
import shutil
import socket
import subprocess
import sys
import tempfile
import threading
import time
import urllib.error
import urllib.request
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
EVIDENCIAS = ROOT / 'docs' / 'taller12'
EVIDENCIAS.mkdir(parents=True, exist_ok=True)

class API:
    def __init__(self, executable, port, usar_postgres):
        self.executable, self.port, self.usar_postgres = executable, port, usar_postgres
        self.proc = None
        self.lines = []

    def start(self):
        env = os.environ.copy()
        env['PORT'] = str(self.port)
        env.pop('RUN_MIGRATIONS', None)
        if not self.usar_postgres:
            env.pop('DATABASE_URL', None)
        self.lines = []
        self.proc = subprocess.Popen([str(self.executable)], cwd=ROOT, env=env,
                                     stdout=subprocess.PIPE, stderr=subprocess.STDOUT,
                                     text=True, encoding='utf-8', errors='replace')
        def read_logs():
            for line in self.proc.stdout:
                self.lines.append(line)
        threading.Thread(target=read_logs, daemon=True).start()
        for _ in range(100):
            if self.proc.poll() is not None:
                raise RuntimeError('La API terminó al arrancar:\n' + ''.join(self.lines))
            try:
                with socket.create_connection(('127.0.0.1', self.port), timeout=.2):
                    return
            except OSError:
                time.sleep(.1)
        raise RuntimeError('La API no abrió el puerto: ' + ''.join(self.lines))

    def stop(self):
        if self.proc is not None:
            self.proc.terminate()
            try:
                self.proc.wait(timeout=5)
            except subprocess.TimeoutExpired:
                self.proc.kill()
                self.proc.wait()
            time.sleep(.15)
            self.proc = None


def free_port():
    s = socket.socket()
    s.bind(('127.0.0.1', 0))
    v = s.getsockname()[1]
    s.close()
    return v


def request(port, method, path, body=None, headers=None):
    data = (json.dumps(body, ensure_ascii=False).encode('utf-8') if body is not None else None)
    h = dict(headers or {})
    if body is not None: h['Content-Type'] = 'application/json'
    req = urllib.request.Request(f'http://127.0.0.1:{port}{path}', data=data,
                                 headers=h, method=method)
    try:
        response = urllib.request.urlopen(req, timeout=10)
    except urllib.error.HTTPError as e:
        response = e
    with response:
        content = response.read().decode('utf-8', 'replace')
        result = {
            'status': response.status,
            'headers': list(response.headers.items()),
            'body': content,
        }
    result['raw'] = (f'HTTP/1.1 {result["status"]} {response.reason}\n'
                     + ''.join(f'{k}: {v}\n' for k,v in result['headers'])
                     + '\n' + content + '\n')
    try:
        result['json'] = json.loads(content)
    except ValueError:
        result['json'] = {}
    return result


def write(name, content):
    (EVIDENCIAS / name).write_text(content, encoding='utf-8')


def save_http(name, results):
    write(name, '\n'.join(r['raw'] for r in results))


def token(result):
    return result.get('json', {}).get('datos', {}).get('token', '')


def crear_reserva(port, bearer, oferta_id):
    return request(port, 'POST', '/reservas', {
        'oferta_id': oferta_id, 'inicio': '2030-10-15T19:00:00Z',
        'fin': '2030-10-15T20:30:00Z', 'monto_acordado': 35,
        'pago_estado': 'verificado', 'pago_url': 'demo://comprobante-semana12'
    }, {'Authorization': 'Bearer ' + bearer})


def main():
    args = argparse.ArgumentParser(description=__doc__)
    args.add_argument('--postgres', action='store_true', help='usar DATABASE_URL del entorno (reto)')
    args.add_argument('--oferta', type=int, default=1, help='ID de una oferta disponible para crear Reserva (defecto: 1)')
    args.add_argument('--oferta-cors', type=int, default=2, help='ID de otra oferta disponible para la prueba CORS (defecto: 2)')
    opts = args.parse_args()
    if opts.postgres and not os.environ.get('DATABASE_URL'):
        raise RuntimeError('--postgres requiere DATABASE_URL del usuario reservas_api')
    if not shutil.which('go'):
        raise RuntimeError('No se encontró Go en PATH')

    with tempfile.TemporaryDirectory(prefix='quijije-franco-s12-') as temp:
        binfile = Path(temp) / ('api.exe' if os.name == 'nt' else 'api')
        registros = []
        for cmd in [['go', 'build', '-o', str(binfile), '.'], ['go', 'test', './...']]:
            run = subprocess.run(cmd, cwd=ROOT, text=True, stdout=subprocess.PIPE,
                                 stderr=subprocess.STDOUT, encoding='utf-8', errors='replace')
            registros.append('$ ' + ' '.join(cmd[:2]) + '\n' + run.stdout + f'Exit code: {run.returncode}\n')
            if run.returncode != 0:
                write('build_y_test.txt', '\n'.join(registros))
                raise RuntimeError('Compilación/pruebas fallaron; ver docs/taller12/build_y_test.txt')
        write('build_y_test.txt', '\n'.join(registros))
        port = free_port()
        server = API(binfile, port, opts.postgres)
        estado = {}
        try:
            server.start()
            # 7 intentos erróneos REALES: no permitir logins anteriores en esta instancia.
            rafaga = [request(port, 'POST', '/login', {'correo':'admin@reservas.local', 'password':'incorrecta-12'})
                      for _ in range(7)]
            save_http('rafaga.txt', rafaga)
            estado['ráfaga'] = [r['status'] for r in rafaga]
            server.stop()
            server.start()  # contador de IP limpio

            # Las 6 peticiones corresponden a docs/hito2_hoja.md del proyecto.
            r1 = request(port, 'GET', '/reservas')
            r2 = request(port, 'POST', '/login', {'correo':'admin@reservas.local', 'password':'admin-2026'})
            admin = token(r2)
            r3 = request(port, 'GET', '/reservas', headers={'Authorization':'Bearer ' + admin})
            cliente_login = request(port, 'POST', '/login', {'correo':'cliente@reservas.local','password':'cliente-2026'})
            cliente = token(cliente_login)
            creada = crear_reserva(port, admin, opts.oferta)
            id_reserva = creada.get('json',{}).get('datos',{}).get('reserva_id')
            if creada['status'] != 201 or not id_reserva:
                write('problemas_ejecucion.md', '# Reserva de prueba no creada\n\nRespuesta REAL:\n\n```http\n'+creada['raw']+'\n```\n\nEn PostgreSQL, primero debe existir la oferta elegida y sus relaciones.\n')
                raise RuntimeError('No se pudo crear reserva; ver docs/taller12/problemas_ejecucion.md')
            url = f'/reservas/{id_reserva}/estado'
            r4 = request(port, 'PATCH', url, {'estado':'confirmada'}, {'Authorization':'Bearer ' + cliente})
            r5 = request(port, 'PATCH', url, {'estado':'confirmada'}, {'Authorization':'Bearer ' + admin})
            r6 = request(port, 'PATCH', url, {'estado':'confirmada'}, {'Authorization':'Bearer ' + admin})
            seis = [r1,r2,r3,r4,r5,r6]
            estado['seis'] = [r['status'] for r in seis]
            save_http('peticiones_seis_http.txt', seis)
            time.sleep(.3)
            write('logs_seis.txt', ''.join(server.lines))
            # Aquí no hay filtro de estado por URL: se verifica el ID numérico.
            normal = request(port, 'GET', f'/reservas/{id_reserva}', headers={'Authorization':'Bearer ' + admin})
            ataque = request(port, 'GET', '/reservas/1%20OR%201%3D1', headers={'Authorization':'Bearer ' + admin})
            save_http('inyeccion.txt', [normal, ataque])
            estado['inyección'] = [normal['status'], ataque['status']]
            cors_headers = {'Access-Control-Request-Method':'PATCH', 'Access-Control-Request-Headers':'authorization,content-type'}
            permitido = request(port,'OPTIONS',url,headers={**cors_headers,'Origin':'http://localhost:5173'})
            no_permitido = request(port,'OPTIONS',url,headers={**cors_headers,'Origin':'http://malo.example'})
            save_http('cors_permitido.txt',[permitido])
            save_http('cors_no_permitido.txt',[no_permitido])
            estado['cors'] = [permitido['status'],no_permitido['status']]
            # Un registro NUEVO: el confirmado ya no admite confirmar de nuevo.
            segunda = crear_reserva(port, admin, opts.oferta_cors)
            if segunda['status'] == 201:
                nuevo_id = segunda['json']['datos']['reserva_id']
                curl = request(port, 'PATCH', f'/reservas/{nuevo_id}/estado', {'estado':'confirmada'},
                               {'Origin':'http://malo.example','Authorization':'Bearer '+admin})
                save_http('cors_curl_ejecuta.txt',[curl])
                estado['curl'] = curl['status']
            else:
                write('cors_curl_ejecuta_pendiente.md', '# Pendiente: crear otra oferta/reserva\n\n'+segunda['raw'])
            # Tabla de resultados que corresponde a peticiones enviadas realmente.
            # No borrar la evidencia de la otra modalidad al ejecutar de nuevo.
            anteriores = {}
            ruta_tabla = EVIDENCIAS / 'peticiones.md'
            if ruta_tabla.exists():
                for linea in ruta_tabla.read_text(encoding='utf-8').splitlines():
                    cols = [x.strip() for x in linea.strip().strip('|').split('|')]
                    if len(cols) == 4 and cols[0][:1].isdigit():
                        anteriores[int(cols[0].split('.')[0])] = cols
            rows = ['# Taller 12 · Reserva deportiva (Quijije y Franco)', '',
                    'Seis peticiones de la API; solo las columnas ejecutadas son evidencia.', '',
                    '| Petición | Antes (taller 11) | Con las defensas | Con usuario API (reto) |',
                    '|---|---|---|---|']
            for n,(desc,v) in enumerate(zip(['Sin JWT','Login admin','Listado con token','Cliente no confirma','Admin confirma','Transición prohibida'],estado['seis']),1):
                antes = anteriores.get(n, ['', 'Sin salida anterior verificable', 'Pendiente', 'Pendiente'])
                col = (antes[2] if opts.postgres else str(v))
                rol = (str(v) if opts.postgres else antes[3])
                rows.append(f'| {n}. {desc} | {antes[1]} | {col} | {rol} |')
            rows.extend(['', 'El repositorio de origen no traía salidas originales de la semana 11 para estas seis peticiones.',
                         'La petición previa 0 crea una reserva con un pago verificado; se ejecutó para esta serie.',
                         '', f'**Modalidad de esta captura:** {"PostgreSQL, DATABASE_URL del entorno" if opts.postgres else "Repositorio en memoria"}.',
                         '', 'Los textos completos de las seis respuestas están en `peticiones_seis_http.txt`.'])
            write('peticiones.md','\n'.join(rows)+'\n')
            print('Peticiones HTTP reales completadas. Respuestas:', estado)
        finally:
            server.stop()
    if estado.get('ráfaga') != [401]*5 + [429]*2 or estado.get('seis') != [401,200,200,403,200,422]:
        print('ATENCIÓN: códigos distintos de la guía. Revisar archivos reales guardados.',file=sys.stderr)
        return 1
    return 0

if __name__ == '__main__':
    try:
        sys.exit(main())
    except Exception as exc:
        print('EVIDENCIAS INCOMPLETAS: '+str(exc), file=sys.stderr)
        sys.exit(1)
