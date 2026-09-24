# ADR-0001: Monolito modular con arquitectura por paquetes de dominio

- **Estado:** aceptado
- **Fecha:** 2026-09-24

## Contexto
Habitat centraliza inquilinos, unidades, zonas comunes, reservas, tickets e historial de un edificio.
Es un proyecto de un solo desarrollador que debe poder evolucionar a algo más grande.

## Opciones consideradas
1. **Microservicios desde el inicio** — complejidad operativa enorme sin beneficio real a esta escala.
2. **Capas horizontales** (`handlers/`, `services/`, `repositories/`) — fácil al principio, pero cada feature toca todo el árbol.
3. **Monolito modular por dominio** (`internal/tenant`, `internal/reservation`...) — cada paquete tiene su dominio, servicio, repositorio y handler.

## Decisión
Opción 3. Cada módulo expone un `Service` y una interfaz `Repository`. Los módulos no importan
los repositorios de otros módulos; si necesitan datos ajenos, llaman a su `Service` o reaccionan a eventos.

## Consecuencias
- Fácil de extraer un módulo a un servicio aparte si algún día hace falta (T-23).
- Hay que vigilar el acoplamiento entre paquetes (ciclos de import = señal de mal diseño).
