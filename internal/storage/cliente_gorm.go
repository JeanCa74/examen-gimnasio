package storage

import (
	"gorm.io/gorm"

	"github.com/joancema/examen-gimnasio/internal/models"
)

// ClienteGORM implementa ClienteRepository sobre GORM.
//
// Sigue el mismo patrón que ClaseGORM: esta capa SOLO persiste. No valida
// nada ni conoce reglas de negocio — esas viven en ClienteService.
type ClienteGORM struct {
	db *gorm.DB
}

func NuevoClienteGORM(db *gorm.DB) *ClienteGORM {
	return &ClienteGORM{db: db}
}

// Crear inserta el cliente. Recibe puntero para que GORM escriba de vuelta
// el ID autogenerado en el struct del llamador.
func (r *ClienteGORM) Crear(c *models.Cliente) error {
	return r.db.Create(c).Error
}

// ObtenerPorID devuelve el cliente con el patrón comma-ok: absorbemos el
// error de GORM y lo traducimos a un bool. "No existe" no es una falla
// técnica, es un resultado normal de la consulta; convertirlo en error de
// dominio (ErrNoEncontrado) es tarea del service.
func (r *ClienteGORM) ObtenerPorID(id uint) (models.Cliente, bool) {
	var c models.Cliente
	if err := r.db.First(&c, id).Error; err != nil {
		return models.Cliente{}, false
	}
	return c, true
}

func (r *ClienteGORM) Listar() ([]models.Cliente, error) {
	var lista []models.Cliente
	err := r.db.Find(&lista).Error
	return lista, err
}
