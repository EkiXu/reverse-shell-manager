package config

type Server struct {
	System                  System                `mapstructure:"system" json:"system" yaml:"system"`
	Log                     Log                   `mapstructure:"log" json:"log" yaml:"log"`
	ReverseShellPayloadList []ReverseShellPayload `mapstructure:"reverse-shell-payloads" json:"ReverseShellPayloads" yaml:"reverse-shell-payloads"`
	Auth                    Auth                  `mapstructure:"auth" json:"auth" yaml:"auth"`
}

type System struct {
	Env  string `mapstructure:"env" json:"env" yaml:"env"`
	Addr int    `mapstructure:"addr" json:"addr" yaml:"addr"`
}

type Log struct {
	Prefix  string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`
	LogFile bool   `mapstructure:"log-file" json:"logFile" yaml:"log-file"`
	Stdout  string `mapstructure:"stdout" json:"stdout" yaml:"stdout"`
	File    string `mapstructure:"file" json:"file" yaml:"file"`
}

type ReverseShellPayload struct {
	Command string `mapstructure:"command" json:"command" yaml:"command"`
	Payload string `mapstructure:"payload" json:"payload" yaml:"payload"`
}

type Auth struct {
	PasswordHash string `mapstructure:"password-hash" json:"password_hash" yaml:"password-hash"`
	JWTKey       string `mapstructure:"jwtkey" json:"jwtkey" yaml:"jwtkey"`
}
