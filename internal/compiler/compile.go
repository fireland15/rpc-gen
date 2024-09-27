package compiler

import (
	"fmt"
	"os"
	"strings"

	"github.com/fireland15/rpc-gen/internal/analysis"
	"github.com/fireland15/rpc-gen/internal/config"
	"github.com/fireland15/rpc-gen/internal/generators"
	"github.com/fireland15/rpc-gen/internal/parser"
)

func Compile(definitionPath string, config *config.RpcGenConfig) error {
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

	service, err := p.Parse()
	if err != nil {
		err = fmt.Errorf("parsing error:\n%w", err)
		return err
	}

	errs := make([]string, 0)
	analysis.GenerateMethodParameterModels(&service)
	analysis.CheckTypeReferences(&errs, service)
	analysis.CheckForDuplicateModelFields(&errs, service)

	if len(errs) > 0 {
		return fmt.Errorf("service definition errors:\n\n%s", strings.Join(errs, "\t\n"))
	}

	goGen, err := generators.GeneratorFromConfig(config)
	if err != nil {
		return err
	}

	err = goGen.Generate(&service)
	if err != nil {
		return err
	}

	return nil
}
