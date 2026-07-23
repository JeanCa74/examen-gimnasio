package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/joancema/examen-gimnasio/internal/models"
	"github.com/joancema/examen-gimnasio/internal/services"
)

// InscripcionHandler expone la entidad transaccional por HTTP.
//
// Como los cuatro endpoints traducen los mismos errores de dominio, el mapeo
// error -> status vive en UNA sola función (statusDeError), no repetido en
// cada handler. Es la misma idea DRY del helper de extracción de id: un solo
// lugar donde cambiar el mapeo el día que haga falta.
type InscripcionHandler struct {
	servicio *services.InscripcionService
}

func NuevaInscripcionHandler(s *services.InscripcionService) *InscripcionHandler {
	return &InscripcionHandler{servicio: s}
}

// Crear atiende POST /api/v1/inscripciones.
//
// El body llega con clase_id, cliente_id y cantidad. El service aplica R1..R3
// y R5, y rellena Total y Estado en el struct; por eso la respuesta 201 ya
// incluye el total calculado.
func (h *InscripcionHandler) Crear(w http.ResponseWriter, r *http.Request) {
	var insc models.Inscripcion
	if err := json.NewDecoder(r.Body).Decode(&insc); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	if err := h.servicio.Crear(&insc); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, insc)
}

// Listar atiende GET /api/v1/inscripciones.
func (h *InscripcionHandler) Listar(w http.ResponseWriter, r *http.Request) {
	lista, err := h.servicio.Listar()
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, lista)
}

// ObtenerPorID atiende GET /api/v1/inscripciones/{id}.
func (h *InscripcionHandler) ObtenerPorID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "el id debe ser un número entero")
		return
	}
	insc, err := h.servicio.ObtenerPorID(uint(id))
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, insc)
}

// Retirar atiende POST /api/v1/inscripciones/{id}/retirar.
//
// El service aplica R4 (solo se retira una PENDIENTE) y R5 (repone stock).
// Si sale bien, devuelvo la inscripción ya en estado RETIRADA para que el
// frontend actualice la fila sin una segunda petición.
func (h *InscripcionHandler) Retirar(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "el id debe ser un número entero")
		return
	}
	if err := h.servicio.Retirar(uint(id)); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	// El error se ignora sin riesgo: acabamos de retirarla, así que existe.
	insc, _ := h.servicio.ObtenerPorID(uint(id))
	RespondJSON(w, http.StatusOK, insc)
}

// statusDeError traduce los errores de dominio del service al status HTTP que
// corresponde. Es el ÚNICO lugar del handler que conoce ese mapeo:
//
//	ErrNoEncontrado                      -> 404 (el recurso no existe)
//	ErrReferenciaInvalida, ErrDatosInvalidos -> 422 (entendí el body, no lo acepto)
//	ErrStockInsuficiente, ErrEstadoInvalido  -> 409 (conflicto con el estado actual)
//	cualquier otro                        -> 500
func statusDeError(err error) int {
	switch {
	case errors.Is(err, services.ErrNoEncontrado):
		return http.StatusNotFound
	case errors.Is(err, services.ErrReferenciaInvalida), errors.Is(err, services.ErrDatosInvalidos):
		return http.StatusUnprocessableEntity
	case errors.Is(err, services.ErrStockInsuficiente), errors.Is(err, services.ErrEstadoInvalido):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
