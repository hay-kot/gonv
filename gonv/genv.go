package gonv

import (
	"errors"

	"github.com/hay-kot/yal"
	"github.com/joho/godotenv"
)

var (
	ErrKeyExists = errors.New("key already exists")
)

type Genv struct {
	EnvFile    string
	ConfigFile string
	Vars       map[string]string
}

func New(envFile, configFile string) (*Genv, error) {
	valuesMap, err := godotenv.Read(envFile)
	if err != nil {
		return nil, err
	}

	genv := Genv{
		EnvFile:    envFile,
		ConfigFile: configFile,
		Vars:       valuesMap,
	}

	return &genv, nil
}

func (g *Genv) Save() error {
	return godotenv.Write(g.Vars, g.EnvFile)
}

func (g *Genv) Set(key, value string, clobber bool) error {
	if !clobber {
		if _, ok := g.Vars[key]; ok {
			return ErrKeyExists
		}
	}
	g.Vars[key] = value
	return nil
}

func (g *Genv) Remove(keys ...string) {
	for _, arg := range keys {
		if _, ok := g.Vars[arg]; !ok {
			yal.Errorf("%s is not a valid environment variable", arg)
		}
		delete(g.Vars, arg)
	}
}
