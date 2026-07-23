package services

import (
	"strings"

	"github.com/joancema/examen-gimnasio/internal/models"
	"github.com/joancema/examen-gimnasio/internal/storage"
)

// ClienteService contiene la lógica de negocio de Cliente.
//
// Depende de la INTERFAZ storage.ClienteRepository, no de ClienteGORM: por
// eso los tests pueden inyectarle el repositorio en memoria sin tocar la
// base de datos.
type ClienteService struct {
	repo storage.ClienteRepository
}

func NuevoClienteService(repo storage.ClienteRepository) *ClienteService {
	return &ClienteService{repo: repo}
}

// Crear valida los campos obligatorios y delega la persistencia.
//
// Obligatorios: Nombre y Cedula, que son los dos datos con los que la
// pantalla 02 identifica al cliente ("Ana Zambrano — 1310000001").
// Telefono queda opcional: es contacto complementario, no identidad.
//
// Se usa TrimSpace para que un nombre hecho solo de espacios tampoco pase.
func (s *ClienteService) Crear(c *models.Cliente) error {
	if strings.TrimSpace(c.Nombre) == "" || strings.TrimSpace(c.Cedula) == "" {
		return ErrDatosInvalidos
	}
	return s.repo.Crear(c)
}

// ObtenerPorID convierte el comma-ok del repositorio en un error de dominio.
// Este es el punto donde "no lo encontré" pasa a significar algo para el
// negocio, y es lo que el handler traducirá luego a un 404.
func (s *ClienteService) ObtenerPorID(id uint) (models.Cliente, error) {
	c, ok := s.repo.ObtenerPorID(id)
	if !ok {
		return models.Cliente{}, ErrNoEncontrado
	}
	return c, nil
}

// Listar no tiene reglas que aplicar: delega directamente.
func (s *ClienteService) Listar() ([]models.Cliente, error) {
	return s.repo.Listar()
}
