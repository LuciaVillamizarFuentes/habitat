# 🏢 Habitat

**Plataforma centralizada para la gestión de edificios residenciales**: inquilinos, unidades, zonas comunes,
reservas, tickets de mantenimiento, cuotas de administración e historial completo de lo que pasa en el edificio.

Construido en **Go** con la librería estándar como base, PostgreSQL y una arquitectura de monolito modular
([ADR-0001](docs/adr/0001-monolito-modular.md)).

> Proyecto de portafolio construido tarea a tarea siguiendo un [roadmap de Junior a Senior](docs/ROADMAP.md).

## Módulos

| Módulo | Qué gestiona | Estado |
|---|---|---|
| `tenant` | Inquilinos / residentes | ✅ Módulo de referencia (completo, con tests) |
| `unit` | Apartamentos, parqueaderos, depósitos, locales | 🧱 Dominio definido |
| `resource` | Zonas comunes: salón social, BBQ, gimnasio… | 🧱 Dominio definido |
| `reservation` | Reservas de zonas comunes sin doble reserva | 🧱 Dominio definido |
| `ticket` | Mantenimiento, PQRS, SLA por prioridad | 🧱 Dominio definido |
| `audit` | Historial append-only de eventos | 🧱 Dominio definido |

## Arrancar

```bash
# Requisitos: Go 1.24+, Docker
make run          # API en :8080 con repositorio en memoria
make test-race    # tests con detector de carreras
make up           # Postgres + API en Docker
```

```bash
curl -X POST localhost:8080/api/v1/tenants \
  -H 'Content-Type: application/json' \
  -d '{"full_name":"Ana Gómez","email":"ana@example.com"}'

curl localhost:8080/api/v1/tenants
```

## Estructura

```
cmd/api/              → main: configuración, wiring y graceful shutdown
internal/
  config/             → variables de entorno
  platform/httpx/     → JSON, errores y middlewares compartidos
  tenant/             → dominio + service + repository + handler + tests (REFERENCIA)
  unit/ resource/ reservation/ ticket/ audit/   → módulos por construir
migrations/           → SQL
docs/
  ROADMAP.md          → las tareas
  adr/                → decisiones de arquitectura
.claude/skills/mentor-go/  → skill de mentor para revisar cada tarea
```

## Flujo de trabajo

1. Toma la siguiente tarea de [`docs/ROADMAP.md`](docs/ROADMAP.md) y busca sus `TODO(T-XX)` en el código.
2. Rama `feat/T-XX-...` → implementa → tests → PR.
3. Pide revisión al mentor: `/mentor-go T-XX`.

## Progreso

- [ ] 🟢 Junior: T-00 → T-06
- [ ] 🟡 Mid: T-07 → T-14
- [ ] 🟠 Senior: T-15 → T-22
- [ ] 🔴 Senior+: T-23 → T-28
