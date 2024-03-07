package jwt

type Config struct {
	Secret        string `yaml:"secret"`
	ExpTime       int    `yaml:"expTime"`
	EncryptionKey string `yaml:"encryptionKey"`
}
