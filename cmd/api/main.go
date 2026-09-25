// Command api arranca el servidor HTTP de Habitat.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tu-usuario/habitat/internal/config"
	"github.com/tu-usuario/habitat/internal/platform/httpx"
	"github.com/tu-usuario/habitat/internal/tenant"
	"github.com/tu-usuario/habitat/internal/unit"
)

func main() {
	if err := run(); err != nil {
		slog.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run() error {
	cfg := config.Load()
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	// --- Wiring de dependencias (composition root) ---
	// TODO(T-08): usar PostgresRepository cuando DATABASE_URL esté configurado.
	tenantRepo := tenant.NewMemoryRepository()
	tenantSvc := tenant.NewService(tenantRepo)

	unitRepo := unit.NewMemoryRepository()
	unitSvc := unit.NewService(unitRepo)

	var handler http.Handler = newMux(tenantSvc, unitSvc)
	handler = httpx.Logger(log, handler)
	handler = httpx.Recoverer(log, handler)

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Graceful shutdown: esperamos SIGINT/SIGTERM y cerramos conexiones en curso.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		log.Info("server listening", "addr", cfg.HTTPAddr, "env", cfg.Env)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		log.Info("shutting down")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}

// newMux registra todas las rutas de la API. Separado de run() para poder
// testear el wiring sin levantar un servidor real.
func newMux(tenantSvc *tenant.Service, unitSvc *unit.Service) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		httpx.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	// TODO(T-07): GET /readyz que verifique la conexión a la base de datos.
	tenant.NewHandler(tenantSvc).Routes(mux)
	unit.NewHandler(unitSvc).Routes(mux)
	return mux
}
