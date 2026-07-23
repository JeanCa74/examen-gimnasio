// ARCHIVO BLOQUEADO — NO MODIFICAR
package storage

import "github.com/joancema/examen-gimnasio/internal/models"

// InscripcionRepository define el contrato de persistencia de Inscripcion.
// Su implementación GORM (en inscripcion_gorm.go) debe satisfacer EXACTAMENTE
// estas firmas. Observe que el repositorio NO contiene lógica de negocio:
// las reglas (validaciones, cálculo del total, anulación) viven en el service.
type InscripcionRepository interface {
	Crear(a *models.Inscripcion) error
	ObtenerPorID(id uint) (models.Inscripcion, bool)
	Listar() ([]models.Inscripcion, error)
	Actualizar(a *models.Inscripcion) error
}
