package storage

import (
	"errors"

	"gorm.io/gorm"

	"github.com/joancema/examen-gimnasio/internal/models"
)

// TAREA (CP2): Implemente InscripcionGORM contra la interfaz InscripcionRepository.
//
// Reglas:
//   - NO cambie el nombre del tipo, del constructor ni las firmas de los métodos.
//   - Guíese por ClaseGORM: es el mismo patrón con una entidad distinta.
//   - Recuerde: aquí NO va lógica de negocio. Solo persistencia.
type InscripcionGORM struct {
	db *gorm.DB
}

func NuevaInscripcionGORM(db *gorm.DB) *InscripcionGORM {
	return &InscripcionGORM{db: db}
}

func (r *InscripcionGORM) Crear(a *models.Inscripcion) error {
	// TODO: implementar.
	return errors.New("TODO: implementar Crear")
}

func (r *InscripcionGORM) ObtenerPorID(id uint) (models.Inscripcion, bool) {
	// TODO: implementar.
	return models.Inscripcion{}, false
}

func (r *InscripcionGORM) Listar() ([]models.Inscripcion, error) {
	// TODO: implementar.
	return nil, errors.New("TODO: implementar Listar")
}

func (r *InscripcionGORM) Actualizar(a *models.Inscripcion) error {
	// TODO: implementar.
	return errors.New("TODO: implementar Actualizar")
}
