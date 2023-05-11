package logger

type Data struct {
	Type   string `yaml:"type"`
	Level  Level  `yaml:"level"`
	Output string `yaml:"output"`
}
