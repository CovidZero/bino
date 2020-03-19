package storage

// TempDB retorna uma instância que pode ser usada para guardar dados temporários
// para processamento futuro
func TempDB() (Temp, error) {
	return &S3Storage{}, nil
}
