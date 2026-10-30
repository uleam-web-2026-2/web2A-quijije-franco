# Hito 2 — Su API, demostrada en vivo

> Transcripción de lectura del PDF original `S8_Hito2_Servidor_Enunciado_y_rubrica.pdf`. Se conserva el texto y el orden de las 10 páginas, incluidas tablas en disposición monoespaciada. En caso de discrepancia visual, prevalece el PDF original.

## Página 1 de 10

```text
 Proyecto autónomo · Hito 2 — Su API,
 demostrada en vivo
 ULEAM · Facultad de Ciencias de la Vida y Tecnologías · Período 2026-2
 Hilo Servidor: Aplicaciones Web II (TI, paralelos A y B) · Aplicación para el Servidor Web
 (SW, paralelo A)
 Docente: Ing. John Antonio Cevallos Macías, Mg.


 Lo esencial
   Dato                      Detalle

   Peso                      7 % de la nota del período. Es el segundo de los tres hitos del proyecto
                             autónomo (semanas 5, 9 y 15).

   Qué se hace               Una demostración en vivo de cinco minutos: la pareja envía seis peticiones a
                             su propia API, corre sus pruebas y muestra su código organizado en capas. No
                             hay diapositivas.

   Qué se                    En la carpeta docs/ del repositorio de la pareja: el contrato ( openapi.yaml ),
   entrega                   la hoja del hito ( hito2_hoja.md ), la colección de Bruno ( hito2_bruno/ ) y las
                             capturas de respaldo ( hito2_capturas/ ). Además, el código de la entidad
                             principal organizado en capas, con la estructura explicada en
                              docs/decisiones.md .

   Cierre de la              La víspera de su demostración, 23:59. Servidor Web A y Web II B cierran el
   entrega                   miércoles 28 de octubre; Web II A, el jueves 29 de octubre. La entrega es el
                             repositorio de la pareja dentro de la organización uleam-web-2026-2 : no
                             vale un repositorio personal ni uno de fuera de la organización. Si a esa hora la
                             entrega no está ahí, la pareja no demuestra ese día: lo hace en horario de
                             tutorías por el 50 %.

   Demostración              Servidor Web A: jueves 29 de octubre, 08:00. Web II B: jueves 29 de octubre,
                             13:00. Web II A: viernes 30 de octubre, 09:00. La sesión B completa se dedica
                             al hito: cinco minutos por pareja, en el orden de la lista.

   Si no se                  La pareja que no está en el aula cuando el docente la llama, aunque llegue más
   presentan                 tarde ese mismo día, demuestra en horario de tutorías y su nota es el 50 % de
                             la que obtenga. La única excepción es una justificación validada por la
                             Dirección de Carrera: el docente no puede aceptar justificaciones.




Hito 2 · Proyecto autónomo · Hilo Servidor · 2026-2                                                         1 / 10
```

## Página 2 de 10

```text
 1. De dónde parte este hito
 En el hito 1 presentaron el negocio, el modelo, la máquina de estados y el mapa de
 endpoints en papel. El hito 2 muestra eso funcionando, con lo que construyen en los
 talleres:

   Taller                    Lo que queda en su producto

   Semana 4                  Configuración fuera del código y pruebas con testing

   Semana 7                  El contrato de la API en docs/openapi.yaml , validado · el cambio de estado
   (taller                   con PATCH , que rechaza la transición prohibida · inicio de sesión, token, rutas
   integrador)               protegidas y permiso por rol · una prueba de integración

 | Semana 8 (clase y guía de pasos) | El código de su entidad principal organizado en capas
 |
 El hito 2 no pide nada que esos talleres y la guía no dejen hecho. Pide juntarlo y
 demostrar que funciona. Es lo mismo que hicieron en la misión de la semana 7 contra la
 mesa de ayuda, pero ahora contra su propia API y esperando que responda bien.

 Alcance: qué tiene que funcionar y qué va solo escrito

   Qué                             Cuánto

   Structs y tablas                Todos. Los structs de todas las entidades del modelo que presentaron en
                                   el hito 1, con sus relaciones, y sus tablas creadas en la base. Si el modelo
                                   cambió desde el hito 1, actualicen la sección 4 de la ficha.

   Entidad principal,              Los endpoints que ya tiene desde la semana 3 (como mínimo listar, ver
   la que tiene estados            una y crear; modificar y borrar si los hicieron), más cambiar de estado.
                                   Es la que se demuestra y se califica.

   Segunda entidad,                Opcional, y es el tope: solo crear y listar dentro de la principal (por
   la relacionada                  ejemplo, los comentarios de un ticket). No suma ni resta nota.

   Usuarios                        Una tabla, solo para iniciar sesión, con dos usuarios de ejemplo, uno por
                                   rol. Sin registro ni CRUD de usuarios. Si su modelo ya tiene usuarios
                                   propios, esta tabla es aparte (en el taller se llamó cuentas ).

   Las demás                       Sin endpoints programados por ahora: se construyen después del hito. En
   entidades                       el contrato sí van, marcadas como plan (apartado 2).

 No construyan endpoints de más de dos entidades antes del hito: la guía de la semana 8
 reorganiza en capas cada entidad con endpoints, y con más de dos ese trabajo se va en
 mover archivos.

 La guía de la semana 8: capas
 En la clase de la semana 8 (martes 20 y miércoles 21 de octubre) se explica cómo
 organizar el código en capas. Ese mismo día reciben una guía de pasos para aplicarlo a su
 proyecto, por su cuenta.
 Las capas se califican en el hito 2: son la demostración 5. Lo que se pide:




Hito 2 · Proyecto autónomo · Hilo Servidor · 2026-2                                                            2 / 10
```

## Página 3 de 10

```text
      El código de su entidad principal está repartido en partes, cada una con una sola
      responsabilidad. Como mínimo: la función que atiende la petición HTTP no consulta la
      base de datos directamente, y la regla de su negocio (qué transición vale y cuál no) está
      en un lugar que no sabe nada de HTTP.
      La estructura la eligen ustedes. Puede ser una de las que se revisan en clase u otra
      distinta, con otros nombres de carpetas, siempre que cumpla lo anterior y la puedan
      explicar.
      En docs/decisiones.md , una sección «Semana 8 · estructura»: qué estructura eligieron,
      qué va en cada carpeta y por qué.
 Al mover archivos es fácil romper algo. Trabajen en commits pequeños y, al terminar,
 vuelvan a enviar las seis peticiones y a correr go test ./... : todo tiene que responder
 igual que antes. Revisen también la pregunta 1 («¿en qué línea se decide que esta ruta
 pide token?»): esa línea habrá cambiado de archivo.


 2. Qué entregan
   En docs/                      Qué es                                               Para qué sirve

    openapi.yaml y               El contrato de su API, como quedó en el taller       Es contra lo que se
    api.html                     integrador: declara security y las respuestas        compara cada
                                 401, 403 y 422 en las rutas que las dan. Lleva       respuesta de la
                                 todos los endpoints del mapa del hito 1 (sección     demostración.
                                 7 de la ficha), de todas las entidades: los que ya
                                 existen en su servidor y los que todavía no. Los
                                 que no existen dicen «(plan)» al final de su
                                 summary .

    hito2_hoja.md                La hoja del hito, de una página (apartado 4).        Dice qué van a enviar y
                                                                                      qué debe responder. El
                                                                                      docente la tiene abierta
                                                                                      mientras ustedes
                                                                                      demuestran.

    hito2_bruno/                 La colección de Bruno con las seis peticiones,       En la demostración se
                                 numeradas del 1 al 6, y una petición 0 de            envían desde aquí. No
                                 preparación, que crea un registro nuevo.             se escribe nada en el
                                                                                      momento.

    hito2_capturas/              Una captura de cada respuesta y una de go test       Es el respaldo si algo
                                 ./... .                                              no corre ese día.

    decisiones.md                Una sección «Semana 8 · estructura»: qué             Es lo que se compara
                                 estructura eligieron, qué va en cada carpeta y por   con su código en la
                                 qué.                                                 demostración 5.

 El contrato del taller integrador no alcanza: hay que completarlo. Del taller sale con
 las operaciones de su entidad principal y el inicio de sesión. Para el hito le agregan los
 endpoints que faltan del mapa. Para cada uno: copien la operación de su entidad principal
 que más se le parezca (listar, ver una, crear), cambien la ruta, el operationId y los
 esquemas por los de esa entidad, y terminen el summary con «(plan)». Cada entidad nueva




Hito 2 · Proyecto autónomo · Hilo Servidor · 2026-2                                                            3 / 10
```

## Página 4 de 10

```text
 necesita su esquema en components.schemas , con los campos de su struct. No hay que
 programar esos endpoints: el contrato dice lo que la API va a ofrecer, y el validador tiene
 que seguir terminando sin errores ni avisos.
 Además, en la raíz del repositorio tiene que estar redocly.yaml , con la regla security-
 defined ya encendida (paso 2.9 del taller integrador): con esas reglas corre el validador en
 la demostración 1.


 3. Las cinco demostraciones
   #      Demostración             Qué se tiene que ver

   1      El contrato              Corren el validador ( npx @redocly/cli lint docs/openapi.yaml --
                                   format=stylish ) y termina sin errores ni avisos. Señalan en el contrato
                                   dónde está security y dónde están el 401, el 403 y el 422.

   2      Quién es                 Petición 1: una ruta protegida, sin token → 401. Petición 2: inicio de
                                   sesión → entrega el token. Petición 3: la misma ruta, con el token → 200.

   3      Qué puede                Petición 4: el cambio de estado, que en su servidor es de un solo rol, con
          hacer                    el token del otro rol → 403. Petición 5: la misma petición, con el rol
                                   correcto → 200.

   4      La regla del             Petición 6: la transición prohibida de su máquina de estados, con el rol
          negocio                  correcto → 422. Después, go test ./... termina en verde.

   5      La estructura            Muestran las carpetas de su entidad principal y abren dos archivos: el que
                                   atiende la petición HTTP y el que tiene la regla del cambio de estado. Se
                                   ve que el primero no consulta la base directamente y que el segundo no
                                   sabe nada de HTTP. Coincide con lo que dice docs/decisiones.md .

 Los códigos que se esperan son los que dice su contrato. Si su contrato dice otro código y
 lo pueden justificar, vale: lo que no vale es que la API responda una cosa y el contrato diga
 otra.


 4. La hoja del hito
  docs/hito2_hoja.md tiene cuatro partes con estos títulos. Aquí va llena con la mesa de
 ayuda, como ejemplo.
 1 · Usuarios de ejemplo

   Usuario                      Rol                       Contraseña de ejemplo

    ana                         solicitante               ana-2026

    luis                        agente                    luis-2026

 Son usuarios de prueba, con contraseñas inventadas para la demostración. Nunca una
 contraseña real.




Hito 2 · Proyecto autónomo · Hilo Servidor · 2026-2                                                            4 / 10
```

## Página 5 de 10

```text
 2 · Las seis peticiones

   #     Petición                                                  Con qué credencial       Debe
                                                                                            responder

   1      GET /tickets                                             Ninguna                  401

   2      POST /login con ana                                      Ninguna                  200 y el token

   3      GET /tickets                                             Token de Ana             200

   4      PATCH /tickets/2 → en_proceso                            Token de Ana             403
                                                                   (solicitante)

   5      PATCH /tickets/2 → en_proceso                            Token de Luis (agente)   200

   6      PATCH /tickets/1 → cerrado , estando                     Token de Luis (agente)   422
          abierto

 3 · Matriz ruta × rol

   Ruta                                               Sin sesión             Solicitante       Agente

    POST /login                                       Sí                     Sí                Sí

    GET /tickets                                      No (401)               Sí                Sí

    POST /tickets                                     No (401)               Sí                Sí

    GET /tickets/{id}                                 No (401)               Sí                Sí

    PUT /tickets/{id}                                 No (401)               No (403)          Sí

    PATCH /tickets/{id}                               No (401)               No (403)          Sí

    DELETE /tickets/{id}                              No (401)               No (403)          Sí

 Va una fila por cada ruta que existe en su servidor: es la misma matriz D1 del taller
 integrador.
 4 · Recuadro de cumplimiento
 ☐ Los structs de todas las entidades están en el código y sus tablas existen · ☐ El validador
 termina sin errores ni avisos · ☐ Las seis peticiones están en la colección, en orden · ☐
 Cada una responde lo que dice esta hoja · ☐ go test ./... en verde, con una prueba de
 integración · ☐ El código de la entidad principal está en capas y la estructura está explicada
 en docs/decisiones.md · ☐ Hay una captura de respaldo de cada paso
 Cambien ☐ por ☑ solo en lo que de verdad cumplen.




Hito 2 · Proyecto autónomo · Hilo Servidor · 2026-2                                                          5 / 10
```

## Página 6 de 10

```text
 5. La demostración, minuto a minuto
 El docente publica el orden de las parejas. La pareja va al escritorio del docente con su
 computadora y con todo ya encendido. Mientras una pareja demuestra, la siguiente de
 la lista está «en cubierta» y deja listo:
      la base de datos y el servidor corriendo;
      un registro recién creado con la petición 0, y su id puesto en las peticiones 4 y 5
      (apartado 6, «Datos de partida»);
      Bruno abierto en la petición 1;
      una terminal en la carpeta del proyecto, donde ya corrieron el validador una vez ese día.
 Cada pareja tiene cinco minutos, y son un tope:

   Minuto          Qué pasa

   0:00 –          Muestran git status (sin cambios pendientes, salvo el id de las peticiones 4 y 5
   0:25            dentro de docs/hito2_bruno/ ) y git log -1 (con fecha anterior al cierre). Se sortea
                   quién maneja el teclado en las demostraciones 1 y 2; el otro integrante maneja la 3, la
                   4 y la 5.

   0:25 –          Demostración 1 · El contrato.
   1:00

   1:00 –          Demostración 2 · Quién es. Peticiones 1, 2 y 3.
   1:50

   1:50 –          Demostración 3 · Qué puede hacer. Peticiones 4 y 5.
   2:35

   2:35 –          Demostración 4 · La regla del negocio. Petición 6 y go test ./... .
   3:20

   3:20 –          Demostración 5 · La estructura. Las carpetas y dos archivos.
   3:50

   3:50 –          Dos preguntas, una a cada integrante, de unos 30 segundos.
   4:50

   4:50 –          El docente registra las notas.
   5:00




Hito 2 · Proyecto autónomo · Hilo Servidor · 2026-2                                                     6 / 10
```

## Página 7 de 10

```text
 Qué dicen en cada demostración
 Quien tiene el teclado envía la petición y dice solo una frase, completando los corchetes. No
 se explica nada más: las preguntas vienen al final.

   Demostración                  Frase

   1. El contrato                «Este es nuestro contrato. El validador no encuentra nada. Aquí está
                                  security , y aquí el 401, el 403 y el 422.»

   2. Quién es                   «Sin token, [ruta] responde 401. Iniciamos sesión como [usuario] y
                                 recibimos el token. Con el token, responde 200.»

   3. Qué puede                  «[Usuario] es [rol] y no puede [acción]: 403. [Otro usuario] es [rol] y sí
   hacer                         puede: [código].»

   4. La regla del               «Pasar de [estado] a [estado] no está permitido en nuestro negocio: 422. Y
   negocio                       estas son nuestras pruebas, en verde.»

   5. La estructura              «Organizamos el código así: [carpetas]. Aquí se atiende la petición y aquí
                                 está la regla del cambio de estado: no se mezclan.»


 Las preguntas
 Se conocen desde ahora para que las preparen. El docente elige dos, una para cada
 integrante, y completa lo que va entre corchetes con algo de su demostración.
  1. ¿En qué línea de su código se decide que esta ruta pide token? Muéstrela.
  2. ¿Qué lleva su token por dentro? Si alguien le cambia el rol, ¿qué responde su servidor y
     por qué?
  3. ¿Por qué esta ruta es solo de [rol] en su negocio y no de los dos?
  4. Si mañana alguien saca esta ruta del grupo protegido, ¿cuál de sus pruebas avisa?
  5. En la transición prohibida, ¿por qué 422 y no 400 ni 403? ¿Dónde se valida?
  6. ¿Por qué eligieron esta estructura y no otra? Si mañana cambian de base de datos, ¿qué
      archivos tocan y cuáles no?

 Reglas del turno
      Las peticiones se envían desde la colección y en su orden. Lo único que se cambia ese
      día es el id del registro nuevo en las peticiones 4 y 5, y se cambia antes del turno, no
      durante.
      La entrega es el repositorio de la organización uleam-web-2026-2 . Se da por
      entregado si, a la hora del cierre, ese repositorio tiene la hoja del hito y el código que se
      va a demostrar. Si no, la pareja demuestra en horario de tutorías por el 50 %, igual que
      la que no llega a tiempo. Un repositorio personal o de fuera de la organización no
      cuenta, aunque tenga todo.
      Se demuestra lo que estaba en ese repositorio al cierre. Si su último commit es
      posterior al cierre, antes de empezar vuelven al último commit anterior a esa hora.
      Si una petición falla, hay un solo reintento de 20 segundos. Si no sale, abren la
      captura de respaldo y esa demostración vale la mitad como máximo.
      A los cinco minutos se corta. Lo que no se alcanzó a demostrar vale 0.




Hito 2 · Proyecto autónomo · Hilo Servidor · 2026-2                                                           7 / 10
```

## Página 8 de 10

```text
      Si la pareja está en el aula pero no está lista cuando la llaman, pasa al final de la lista
      una sola vez, sin descuento.
      No se pasa la pregunta al compañero.
      Si falta un integrante, el que está hace todo y responde las dos preguntas.
      Si la pareja no está en el aula cuando el docente la llama, demuestra en horario
      de tutorías y su nota es el 50 % de la que obtenga. Vale igual para la pareja que
      llega tarde ese mismo día. El día y la hora son los de la lista que publica el docente.
      Justificaciones. El docente no puede aceptarlas: se tramitan en la Dirección de
      Carrera. Con una justificación validada por la Dirección de Carrera, la pareja demuestra
      en la fecha que se le asigne y conserva su nota completa. Sin esa validación, vale el 50
      %.
      Hace falta una computadora por pareja. Quien no tenga puede demostrar desde
      Codespaces: avisen antes.


 6. Para que corra ese día
      Puerto de la base. Si su computadora tiene PostgreSQL instalado en Windows, ese
      programa ocupa el puerto 5432 y recibe las conexiones que iban al contenedor.
      Publiquen el contenedor en el 5433 y cambien la dirección de la base en su
      configuración.
      Puerto del servidor. Si el 8080 está ocupado, cambien el puerto del servidor y el de la
      colección de Bruno.
      Token que vence. La petición 2 se envía en vivo, así que el token es nuevo. No peguen
      en la colección un token de otro día.
      Datos de partida. La petición 5 gasta su registro: lo cambia de estado, y enviada
      otra vez sobre el mismo registro responde 422 en lugar de 200. Antes de cada ensayo,
      y antes de su turno, creen un registro nuevo con la petición 0 y pongan su id en las
      peticiones 4 y 5. La petición 6 no cambia nada: su registro solo tiene que estar en el
      estado desde el que parte su transición prohibida, y sirve para todos los ensayos.
      El validador. La primera vez descarga un paquete de internet. Córranlo ese mismo día
      antes de su turno, para que en la demostración responda de inmediato.
      Ensayen con reloj, desde cero: apaguen todo, enciendan y hagan la demostración
      completa en menos de cuatro minutos.


 7. Antes del cierre, revisen
 ☐ Los structs de todas las entidades del modelo están en el código, y al arrancar se crean
 todas sus tablas.
 ☐ docs/openapi.yaml declara security , tiene el 401, el 403 y el 422 donde corresponden, y
 lleva todos los endpoints del mapa del hito 1: los que aún no existen dicen «(plan)».
 ☐ redocly.yaml está en la raíz, sin la línea security-defined: off , y el validador termina
 sin errores ni avisos.
 ☐ docs/hito2_hoja.md tiene sus cuatro partes, con los títulos de este documento.
 ☐ docs/hito2_bruno/ tiene la petición 0 de preparación y las seis peticiones, numeradas y
 en orden.




Hito 2 · Proyecto autónomo · Hilo Servidor · 2026-2                                             8 / 10
```

## Página 9 de 10

```text
 ☐ Cada petición responde lo que dice la hoja.
 ☐ go test ./... termina en verde y hay al menos una prueba de integración.
 ☐ docs/hito2_capturas/ tiene una captura de cada respuesta y una de las pruebas.
 ☐ En el repositorio no hay contraseñas reales ni el secreto con que se firman los tokens.
 ☐ Los dos integrantes tienen commits, la IA está declarada y se hizo git push al
 repositorio de la organización uleam-web-2026-2 antes del cierre.
 ☐ El código de la entidad principal está organizado en capas, docs/decisiones.md explica la
 estructura, y después de moverlo las seis peticiones y las pruebas siguen respondiendo
 igual.
 ☐ Ensayaron con reloj y los dos pueden responder las seis preguntas.


 8. Rúbrica
 Cada demostración vale 14 % y cada pregunta 15 %. Todo se califica con la misma escala
 corta:

   Nivel      En una demostración                                             En una pregunta

   0          No corre, o no se alcanzó en los cinco minutos.                 No responde.

   1          Solo se ve en la captura de respaldo, o corre pero una de       Repite lo que está
              sus partes falla o no responde lo que dicen su hoja y su        escrito, sin explicar por
              contrato.                                                       qué.

   2          Corre en vivo, completa, y responde lo que dicen su hoja y su   Sostiene la decisión con
              contrato.                                                       una razón propia.




Hito 2 · Proyecto autónomo · Hilo Servidor · 2026-2                                                       9 / 10
```

## Página 10 de 10

```text
 Cada demostración tiene varias partes (la 2 son tres peticiones; la 4, una petición y las
 pruebas; la 5, dos archivos y la explicación escrita). Para el nivel 2 tienen que salir todas.

   Qué se                Peso       Qué se revisa
   califica

   1. El contrato        14         Validador sin errores ni avisos · security declarado · 401, 403 y 422 en
                         %          el contrato

   2. Quién es           14         401 sin token · el inicio de sesión entrega un token · 200 con el token
                         %

   3. Qué puede          14         403 con el rol que no puede · respuesta correcta con el rol que sí puede
   hacer                 %

   4. La regla           14         422 en la transición prohibida · go test ./... en verde, con una prueba
   del negocio           %          de integración

   5. La                 14         La función que atiende la petición no consulta la base directamente · la
   estructura            %          regla del cambio de estado no depende de HTTP · coincide con lo escrito
                                    en docs/decisiones.md · vale cualquier estructura que cumpla esto, no
                                    solo las de clase

   Pregunta al           15         Una de las seis preguntas publicadas
   integrante 1          %

   Pregunta al           15         Otra de las seis preguntas
   integrante 2          %

 Se califica la razón, no el acierto: una respuesta que sostiene con un argumento del
 negocio o del diseño una decisión que el docente habría tomado distinto vale 2.
 Cómo se calcula la nota. Cada parte aporta nivel ÷ 2 × peso. La suma da la nota del hito
 sobre 100, que equivale al 7 % de la nota del período.
 Ejemplo: demostraciones en 2, 2, 1, 2 y 2, y preguntas en 2 y 1, dan 14 + 14 + 7 + 14 +
 14 + 15 + 7,5 = 85,5 sobre 100.
 Fuera de fecha. La demostración en horario de tutorías (por no llegar a tiempo o por no
 tener la entrega en el repositorio de la organización al cierre) se califica con esta misma
 rúbrica y el resultado se multiplica por 0,5: un 85,5 queda en 42,75. Con una justificación
 validada por la Dirección de Carrera no se multiplica.




Hito 2 · Proyecto autónomo · Hilo Servidor · 2026-2                                                            10 / 10
```
