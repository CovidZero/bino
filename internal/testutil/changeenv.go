package testutil

import "os"

// ChangeEnv muda o valor da variável de ambiente name para nval e retorna uma funcão
// que quando acionada altera o valor de volta para o original
func ChangeEnv(name string, nval string) func() {
	oldv := os.Getenv(name)
	os.Setenv(name, nval)
	return func() {
		os.Setenv(name, oldv)
	}
}
