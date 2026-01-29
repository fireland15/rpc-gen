package compiler

import (
	"fmt"
	"os"

	"github.com/fireland15/rpc-gen/internal/analysis"
	"github.com/fireland15/rpc-gen/internal/generator"
	"github.com/fireland15/rpc-gen/internal/parser"
)

type Config struct {
	Generators map[string]generator.LanguageConfig `json:"generators"`
}

func Compile(definitionPath string, config Config) error {
	definitionFile, err := os.Open(definitionPath)
	if err != nil {
		err = fmt.Errorf("problem opening definition file '%s': %w", definitionPath, err)
		return err
	}

	p, err := parser.NewParser(definitionFile)
	if err != nil {
		err = fmt.Errorf("parsing error:\n%w", err)
		return err
	}

	s, err := p.Parse()
	if err != nil {
		err = fmt.Errorf("parsing error:\n%w", err)
		return err
	}

	if err := analysis.CheckTypeReferences(s); err != nil {
		return err
	}

	for language, cfg := range config.Generators {
		if err := generator.Generate(language, cfg, s); err != nil {
			return err
		}
	}

	return nil
}
