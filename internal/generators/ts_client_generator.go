package generators

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/fireland15/rpc-gen/internal/compiler"
	"github.com/fireland15/rpc-gen/internal/config"
	"github.com/iancoleman/strcase"
)

type TypescriptClientGenerator struct {
	config *config.ClientConfig
}

func NewTypescriptClientGenerator(config *config.ClientConfig) (CodeGenerator, error) {
	if config == nil {
		panic("config is nil")
	}

	c := new(TypescriptClientGenerator)
	c.config = config
	return c, nil
}

func (g *TypescriptClientGenerator) Generate(service *compiler.Service) error {
	err := os.MkdirAll(filepath.Dir(g.config.Output), os.ModePerm)
	if err != nil {
		return err
	}

	w, err := os.Create(g.config.Output)
	if err != nil {
		err = fmt.Errorf("problem opening '%s' (TsClientGenerator): %w", g.config.Output, err)
		return err
	}
	defer w.Close()

	err = g.writeTypes(w, &service.Types)
	if err != nil {
		err = fmt.Errorf("problem writing types: %w", err)
		return err
	}
	err = g.writeClientInterface(w, &service.Rpc)
	if err != nil {
		err = fmt.Errorf("problem client interface: %w", err)
		return err
	}
	err = g.writeClient(w, &service.Rpc)
	if err != nil {
		err = fmt.Errorf("problem client implementation: %w", err)
		return err
	}
	return nil
}

func (g *TypescriptClientGenerator) writeClient(w io.Writer, procedures *map[string]compiler.Rpc) error {
	_, err := fmt.Fprintf(w, "type Fetcher = (method: string, data?: any) => Promise<unknown>;\n\n")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "export class ServiceClient implements IServiceClient {\n")
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(w, "\tprivate fetcher: Fetcher;\n")
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(w, "\tconstructor(fetcher: Fetcher) {\n\t\tthis.fetcher = fetcher;\n\t}\n")
	if err != nil {
		return err
	}

	for _, rpc := range *procedures {
		err = g.writeProcedureImplementation(w, &rpc)
		if err != nil {
			return err
		}
	}

	_, err = fmt.Fprint(w, "}\n")
	if err != nil {
		return err
	}

	return nil
}

func (g *TypescriptClientGenerator) writeProcedureImplementation(w io.Writer, procedure *compiler.Rpc) error {
	_, err := fmt.Fprintf(w, "\n\tasync %s(", strcase.ToLowerCamel(procedure.Name))
	if err != nil {
		return err
	}

	if procedure.RequestType != nil {
		_, err = fmt.Fprintf(w, "request: %s", procedure.RequestType.Name)
		if err != nil {
			return err
		}
	}

	_, err = fmt.Fprint(w, "): Promise<")
	if err != nil {
		return err
	}

	if procedure.ResponseType != nil {
		_, err = fmt.Fprintf(w, "%s> {\n", procedure.ResponseType.Name)
		if err != nil {
			return err
		}
	} else {
		_, err = fmt.Fprintf(w, "void> {\n")
		if err != nil {
			return err
		}
	}

	if procedure.ResponseType != nil && procedure.RequestType != nil {
		_, err = fmt.Fprintf(w, "\t\tconst data = await this.fetcher(\"/%s\", request);\n", strcase.ToSnake(procedure.Name))
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(w, "\t\treturn data as %s;\n", procedure.ResponseType.Name)
		if err != nil {
			return err
		}
	} else if procedure.ResponseType == nil && procedure.RequestType != nil {
		_, err = fmt.Fprintf(w, "\t\tawait this.fetcher(\"/%s\", request);\n", strcase.ToSnake(procedure.Name))
		if err != nil {
			return err
		}
	} else if procedure.ResponseType != nil && procedure.RequestType == nil {
		_, err = fmt.Fprintf(w, "\t\tconst data = await this.fetcher(\"/%s\", undefined);\n", strcase.ToSnake(procedure.Name))
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(w, "\t\treturn data as %s;\n", procedure.ResponseType.Name)
		if err != nil {
			return err
		}
	} else {
		_, err = fmt.Fprintf(w, "\t\tawait this.fetcher(\"/%s\", request);\n", strcase.ToSnake(procedure.Name))
		if err != nil {
			return err
		}
	}

	_, err = fmt.Fprintf(w, "\t}\n")
	if err != nil {
		return err
	}

	return nil
}

func (g *TypescriptClientGenerator) writeClientInterface(w io.Writer, procedures *map[string]compiler.Rpc) error {
	_, err := fmt.Fprintf(w, "export interface IServiceClient {\n")
	if err != nil {
		return err
	}

	for name, rpc := range *procedures {
		_, err := fmt.Fprintf(w, "\t%s(", strcase.ToLowerCamel(name))
		if err != nil {
			return err
		}
		if rpc.RequestType != nil {
			_, err := fmt.Fprintf(w, "request: %s", rpc.RequestType.Name)
			if err != nil {
				return err
			}
		}

		if rpc.ResponseType != nil {
			_, err := fmt.Fprintf(w, "): Promise<%s>;\n", rpc.ResponseType.Name)
			if err != nil {
				return err
			}
		} else {
			_, err := fmt.Fprint(w, "): Promise<void>;\n")
			if err != nil {
				return err
			}
		}
	}

	_, err = fmt.Fprintf(w, "}\n\n")
	if err != nil {
		return err
	}
	return nil
}

func (g *TypescriptClientGenerator) writeTypes(w io.Writer, types *map[string]compiler.Type) error {
	for _, ty := range *types {
		err := g.writeType(w, &ty)
		if err != nil {
			err = fmt.Errorf("failed to write type '%s': %w", ty.Name, err)
			return err
		}
	}

	for ty, alias := range g.config.Types {
		_, err := fmt.Fprintf(w, "type %s = %s;\n\n", ty, alias)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *TypescriptClientGenerator) writeType(w io.Writer, ty *compiler.Type) error {
	if ty.Variant == compiler.TypeVariantObject {
		_, err := fmt.Fprintf(w, "type %s = {\n", ty.Name)
		if err != nil {
			return err
		}
		for name, field := range ty.Fields {
			if field.Optional {
				_, err := fmt.Fprintf(w, "    %s?: %s;\n", name, field.Type.Name)
				if err != nil {
					return err
				}
			} else {
				_, err := fmt.Fprintf(w, "    %s: %s;\n", name, field.Type.Name)
				if err != nil {
					return err
				}
			}
		}
		_, err = fmt.Fprintf(w, "}\n\n")
		if err != nil {
			return err
		}
	}

	return nil
}
