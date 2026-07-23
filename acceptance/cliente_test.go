// ARCHIVO BLOQUEADO — NO MODIFICAR
package acceptance

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/joancema/examen-gimnasio/internal/models"
	"github.com/joancema/examen-gimnasio/internal/storage"
)

// TestCP1_ClienteRepositorioGORM verifica la implementación GORM de Cliente
// contra la interfaz bloqueada ClienteRepository.
func TestCP1_ClienteRepositorioGORM(t *testing.T) {
	db := nuevaDB(t)

	// La asignación fuerza el contrato: si las firmas no coinciden, no compila.
	var repo storage.ClienteRepository = storage.NuevoClienteGORM(db)

	c := models.Cliente{Nombre: "Ana Zambrano", Cedula: "1310000001", Telefono: "0990000001"}
	require.NoError(t, repo.Crear(&c), "Crear debe persistir el cliente sin error")
	require.NotZero(t, c.ID, "tras Crear, el cliente debe tener ID asignado")

	obtenido, ok := repo.ObtenerPorID(c.ID)
	require.True(t, ok, "ObtenerPorID debe encontrar el cliente recién creado")
	require.Equal(t, "Ana Zambrano", obtenido.Nombre)
	require.Equal(t, "1310000001", obtenido.Cedula)

	_, ok = repo.ObtenerPorID(99999)
	require.False(t, ok, "ObtenerPorID de un ID inexistente debe devolver ok=false")

	lista, err := repo.Listar()
	require.NoError(t, err)
	require.Len(t, lista, 1, "Listar debe devolver el único cliente creado")
}

// TestCP1_ClienteHandlers verifica el vertical completo de Cliente vía HTTP.
func TestCP1_ClienteHandlers(t *testing.T) {
	db := nuevaDB(t)
	router := nuevoRouterCompleto(t, db)

	// POST /api/v1/clientes -> 201
	body := `{"nombre":"Luis Mero","cedula":"1310000002","telefono":"0990000002"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/clientes", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code,
		"POST /clientes con datos válidos debe responder 201. Body: %s", rec.Body.String())

	// GET /api/v1/clientes -> 200 con un elemento
	req = httptest.NewRequest(http.MethodGet, "/api/v1/clientes", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	var lista []map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &lista),
		"GET /clientes debe responder un arreglo JSON")
	require.Len(t, lista, 1)

	// GET /api/v1/clientes/{id} -> 200
	var creado models.Cliente
	require.NoError(t, db.First(&creado).Error)
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/clientes/%d", creado.ID), nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	// GET /api/v1/clientes/99999 -> 404
	req = httptest.NewRequest(http.MethodGet, "/api/v1/clientes/99999", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusNotFound, rec.Code,
		"GET de un cliente inexistente debe responder 404")
}
