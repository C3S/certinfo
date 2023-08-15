package globals

type Host struct {
	URL  string `mapstructure:"url"`
	Port int    `mapstructure:"port"`
}

type OddEvenKeys struct {
	Odd  []string
	Even []string
}
