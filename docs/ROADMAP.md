# 🗺️ Roadmap Habitat: de Junior a Senior en Go

Cada tarea es un **reto real** que te podrías encontrar en un equipo backend. Están ordenadas por nivel de
seniority: las primeras practican el lenguaje, las últimas practican **criterio** (diseño, trade-offs, operación).

## Cómo trabajar cada tarea

1. Crea una rama `feat/T-XX-descripcion`.
2. Lee la teoría sugerida **antes** de programar (máx. 30–45 min; es para desbloquearte, no para postergar).
3. Implementa hasta cumplir los **criterios de aceptación**.
4. Responde por escrito las **preguntas de entrevista** en `docs/journal/T-XX.md` (esto es oro para tus entrevistas).
5. Abre un PR contra `main` (aunque trabajes solo) y pide revisión al **mentor**: `/mentor-go T-XX`.
6. Merge solo cuando el mentor dé "Aprobado" o "Aprobado con comentarios".

> Busca los `TODO(T-XX)` en el código: `grep -rn "TODO(T-" .`

**Leyenda de nivel:** 🟢 Junior · 🟡 Mid · 🟠 Senior · 🔴 Senior+/Staff

---

## 🟢 Nivel 1 — Junior: domina el lenguaje y el patrón del proyecto

### T-00 · Entiende el template
- **Qué hacer:** clona, corre `make run` y `make test-race`. Crea un inquilino con `curl`. Lee el módulo `internal/tenant` completo y dibuja (a mano o en Excalidraw) cómo viaja un request: `main → middleware → handler → service → repository`.
- **Aceptación:** diagrama en `docs/architecture.png` + un párrafo explicando por qué el `Service` recibe una interfaz y no el `MemoryRepository` concreto.
- **Preguntas de entrevista:** ¿Qué es la "composition root"? ¿Por qué `internal/`? ¿Qué pasa si quitas `DisallowUnknownFields`?
- **Teoría:**
  - Platzi · [Struct Server](https://platzi.com/cursos/go-rest-websockets/struct-server/) y [Patrón repository](https://platzi.com/cursos/go-rest-websockets/patron-repository/)
  - Docs · [Effective Go](https://go.dev/doc/effective_go) (secciones Interfaces y Errors)
  - YouTube · [Kat Zien — How Do You Structure Your Go Apps (GopherCon 2018)](https://www.youtube.com/watch?v=oL6JBUk6tj0)

### T-01 · Validación del dominio `Unit`
- **Qué hacer:** implementa `Validate()` para `unit.Unit` con las reglas del `TODO`. Errores de dominio como variables (`var ErrInvalidKind = errors.New(...)`).
- **Aceptación:** tests table-driven con al menos 8 casos, incluidos bordes (coeficiente = 1, área = 0, kind desconocido).
- **Preguntas:** ¿`errors.Is` vs `==`? ¿Cuándo usar un tipo de error propio (`type ValidationError struct`) en vez de un sentinel?
- **Teoría:**
  - Docs · [Working with Errors in Go 1.13](https://go.dev/blog/go1.13-errors)
  - Docs · [Table-driven tests (Go Wiki)](https://go.dev/wiki/TableDrivenTests)
  - Platzi · [Curso de Unit Testing en Go](https://platzi.com/cursos/go-unit-testing/)

### T-02 · CRUD completo de unidades
- **Qué hacer:** `Repository` (memoria), `Service`, `Handler` para `units` copiando el patrón de `tenant`. Rutas: `POST/GET /api/v1/units`, `GET /api/v1/units/{id}`.
- **Aceptación:** tests de handler con `httptest`; `code` duplicado → `409`; datos inválidos → `422`.
- **Preguntas:** ¿`RWMutex` vs `Mutex`, cuándo compensa? ¿Qué hace el router de Go 1.22 con `{id}`?
- **Teoría:**
  - Platzi · [CRUD para Posts](https://platzi.com/cursos/go-rest-websockets/crud-para-posts/)
  - Platzi · [Lectura y Escritura Concurrente con RWMutex](https://platzi.com/cursos/golang-avanzado/mutex-de-lectura-y-escritura/)
  - Docs · [Routing Enhancements for Go 1.22](https://go.dev/blog/routing-enhancements)

### T-03 · Actualizar y dar de baja inquilinos
- **Qué hacer:** `PATCH /api/v1/tenants/{id}` (actualización parcial) y `DELETE` como **soft delete** (`active=false`).
- **Reto real:** en un `PATCH`, ¿cómo distingues "el cliente no envió `phone`" de "el cliente quiere `phone` vacío"? Resuélvelo (pista: punteros o un tipo `Optional[T]` con genéricos).
- **Aceptación:** tests que cubran ambos casos; un inquilino dado de baja no aparece en el listado por defecto.
- **Preguntas:** ¿Por qué soft delete y no borrado físico en un sistema con historial? ¿Implicaciones legales (Habeas Data, Ley 1581)?
- **Teoría:** Docs · [Tutorial de genéricos](https://go.dev/doc/tutorial/generics)

### T-04 · Paginación y filtros
- **Qué hacer:** `GET /api/v1/tenants?limit=20&cursor=...&active=true&unit_id=...`.
- **Reto real:** implementa **paginación por cursor** (no offset). Devuelve `next_cursor`.
- **Aceptación:** `limit` máx. 100; cursor opaco (base64); test que demuestre que no se repiten ni saltan elementos al insertar durante la paginación.
- **Preguntas:** ¿Por qué offset es lento en tablas grandes? ¿Qué problemas de consistencia tiene?
- **Teoría:** Platzi · [Paginación para Posts](https://platzi.com/cursos/go-rest-websockets/paginacion-para-posts/)

### T-05 · Relación inquilino ↔ unidad
- **Qué hacer:** al registrar un inquilino con `unit_id`, verificar que la unidad existe y es de tipo `apartment`.
- **Reto de diseño:** `tenant` no debe importar el repositorio de `unit`. Define una interfaz pequeña **en el paquete consumidor** (`type UnitChecker interface { Exists(ctx, id) (bool, error) }`).
- **Aceptación:** test con un fake de `UnitChecker`; cero ciclos de import.
- **Preguntas:** "Accept interfaces, return structs": explícalo con este caso.
- **Teoría:** YouTube · [Rob Pike — Go Proverbs](https://www.youtube.com/watch?v=PAAkCSZUG1c)

### T-06 · IDs con UUIDv7 + tu primer ADR
- **Qué hacer:** reemplaza `newID()` por UUIDv7 (`github.com/google/uuid`). Escribe `docs/adr/0002-identificadores.md` comparando: autoincrement, UUIDv4, UUIDv7, ULID.
- **Aceptación:** ADR con contexto, opciones y consecuencias (impacto en índices B-tree).
- **Preguntas:** ¿Por qué UUIDv4 fragmenta índices? ¿Es un riesgo exponer IDs secuenciales?

---

## 🟡 Nivel 2 — Mid: persistencia real, errores bien hechos, dominio con reglas

### T-07 · Configuración robusta y health checks
- **Qué hacer:** valida la config al arrancar (falla rápido). Añade `GET /readyz` que haga ping a la base de datos con timeout de 2s.
- **Aceptación:** `/healthz` (liveness) nunca toca la BD; `/readyz` (readiness) sí. Explícalo en el README.
- **Preguntas:** diferencia liveness vs readiness en Kubernetes; ¿qué pasa si tu liveness depende de la BD y la BD cae?
- **Teoría:** Docs · [Go Concurrency Patterns: Context](https://go.dev/blog/context)

### T-08 · Postgres de verdad
- **Qué hacer:** `PostgresRepository` para `tenant` y `unit` con [pgx](https://github.com/jackc/pgx) (o [sqlc](https://sqlc.dev) si quieres generar el código). Migraciones con [goose](https://github.com/pressly/goose) o `golang-migrate`. `main.go` elige la implementación según la config.
- **Reto real:** traduce el error de unique violation de Postgres (`23505`) a `tenant.ErrDuplicateMail` **sin** que el servicio sepa que existe Postgres.
- **Aceptación:** los mismos tests del servicio pasan con ambas implementaciones; pool de conexiones configurado (max conns, timeouts).
- **Preguntas:** ¿`database/sql` vs pgx nativo? ¿Qué pasa si no cierras `rows`? ¿Cómo dimensionas el pool?
- **Teoría:**
  - Docs · [Accessing a relational database (tutorial oficial)](https://go.dev/doc/tutorial/database-access)
  - Platzi · [Implementando el registro](https://platzi.com/cursos/go-rest-websockets/implementando-el-registro/) (repositorio con Postgres)
  - Platzi · [Implementación de Singleton para conexiones de BD](https://platzi.com/cursos/golang-avanzado/singleton/) — y luego pregúntate por qué un singleton global **no** es buena idea aquí.

### T-09 · Zonas comunes (resources)
- **Qué hacer:** CRUD de `resource` con horarios de apertura, capacidad y tarifa en **centavos (`int64`)**.
- **Reto real:** los horarios se guardan en hora local del edificio (`America/Bogota`). Usa `time.LoadLocation` y no mezcles UTC con local.
- **Aceptación:** tests que prueben un horario 08:00–22:00 en Bogotá con instantes UTC.
- **Preguntas:** ¿Por qué nunca `float64` para dinero? ¿Por qué guardar `TIMESTAMPTZ`?

### T-10 · Reservas: la versión ingenua
- **Qué hacer:** `POST /api/v1/reservations` que valide horario, duración máxima y **no solapamiento** (consultando reservas existentes antes de insertar).
- **Aceptación:** tests de `Overlaps` con rangos contiguos, contenidos, idénticos; endpoint funcionando.
- **Nota:** esta versión **tiene un bug de concurrencia a propósito**. Lo vas a romper en T-17. 😈
- **Teoría:** Docs · [PostgreSQL Range Types](https://www.postgresql.org/docs/current/rangetypes.html) (léelo ahora, lo usarás luego)

### T-11 · Request ID, logging estructurado y `context`
- **Qué hacer:** middleware que lee `X-Request-ID` o genera uno, lo guarda en el `context` y lo incluye en **todos** los logs de ese request (incluidos los del servicio).
- **Aceptación:** un `slog.Handler` que extraiga el request ID del contexto automáticamente.
- **Preguntas:** ¿Por qué no usar `string` como key de contexto? ¿Qué **no** debería ir en un context?
- **Teoría:**
  - Docs · [log/slog](https://pkg.go.dev/log/slog)
  - Platzi · [Middleware de validación de autenticación](https://platzi.com/cursos/go-rest-websockets/middleware-de-validacion-de-autenticacion/)

### T-12 · Errores de API estándar (RFC 9457)
- **Qué hacer:** migra `httpx.Error` a `application/problem+json` con `type`, `title`, `status`, `detail`, `instance`. Errores de validación con lista de campos.
- **Aceptación:** los errores 500 **nunca** filtran detalles internos; se loguean con el request ID.
- **Teoría:** [RFC 9457 — Problem Details for HTTP APIs](https://www.rfc-editor.org/rfc/rfc9457)

### T-13 · Tickets de mantenimiento con máquina de estados
- **Qué hacer:** módulo `ticket` completo. Las transiciones inválidas devuelven `409`.
- **Reto de diseño:** modela la máquina de estados como datos (mapa de transiciones permitidas), no como un `switch` gigante.
- **Aceptación:** test que recorre **todas** las combinaciones estado×estado y verifica cuáles se permiten.
- **Teoría:** Platzi · [Patrón Strategy en Go](https://platzi.com/cursos/golang-avanzado/strategy/)

### T-14 · SLA y asignación automática
- **Qué hacer:** según la prioridad, calcula `due_at` (urgent 4h, high 24h, medium 72h, low 7d **hábiles**). Endpoint `GET /api/v1/tickets/overdue`.
- **Reto real:** días hábiles en Colombia excluyen festivos. Diseña una interfaz `Calendar` inyectable para no acoplar el dominio a una fuente de festivos.
- **Aceptación:** tests deterministas inyectando el reloj (`now func() time.Time`, como en `tenant.Service`).

---

## 🟠 Nivel 3 — Senior: concurrencia, consistencia, seguridad, operación

### T-15 · Historial del edificio (audit log)
- **Qué hacer:** tabla `audit_events` append-only. Cada caso de uso relevante registra un evento (`tenant.registered`, `reservation.confirmed`, `ticket.status_changed`...). Endpoint `GET /api/v1/history?aggregate_id=...`.
- **Reto real:** el evento debe guardarse **en la misma transacción** que el cambio. Diseña un `TxManager` / Unit of Work sin que los servicios importen `pgx`.
- **Aceptación:** si falla el insert del evento, el cambio de negocio hace rollback (test de integración).
- **Preguntas:** ¿Cómo pasas la transacción entre repositorios sin contaminar las interfaces? (hay varias escuelas: tx en context, funciones `WithTx`, repos transaccionales; defiende una).

### T-16 · Tests de integración y CI serio
- **Qué hacer:** tests de integración contra Postgres real con [Testcontainers for Go](https://golang.testcontainers.org/) (build tag `integration`). CI con lint, `-race`, integración y cobertura.
- **Aceptación:** pipeline verde; cobertura ≥ 75% en `internal/`; un fuzz test para el parser de cursor (T-04).
- **Teoría:**
  - Docs · [Tutorial de Fuzzing](https://go.dev/doc/tutorial/fuzz)
  - Platzi · [¿Qué son los mocks? ¿Cómo y cuándo mockear?](https://platzi.com/cursos/go-unit-testing/) (clase del curso de Unit Testing)

### T-17 · 💥 Rompe y arregla la doble reserva
- **Qué hacer:**
  1. Escribe un test que lance 50 goroutines intentando reservar el **mismo** salón a la **misma** hora. Demuestra que con T-10 se crean varias reservas.
  2. Arréglalo a nivel de base de datos con `EXCLUDE USING gist (resource_id WITH =, tstzrange(starts_at, ends_at) WITH &&) WHERE (status <> 'cancelled')`.
  3. Añade **optimistic locking** (`version`) para la edición de reservas.
- **Aceptación:** el test de 50 goroutines crea exactamente 1 reserva, las demás reciben `409`. Documenta en un ADR por qué no bastaba con un `sync.Mutex` en Go (pista: ¿y si corren 3 réplicas del servicio?).
- **Preguntas:** optimistic vs pessimistic locking (`SELECT ... FOR UPDATE`); niveles de aislamiento; ¿qué es un *write skew*?
- **Teoría:**
  - Platzi · [Condiciones de Carrera en Programación Concurrente](https://platzi.com/cursos/golang-avanzado/condiciones-de-carrera-en-programacion-c/) y [Evitar Condiciones de Carrera con Mutex y WaitGroup](https://platzi.com/cursos/golang-avanzado/sync-mutex-lock-y-unlock/)
  - Docs · [Data Race Detector](https://go.dev/doc/articles/race_detector)
  - Docs · [PostgreSQL: Exclusion Constraints](https://www.postgresql.org/docs/current/ddl-constraints.html#DDL-CONSTRAINTS-EXCLUSION)

### T-18 · Autenticación y autorización (RBAC)
- **Qué hacer:** login con JWT (access corto + refresh token rotado). Roles: `admin`, `staff` (portería/mantenimiento), `tenant`. Un inquilino solo ve **sus** reservas y tickets.
- **Reto real:** la autorización a nivel de recurso ("¿esta reserva es tuya?") vive en el servicio, no solo en el middleware. Test que pruebe que un tenant **no** puede ver la reserva de otro (IDOR).
- **Aceptación:** contraseñas con bcrypt/argon2id; refresh tokens revocables; tests de autorización por rol.
- **Teoría:**
  - Platzi · [Autenticación de usuarios](https://platzi.com/cursos/go-rest-websockets/autenticacion-de-usuarios/) e [Implementando el middleware](https://platzi.com/cursos/go-rest-websockets/implementando-el-middleware/)
  - [OWASP API Security Top 10](https://owasp.org/API-Security/) (lee API1: BOLA)

### T-19 · Outbox pattern + eventos
- **Qué hacer:** cuando se confirma una reserva, hay que notificar (email/push). Implementa **transactional outbox**: el evento se escribe en `outbox` en la misma transacción; un worker lo publica (a NATS, o a un log por ahora) y lo marca como enviado.
- **Aceptación:** entrega *at-least-once*; los consumidores son idempotentes; el worker usa `SELECT ... FOR UPDATE SKIP LOCKED` para poder correr varias réplicas.
- **Preguntas:** ¿Por qué no publicar a NATS directamente después del commit? ¿Qué es exactly-once y por qué es (casi) un mito?
- **Teoría:**
  - [Pattern: Transactional outbox (microservices.io)](https://microservices.io/patterns/data/transactional-outbox.html)
  - Platzi · [Microservicios y arquitecturas basadas en eventos](https://platzi.com/cursos/go-eventos-cqrs/microservicios-y-arquitecturas-basadas-en-eventos/) y [Definiendo mensajes y eventos](https://platzi.com/cursos/go-eventos-cqrs/definiendo-mensajes-y-eventos/)

### T-20 · Rate limiting y caché
- **Qué hacer:** rate limit por usuario en endpoints de escritura (token bucket, `golang.org/x/time/rate` en memoria → luego Redis para múltiples réplicas). Caché de disponibilidad de zonas comunes con invalidación al reservar.
- **Aceptación:** respuesta `429` con header `Retry-After`; test que demuestre la invalidación de caché.
- **Preguntas:** cache-aside vs write-through; ¿qué es un *cache stampede* y cómo lo evitas? (pista: `singleflight`).
- **Teoría:** Platzi · [Creación de un Sistema de Caché Concurrente](https://platzi.com/cursos/golang-avanzado/sistema-de-cache-sin-concurrencia/) y [Cache concurrente para cálculos intensivos](https://platzi.com/cursos/golang-avanzado/reutilizacion-de-computacion-intensiva/)

### T-21 · Jobs en segundo plano con worker pool
- **Qué hacer:** reservas `pending` que no se pagan en 30 min se cancelan solas; recordatorios 24h antes de cada reserva; tickets vencidos escalan prioridad.
- **Reto real:** scheduler + worker pool con número de workers configurable, **graceful shutdown** (terminar el trabajo en curso al recibir SIGTERM) y sin goroutine leaks.
- **Aceptación:** test con `goleak` (o equivalente) que verifique que no quedan goroutines; test del shutdown.
- **Teoría:**
  - Docs · [Go Concurrency Patterns: Pipelines and cancellation](https://go.dev/blog/pipelines)
  - YouTube · [Rob Pike — Concurrency is not Parallelism](https://www.youtube.com/watch?v=oV9rvDllKEg)
  - Platzi · [Curso de Go Intermedio](https://platzi.com/cursos/golang-intermedio/) (worker pools)

### T-22 · Notificaciones en tiempo real
- **Qué hacer:** la portería ve en vivo los tickets urgentes y las reservas del día. Implementa WebSockets **o** SSE (justifica en un ADR cuál y por qué).
- **Reto real:** con 3 réplicas del servicio, un evento generado en la réplica A debe llegar a clientes conectados a la réplica C (fan-out vía NATS/Redis pub-sub). Maneja backpressure: un cliente lento no puede bloquear a los demás.
- **Teoría:**
  - Platzi · [Websockets](https://platzi.com/cursos/go-rest-websockets/websockets/), [Struct de Hub para conexiones](https://platzi.com/cursos/go-rest-websockets/struct-de-hub-para-conexiones/), [Implementando el Broadcast](https://platzi.com/cursos/go-rest-websockets/implementando-el-broadcast/)
  - Platzi · [Creando el Servicio Pusher](https://platzi.com/cursos/go-eventos-cqrs/creando-el-servicio-pusher/)

---

## 🔴 Nivel 4 — Senior+ / Staff: sistemas, rendimiento y decisiones

### T-23 · Observabilidad completa
- **Qué hacer:** OpenTelemetry: trazas (HTTP → servicio → SQL), métricas RED (rate, errors, duration) y logs correlacionados con trace ID. `pprof` detrás de auth. Docker compose con Grafana/Tempo/Prometheus (o Jaeger).
- **Aceptación:** puedes seguir un request lento desde el dashboard hasta la query SQL culpable. Define 2 SLOs (ej: p99 de `POST /reservations` < 300 ms).
- **Teoría:**
  - [OpenTelemetry Go](https://opentelemetry.io/docs/languages/go/)
  - Docs · [Diagnostics (profiling, tracing)](https://go.dev/doc/diagnostics)

### T-24 · Pruebas de carga y optimización
- **Qué hacer:** escenario con [k6](https://grafana.com/docs/k6/latest/): 500 usuarios virtuales consultando disponibilidad y reservando. Encuentra el cuello de botella con `pprof` y mejóralo (índices, N+1, allocations, tamaño del pool).
- **Aceptación:** informe antes/después con números (p50/p95/p99, RPS, CPU, allocs) en `docs/perf/`. Al menos un benchmark `go test -bench` con `-benchmem`.
- **Teoría:** Docs · [Profiling Go Programs](https://go.dev/blog/pprof)

### T-25 · Cuotas de administración e idempotencia
- **Qué hacer:** módulo de cobros: cuota mensual por unidad según **coeficiente de copropiedad**, pagos, saldo y paz y salvo. Endpoint de pago con header `Idempotency-Key`.
- **Reto real:** repartir $10.000.000 entre unidades con coeficientes debe sumar **exactamente** $10.000.000 (reparto de centavos sobrantes). Un pago reintentado 3 veces por timeout se cobra **una** vez.
- **Aceptación:** property-based test del reparto; test de reintentos concurrentes con la misma key.
- **Teoría:** [Stripe — Designing robust and predictable APIs with idempotency](https://stripe.com/blog/idempotency)

### T-26 · Multi-edificio (SaaS multi-tenant)
- **Qué hacer:** una administradora gestiona 30 edificios. Añade `building_id` a todo y garantiza aislamiento: un admin del edificio A jamás ve datos del B.
- **Reto de diseño:** compara 3 estrategias (BD por cliente, schema por cliente, columna + Row-Level Security) en un ADR e implementa una. Si eliges RLS, fija `app.building_id` por transacción.
- **Aceptación:** test automatizado de fuga entre edificios en **todos** los endpoints de lectura.
- **Teoría:** [PostgreSQL Row Security Policies](https://www.postgresql.org/docs/current/ddl-rowsecurity.html)

### T-27 · Extraer un servicio (y saber cuándo NO hacerlo)
- **Qué hacer:** extrae `notifications` a un servicio independiente comunicado por gRPC (comandos) y NATS (eventos). Protobuf versionado, timeouts, retries con backoff, circuit breaker.
- **Aceptación:** si el servicio de notificaciones cae, reservar sigue funcionando. ADR final: "¿valió la pena?" con argumentos honestos.
- **Teoría:**
  - Platzi · [Curso de Go Avanzado: Protobuffers y gRPC](https://platzi.com/cursos/go-protobuffers-grpc/)
  - Platzi · [CQRS: Command Query Responsibility Segregation](https://platzi.com/cursos/go-eventos-cqrs/cqrs-command-query-responsibility-segregation/)

### T-28 · Design doc + presentación (el entregable de un senior)
- **Qué hacer:** escribe `docs/DESIGN.md` como si fueras a presentarlo a un equipo: problema, arquitectura, modelo de datos, decisiones clave (enlaza tus ADRs), riesgos y qué harías con 3 meses más. Graba un video de 5 min explicándolo.
- **Aceptación:** alguien que no conoce el proyecto entiende en 10 min por qué está construido así.
- **Por qué importa:** en entrevistas senior te van a pedir **system design** y "cuéntame de una decisión técnica difícil". Este documento es tu respuesta preparada.

---

## 📚 Rutas completas (si quieres ir más allá de clases sueltas)
- Platzi · [Ruta Desarrollo Backend con Go](https://platzi.com/ruta/backend-go/)
- Platzi · [Curso de Go Avanzado: Concurrencia y Patrones de Diseño](https://platzi.com/cursos/golang-avanzado/)
- Platzi · [Curso de Go Avanzado: REST y WebSockets](https://platzi.com/cursos/go-rest-websockets/)
- Platzi · [Curso de Go Avanzado: Arquitectura de Eventos y CQRS](https://platzi.com/cursos/go-eventos-cqrs/)
- Blog · [Mat Ryer — How I write HTTP services in Go after 13 years](https://grafana.com/blog/how-i-write-http-services-in-go-after-13-years/)
- Docs · [A Tour of Go](https://go.dev/tour/) (para repasar rápido lo que se haya oxidado)
