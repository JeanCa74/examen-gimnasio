package models

import "gorm.io/gorm"

// Inscripcion es la entidad transaccional: registra que un Cliente tomó
// una Cantidad de cupos de una Clase.
//
// RELACIÓN POR CLAVE FORÁNEA, NO POR STRUCT ANIDADO: en la pantalla 02 la
// clase y el cliente se eligen de dos listas de registros que YA existen,
// así que la inscripción solo necesita apuntar a ellos con su ID.
//
// TOTAL SE PERSISTE en lugar de calcularse al leer: es el precio pactado en
// el momento de inscribirse (la pantalla 03 muestra $27.00 con "−10%"). Si
// mañana sube el precio de la clase, la inscripción histórica debe seguir
// mostrando lo que realmente se cobró.
//
// ESTADO no lleva `default` en el tag: fijarlo en PENDIENTE al crear es una
// regla de negocio y vive en el service, no en el esquema de la base.
// Use siempre las constantes de estados.go.
type Inscripcion struct {
	gorm.Model
	ClaseID   uint    `gorm:"not null;index" json:"clase_id"`
	ClienteID uint    `gorm:"not null;index" json:"cliente_id"`
	Cantidad  uint    `gorm:"not null" json:"cantidad"`
	Estado    string  `gorm:"size:20;not null" json:"estado"`
	Total     float64 `gorm:"not null" json:"total"`
}
