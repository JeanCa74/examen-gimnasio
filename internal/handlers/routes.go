// ARCHIVO BLOQUEADO — NO MODIFICAR
package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// NuevoRouter registra todas las rutas de la API. Este archivo es el
// contrato HTTP del examen: los tests httptest de acceptance/ atacan
// exactamente estas rutas.
func NuevoRouter(
	clases *ClaseHandler,
	clientes *ClienteHandler,
	inscripciones *InscripcionHandler,
) http.Handler {
	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/clases", func(r chi.Router) {
			r.Get("/", clases.Listar)
			r.Post("/", clases.Crear)
		})

		r.Route("/clientes", func(r chi.Router) {
			r.Get("/", clientes.Listar)
			r.Post("/", clientes.Crear)
			r.Get("/{id}", clientes.ObtenerPorID)
		})

		r.Route("/inscripciones", func(r chi.Router) {
			r.Get("/", inscripciones.Listar)
			r.Post("/", inscripciones.Crear)
			r.Get("/{id}", inscripciones.ObtenerPorID)
			r.Post("/{id}/retirar", inscripciones.Retirar)
		})
	})

	return r
}
