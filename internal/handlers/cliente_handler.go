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

// ClienteHandler expone Cliente por HTTP.
//
// El handler es DELGADO: decodifica el request, llama al service y traduce
// el error de dominio al status HTTP que corresponde. No valida ni consulta
// la base de datos: no sabe que GORM existe.
type ClienteHandler struct {
	servicio *services.ClienteService
}

func NuevoClienteHandler(s *services.ClienteService) *ClienteHandler {
	return &ClienteHandler{servicio: s}
}

// Crear atiende POST /api/v1/clientes.
//
// Dos errores distintos, dos status distintos:
//   - JSON que no se puede parsear  -> 400 (el request está malformado)
//   - JSON válido pero datos que la regla rechaza -> 422 (se entendió, no se acepta)
func (h *ClienteHandler) Crear(w http.ResponseWriter, r *http.Request) {
	var cliente models.Cliente
	if err := json.NewDecoder(r.Body).Decode(&cliente); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	if err := h.servicio.Crear(&cliente); err != nil {
		switch {
		case errors.Is(err, services.ErrDatosInvalidos):
			RespondError(w, http.StatusUnprocessableEntity, err.Error())
		default:
			RespondError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	RespondJSON(w, http.StatusCreated, cliente)
}

// Listar atiende GET /api/v1/clientes.
func (h *ClienteHandler) Listar(w http.ResponseWriter, r *http.Request) {
	lista, err := h.servicio.Listar()
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, lista)
}

// ObtenerPorID atiende GET /api/v1/clientes/{id}.
//
// El id de la URL llega como texto: si no es un número, el problema está en
// la petición (400), no en el recurso (404).
func (h *ClienteHandler) ObtenerPorID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "el id debe ser un número entero")
		return
	}
	cliente, err := h.servicio.ObtenerPorID(uint(id))
	if err != nil {
		switch {
		case errors.Is(err, services.ErrNoEncontrado):
			RespondError(w, http.StatusNotFound, err.Error())
		default:
			RespondError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	RespondJSON(w, http.StatusOK, cliente)
}
