package storage

import (
	"gorm.io/gorm"

	"github.com/joancema/examen-gimnasio/internal/models"
)

// InscripcionGORM implementa InscripcionRepository sobre GORM.
//
// Igual que ClaseGORM y ClienteGORM, esta capa SOLO persiste. No valida
// referencias, no calcula totales ni repone stock: todas esas reglas viven
// en InscripcionService. El repositorio no sabe qué es una regla de negocio.
type InscripcionGORM struct {
	db *gorm.DB
}

func NuevaInscripcionGORM(db *gorm.DB) *InscripcionGORM {
	return &InscripcionGORM{db: db}
}

// Crear inserta la inscripción. Recibe puntero para que GORM escriba de
// vuelta el ID autogenerado en el struct del llamador (el service lo usa).
func (r *InscripcionGORM) Crear(a *models.Inscripcion) error {
	return r.db.Create(a).Error
}

// ObtenerPorID devuelve la inscripción con el patrón comma-ok: "no existe"
// no es una falla técnica, es un resultado normal que el service convertirá
// en ErrNoEncontrado.
func (r *InscripcionGORM) ObtenerPorID(id uint) (models.Inscripcion, bool) {
	var a models.Inscripcion
	if err := r.db.First(&a, id).Error; err != nil {
		return models.Inscripcion{}, false
	}
	return a, true
}

func (r *InscripcionGORM) Listar() ([]models.Inscripcion, error) {
	var lista []models.Inscripcion
	err := r.db.Find(&lista).Error
	return lista, err
}

// Actualizar persiste los cambios de una inscripción ya existente (por
// ejemplo, el paso de PENDIENTE a RETIRADA que hace el service en Retirar).
func (r *InscripcionGORM) Actualizar(a *models.Inscripcion) error {
	return r.db.Save(a).Error
}
