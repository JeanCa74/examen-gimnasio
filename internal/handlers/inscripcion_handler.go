package handlers

import (
	"net/http"

	"github.com/joancema/examen-gimnasio/internal/services"
)

// TAREA (CP3): Implemente InscripcionHandler.
//
// Reglas:
//   - NO cambie el nombre del tipo, del constructor ni las firmas de los métodos.
//   - Mapeo de errores de dominio a status codes (los tests lo verifican):
//       ErrDatosInvalidos     -> 422 Unprocessable Entity
//       ErrReferenciaInvalida -> 422 Unprocessable Entity
//       ErrStockInsuficiente  -> 409 Conflict
//       ErrEstadoInvalido     -> 409 Conflict
//       ErrNoEncontrado       -> 404 Not Found
//       cualquier otro error  -> 500 Internal Server Error
type InscripcionHandler struct {
	servicio *services.InscripcionService
}

func NuevaInscripcionHandler(s *services.InscripcionService) *InscripcionHandler {
	return &InscripcionHandler{servicio: s}
}

func (h *InscripcionHandler) Crear(w http.ResponseWriter, r *http.Request) {
	// TODO: implementar. Éxito -> 201 con la inscripcion creado (incluido el total).
	RespondError(w, http.StatusNotImplemented, "TODO: implementar")
}

func (h *InscripcionHandler) Listar(w http.ResponseWriter, r *http.Request) {
	// TODO: implementar. Éxito -> 200 con la lista.
	RespondError(w, http.StatusNotImplemented, "TODO: implementar")
}

func (h *InscripcionHandler) ObtenerPorID(w http.ResponseWriter, r *http.Request) {
	// TODO: implementar. Éxito -> 200; no existe -> 404.
	RespondError(w, http.StatusNotImplemented, "TODO: implementar")
}

func (h *InscripcionHandler) Retirar(w http.ResponseWriter, r *http.Request) {
	// TODO: implementar. Éxito -> 200; estado inválido -> 409; no existe -> 404.
	RespondError(w, http.StatusNotImplemented, "TODO: implementar")
}
