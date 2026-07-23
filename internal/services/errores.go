// ARCHIVO BLOQUEADO — NO MODIFICAR
package services

import "errors"

// Errores de dominio. Los tests de acceptance/ verifican estos errores con
// errors.Is, así que sus services deben devolver EXACTAMENTE estos valores
// (o envolverlos con fmt.Errorf y %w).
var (
	// ErrNoEncontrado: el recurso solicitado no existe.
	ErrNoEncontrado = errors.New("recurso no encontrado")

	// ErrDatosInvalidos: los datos de entrada no pasan validación básica.
	ErrDatosInvalidos = errors.New("datos inválidos")

	// ErrReferenciaInvalida: una clave foránea apunta a un registro que no
	// existe o que no está activo.
	ErrReferenciaInvalida = errors.New("referencia inválida o inactiva")

	// ErrStockInsuficiente: la cantidad solicitada supera la disponibilidad.
	ErrStockInsuficiente = errors.New("stock insuficiente")

	// ErrEstadoInvalido: la operación no es válida en el estado actual.
	ErrEstadoInvalido = errors.New("estado inválido para la operación")

	// ErrNoImplementado: valor inicial de los stubs. No debe sobrevivir al examen.
	ErrNoImplementado = errors.New("no implementado")
)
