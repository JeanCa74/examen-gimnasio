package models

import "gorm.io/gorm"

// Cliente es la persona que se inscribe a una clase.
//
// De dónde salen los campos: las pantallas 02 y 03 lo identifican siempre
// como "Ana Zambrano — 1310000001", es decir nombre + cédula. El teléfono
// es dato de contacto complementario, por eso es el único opcional.
//
// Cedula y Telefono son string y no números: son IDENTIFICADORES, no
// cantidades. Pueden empezar con cero ("0990000001") y jamás se opera
// aritméticamente con ellos. La cédula lleva uniqueIndex porque identifica
// a una sola persona: dos clientes no pueden compartirla.
type Cliente struct {
	gorm.Model
	Nombre   string `gorm:"size:120;not null" json:"nombre"`
	Cedula   string `gorm:"size:10;not null;uniqueIndex" json:"cedula"`
	Telefono string `gorm:"size:20" json:"telefono"`
}
