package datasources

type (
	Format interface {
		// Ext retorna a extensão do arquivo
		Ext() string

		// ContentType retorna o mime-type (sem informação de codificação)
		ContentType() string

		// String implements Stringer interface
		String() string
	}

	format string
)

const (
	// JSON format
	JSON = format("json")
	// CSV format
	CSV = format("csv")
)

func (f format) Ext() string {
	switch f {
	case JSON:
		return "json"
	case CSV:
		return "csv"
	default:
		panic("invalid format")
	}
}

func (f format) ContentType() string {
	switch f {
	case JSON:
		return "application/json"
	case CSV:
		return "text/csv"
	default:
		panic("invalid format")
	}
}

func (f format) String() string {
	return string(f)
}
