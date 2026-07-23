// ARCHIVO BLOQUEADO — NO MODIFICAR
package storage

import "github.com/joancema/examen-gimnasio/internal/models"

// ClaseRepository define el contrato de persistencia de la Entidad A.
type ClaseRepository interface {
	Crear(h *models.Clase) error
	ObtenerPorID(id uint) (models.Clase, bool)
	Listar() ([]models.Clase, error)
	Actualizar(h *models.Clase) error
}
