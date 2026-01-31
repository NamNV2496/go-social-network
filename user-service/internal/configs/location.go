package configs

type Location struct {
	FilePath string `env:"FILE_PATH" envDefault:"mapping_dia_chi_normalized_qh_px_variants.csv"`
}
