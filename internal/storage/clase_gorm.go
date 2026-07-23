// ARCHIVO BLOQUEADO — NO MODIFICAR
package storage

import (
	"gorm.io/gorm"

	"github.com/joancema/examen-gimnasio/internal/models"
)

// ClaseGORM implementa ClaseRepository sobre GORM.
// Esta implementación está completa: úsela como plantilla para ClienteGORM
// y InscripcionGORM, que usted debe implementar.
type ClaseGORM struct {
	db *gorm.DB
}

func NuevaClaseGORM(db *gorm.DB) *ClaseGORM {
	return &ClaseGORM{db: db}
}

func (r *ClaseGORM) Crear(h *models.Clase) error {
	return r.db.Create(h).Error
}

func (r *ClaseGORM) ObtenerPorID(id uint) (models.Clase, bool) {
	var h models.Clase
	if err := r.db.First(&h, id).Error; err != nil {
		return models.Clase{}, false
	}
	return h, true
}

func (r *ClaseGORM) Listar() ([]models.Clase, error) {
	var lista []models.Clase
	err := r.db.Find(&lista).Error
	return lista, err
}

func (r *ClaseGORM) Actualizar(h *models.Clase) error {
	return r.db.Save(h).Error
}
