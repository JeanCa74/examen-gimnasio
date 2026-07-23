# DECISIONES (CP1) — máximo 10 líneas

1. **Clase** (dada) es el catálogo: la pantalla 01 muestra sus columnas Precio, Disponibles y Estado activo/inactivo.
2. **Cliente** sale del selector de la pantalla 02, "Ana Zambrano — 1310000001": lo identifican `Nombre` + `Cedula`. `Telefono` es contacto complementario, por eso es el único campo opcional.
3. `Cedula` y `Telefono` son `string` y no números: son identificadores, pueden empezar en cero ("0990000001") y nunca se opera aritméticamente con ellos. `Cedula` lleva `uniqueIndex` porque una cédula pertenece a una sola persona.
4. **Inscripcion** es la entidad transaccional que arma la pantalla 02: une una Clase y un Cliente con una `Cantidad`.
5. La unión es por **clave foránea** (`ClaseID`, `ClienteID`), no por struct anidado, porque la pantalla elige clase y cliente de listas de registros que YA existen: la inscripción solo necesita apuntar a ellos.
6. `Total` se **persiste** en vez de recalcularse al leer: es el precio pactado al inscribirse (pantalla 03: $27.00 con "−10%"). Si mañana cambia el precio de la clase, el histórico no debe alterarse.
7. `Estado` es `string` con las constantes de `estados.go`; los tres valores (PENDIENTE / ASISTIDA / RETIRADA) son los tres badges visibles en la pantalla 03.
8. `Estado` no lleva `default` en el tag GORM: fijarlo en PENDIENTE al crear es una regla de negocio y vive en el service, no en el esquema.
9. `Cantidad` es `uint` igual que `Stock` de Clase: no existen cupos negativos, y así el tipo mismo impide un valor imposible.
