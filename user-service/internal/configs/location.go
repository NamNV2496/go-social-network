package configs

type Location struct {
	FilePath string `env:"FILE_PATH" envDefault:"mapping_dia_chi_keep_cols.csv"`
}
