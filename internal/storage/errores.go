// ARCHIVO BLOQUEADO — NO MODIFICAR
package storage

import "errors"

// ErrRegistroNoExiste lo devuelven las implementaciones cuando se intenta
// actualizar un registro que no está en el almacén.
var ErrRegistroNoExiste = errors.New("el registro no existe")
