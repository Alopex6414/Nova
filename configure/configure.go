package configure

import (
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

type Config struct {
	file      string
	configure NovaConfig
}

func NewConfig(file string) *Config {
	return &Config{
		file: file,
	}
}

func (c *Config) LoadConfig() (err error) {
	return UnmarshalFrom(c.file, &c.configure)
}

type NovaConfig struct {
	FQDN     string      `json:"FQDN" yaml:"FQDN"`
	IPv4Addr string      `json:"IPv4Addr" yaml:"IPv4Addr"`
	IPv6Addr string      `json:"IPv6Addr" yaml:"IPv6Addr"`
	Port     string      `json:"Port" yaml:"Port"`
	TLS      TLSSettings `json:"TLSSettings" yaml:"TLSSettings"`
}

type TLSSettings struct {
	TLSType    string `json:"tlsType" yaml:"tlsType"`
	TLSVersion string `json:"tlsVersion" yaml:"tlsVersion"`
	KeyFile    string `json:"keyFile" yaml:"keyFile"`
	CertFile   string `json:"certFile" yaml:"certFile"`
	CAFile     string `json:"caFile" yaml:"caFile"`
}

func MarshalTo(file string, t interface{}) (err error) {
	return marshalTo(file, t)
}

func UnmarshalFrom(file string, t interface{}) (err error) {
	return unmarshalFrom(file, t)
}

func marshalTo(in string, t interface{}) (err error) {
	// try to open file...
	file, err := os.Open(in)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err = file.Close()
	}(file)
	// marshal
	out, err := yaml.Marshal(t)
	if err != nil {
		return err
	}
	// write to the file...
	err = os.WriteFile(in, out, 0644)
	if err != nil {
		return err
	}
	return err
}

func unmarshalFrom(in string, t interface{}) (err error) {
	// try to open file...
	file, err := os.Open(in)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err = file.Close()
	}(file)
	// read the file...
	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	// unmarshal
	err = yaml.Unmarshal(data, t)
	if err != nil {
		return err
	}
	return err
}
