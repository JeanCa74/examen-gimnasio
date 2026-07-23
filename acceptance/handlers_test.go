// ARCHIVO BLOQUEADO — NO MODIFICAR
package acceptance

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/joancema/examen-gimnasio/internal/models"
)

// sembrarHTTP prepara datos base directamente en la base de datos.
func sembrarHTTP(t *testing.T, db *gorm.DB) (models.Clase, models.Cliente, models.Inscripcion) {
	t.Helper()
	principal := models.Clase{Nombre: "Clase de spinning", PrecioUnitario: 10, Stock: 5, Activo: true}
	require.NoError(t, db.Create(&principal).Error)
	ana := models.Cliente{Nombre: "Ana Zambrano", Cedula: "1310000001", Telefono: "0990000001"}
	require.NoError(t, db.Create(&ana).Error)
	asistida := models.Inscripcion{
		ClaseID: principal.ID,
		ClienteID:     ana.ID,
		Cantidad:      1,
		Estado:        models.EstadoAsistida,
		Total:         10,
	}
	require.NoError(t, db.Create(&asistida).Error)
	return principal, ana, asistida
}

func postJSON(router http.Handler, ruta, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, ruta, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

// TestCP3_CrearInscripcionHTTP: el flujo feliz responde 201.
func TestCP3_CrearInscripcionHTTP(t *testing.T) {
	db := nuevaDB(t)
	router := nuevoRouterCompleto(t, db)
	principal, ana, _ := sembrarHTTP(t, db)

	body := fmt.Sprintf(`{"clase_id":%d,"cliente_id":%d,"cantidad":2}`, principal.ID, ana.ID)
	rec := postJSON(router, "/api/v1/inscripciones", body)
	require.Equal(t, http.StatusCreated, rec.Code,
		"POST /inscripciones válido debe responder 201. Body: %s", rec.Body.String())
}

// TestCP3_MapeoErroresHTTP: los errores de dominio se mapean a los status
// codes correctos (422 / 409 / 404).
func TestCP3_MapeoErroresHTTP(t *testing.T) {
	db := nuevaDB(t)
	router := nuevoRouterCompleto(t, db)
	principal, ana, asistida := sembrarHTTP(t, db)

	// Referencia inválida -> 422
	body := fmt.Sprintf(`{"clase_id":99999,"cliente_id":%d,"cantidad":1}`, ana.ID)
	rec := postJSON(router, "/api/v1/inscripciones", body)
	require.Equal(t, http.StatusUnprocessableEntity, rec.Code,
		"clase inexistente debe responder 422. Body: %s", rec.Body.String())

	// Stock insuficiente -> 409
	body = fmt.Sprintf(`{"clase_id":%d,"cliente_id":%d,"cantidad":99}`, principal.ID, ana.ID)
	rec = postJSON(router, "/api/v1/inscripciones", body)
	require.Equal(t, http.StatusConflict, rec.Code,
		"cantidad mayor al stock debe responder 409. Body: %s", rec.Body.String())

	// Retirar una ASISTIDA -> 409
	rec = postJSON(router, fmt.Sprintf("/api/v1/inscripciones/%d/retirar", asistida.ID), "")
	require.Equal(t, http.StatusConflict, rec.Code,
		"retirar una inscripcion ASISTIDA debe responder 409. Body: %s", rec.Body.String())

	// Obtener inexistente -> 404
	req := httptest.NewRequest(http.MethodGet, "/api/v1/inscripciones/99999", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	require.Equal(t, http.StatusNotFound, rr.Code,
		"GET de una inscripcion inexistente debe responder 404. Body: %s", rr.Body.String())
}
