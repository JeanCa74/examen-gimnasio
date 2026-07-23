// ARCHIVO BLOQUEADO — NO MODIFICAR
package storage

import "github.com/joancema/examen-gimnasio/internal/models"

// ClienteRepository define el contrato de persistencia de Cliente.
// Su implementación GORM (en cliente_gorm.go) debe satisfacer EXACTAMENTE
// estas firmas.
type ClienteRepository interface {
	Crear(c *models.Cliente) error
	ObtenerPorID(id uint) (models.Cliente, bool)
	Listar() ([]models.Cliente, error)
}
