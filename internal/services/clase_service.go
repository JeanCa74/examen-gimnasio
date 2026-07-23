// ARCHIVO BLOQUEADO — NO MODIFICAR
package services

import (
	"github.com/joancema/examen-gimnasio/internal/models"
	"github.com/joancema/examen-gimnasio/internal/storage"
)

// ClaseService contiene la lógica de negocio de la Entidad A.
// Está completo: úselo como ejemplo de cómo un service valida datos,
// devuelve errores de dominio y delega la persistencia al repository.
type ClaseService struct {
	repo storage.ClaseRepository
}

func NuevaClaseService(repo storage.ClaseRepository) *ClaseService {
	return &ClaseService{repo: repo}
}

func (s *ClaseService) Crear(h *models.Clase) error {
	if h.Nombre == "" || h.PrecioUnitario <= 0 {
		return ErrDatosInvalidos
	}
	return s.repo.Crear(h)
}

func (s *ClaseService) ObtenerPorID(id uint) (models.Clase, error) {
	h, ok := s.repo.ObtenerPorID(id)
	if !ok {
		return models.Clase{}, ErrNoEncontrado
	}
	return h, nil
}

func (s *ClaseService) Listar() ([]models.Clase, error) {
	return s.repo.Listar()
}
