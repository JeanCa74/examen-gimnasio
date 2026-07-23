package handlers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/joancema/examen-gimnasio/internal/handlers"
	"github.com/joancema/examen-gimnasio/internal/models"
	"github.com/joancema/examen-gimnasio/internal/services"
	"github.com/joancema/examen-gimnasio/internal/storage"
)

// routerEnMemoria arma el MISMO router de la API (rutas de routes.go) pero
// cableado con los repositorios en memoria, para probar la capa HTTP sin
// tocar una base de datos real.
func routerEnMemoria() (http.Handler, *storage.ClaseMemoria, *storage.ClienteMemoria) {
	claseRepo := storage.NuevaClaseMemoria()
	clienteRepo := storage.NuevoClienteMemoria()
	inscRepo := storage.NuevaInscripcionMemoria()

	router := handlers.NuevoRouter(
		handlers.NuevoClaseHandler(services.NuevaClaseService(claseRepo)),
		handlers.NuevoClienteHandler(services.NuevoClienteService(clienteRepo)),
		handlers.NuevaInscripcionHandler(services.NuevaInscripcionService(inscRepo, claseRepo, clienteRepo)),
	)
	return router, claseRepo, clienteRepo
}

// TestInscripcionHandler_Crear201 es un test PROPIO de la capa HTTP con
// httptest: manda un POST /api/v1/inscripciones real (sin abrir un puerto) y
// verifica que el flujo feliz responde 201 y que el body trae el total ya
// calculado por el service (3 × 8.50 = 25.50, sin descuento).
func TestInscripcionHandler_Crear201(t *testing.T) {
	router, claseRepo, clienteRepo := routerEnMemoria()

	clase := models.Clase{Nombre: "Clase de spinning", PrecioUnitario: 8.50, Stock: 10, Activo: true}
	require.NoError(t, claseRepo.Crear(&clase))
	cliente := models.Cliente{Nombre: "Ana Zambrano", Cedula: "1310000001"}
	require.NoError(t, clienteRepo.Crear(&cliente))

	body := fmt.Sprintf(`{"clase_id":%d,"cliente_id":%d,"cantidad":3}`, clase.ID, cliente.ID)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/inscripciones", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code, "POST válido debe responder 201. Body: %s", rec.Body.String())

	var creada models.Inscripcion
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &creada))
	require.InDelta(t, 25.50, creada.Total, 0.001, "el 201 debe incluir el total calculado")
	require.Equal(t, models.EstadoPendiente, creada.Estado)
}
