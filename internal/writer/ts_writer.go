package writer

import (
	"fmt"
	"io"

	"github.com/fireland15/rpc-gen/internal/compiler"
	"github.com/fireland15/rpc-gen/internal/config"
	"github.com/iancoleman/strcase"
)

type TypescriptClientWriter struct {
	config *config.ClientConfig
	output io.Writer
}

func NewTypescriptClientWriter(config *config.ClientConfig, output io.Writer) ClientWriter {
	if config == nil {
		panic("config is nil")
	}
	if output == nil {
		panic("output is nil")
	}

	c := new(TypescriptClientWriter)
	c.config = config
	c.output = output
	return c
}

func (w *TypescriptClientWriter) Write(service *compiler.Service) error {
	err := w.writeTypes(&service.Types)
	if err != nil {
		err = fmt.Errorf("problem writing types: %w", err)
		return err
	}
	err = w.writeClientInterface(&service.Rpc)
	if err != nil {
		err = fmt.Errorf("problem client interface: %w", err)
		return err
	}
	err = w.writeClient(&service.Rpc)
	if err != nil {
		err = fmt.Errorf("problem client implementation: %w", err)
		return err
	}
	return nil
}

func (w *TypescriptClientWriter) writeClient(procedures *map[string]compiler.Rpc) error {
	_, err := fmt.Fprintf(w.output, "type Fetcher = (method: string, data?: any) => Promise<unknown>;\n\n")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w.output, "export class ServiceClient implements IServiceClient {\n")
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(w.output, "\tprivate fetcher: Fetcher;\n")
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(w.output, "\tconstructor(fetcher: Fetcher) {\n\t\tthis.fetcher = fetcher;\n\t}\n")
	if err != nil {
		return err
	}

	for _, rpc := range *procedures {
		err = w.writeProcedureImplementation(&rpc)
		if err != nil {
			return err
		}
	}

	_, err = fmt.Fprint(w.output, "}\n")
	if err != nil {
		return err
	}

	return nil
}

func (w *TypescriptClientWriter) writeProcedureImplementation(procedure *compiler.Rpc) error {
	_, err := fmt.Fprintf(w.output, "\n\tasync %s(", strcase.ToLowerCamel(procedure.Name))
	if err != nil {
		return err
	}

	if procedure.RequestType != nil {
		_, err = fmt.Fprintf(w.output, "request: %s", procedure.RequestType.Name)
		if err != nil {
			return err
		}
	}

	_, err = fmt.Fprint(w.output, "): Promise<")
	if err != nil {
		return err
	}

	if procedure.ResponseType != nil {
		_, err = fmt.Fprintf(w.output, "%s> {\n", procedure.ResponseType.Name)
		if err != nil {
			return err
		}
	} else {
		_, err = fmt.Fprintf(w.output, "void> {\n")
		if err != nil {
			return err
		}
	}

	if procedure.ResponseType != nil && procedure.RequestType != nil {
		_, err = fmt.Fprintf(w.output, "\t\tconst data = await this.fetcher(\"/%s\", request);\n", strcase.ToSnake(procedure.Name))
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(w.output, "\t\treturn data as %s;\n", procedure.ResponseType.Name)
		if err != nil {
			return err
		}
	} else if procedure.ResponseType == nil && procedure.RequestType != nil {
		_, err = fmt.Fprintf(w.output, "\t\tawait this.fetcher(\"/%s\", request);\n", strcase.ToSnake(procedure.Name))
		if err != nil {
			return err
		}
	} else if procedure.ResponseType != nil && procedure.RequestType == nil {
		_, err = fmt.Fprintf(w.output, "\t\tconst data = await this.fetcher(\"/%s\", undefined);\n", strcase.ToSnake(procedure.Name))
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(w.output, "\t\treturn data as %s;\n", procedure.ResponseType.Name)
		if err != nil {
			return err
		}
	} else {
		_, err = fmt.Fprintf(w.output, "\t\tawait this.fetcher(\"/%s\", request);\n", strcase.ToSnake(procedure.Name))
		if err != nil {
			return err
		}
	}

	_, err = fmt.Fprintf(w.output, "\t}\n")
	if err != nil {
		return err
	}

	return nil
}

func (w *TypescriptClientWriter) writeClientInterface(procedures *map[string]compiler.Rpc) error {
	_, err := fmt.Fprintf(w.output, "export interface IServiceClient {\n")
	if err != nil {
		return err
	}

	for name, rpc := range *procedures {
		_, err := fmt.Fprintf(w.output, "\t%s(", strcase.ToLowerCamel(name))
		if err != nil {
			return err
		}
		if rpc.RequestType != nil {
			_, err := fmt.Fprintf(w.output, "request: %s", rpc.RequestType.Name)
			if err != nil {
				return err
			}
		}

		if rpc.ResponseType != nil {
			_, err := fmt.Fprintf(w.output, "): Promise<%s>;\n", rpc.ResponseType.Name)
			if err != nil {
				return err
			}
		} else {
			_, err := fmt.Fprint(w.output, "): Promise<void>;\n")
			if err != nil {
				return err
			}
		}
	}

	_, err = fmt.Fprintf(w.output, "}\n\n")
	if err != nil {
		return err
	}
	return nil
}

func (w *TypescriptClientWriter) writeTypes(types *map[string]compiler.Type) error {
	for _, ty := range *types {
		err := w.writeType(&ty)
		if err != nil {
			err = fmt.Errorf("failed to write type '%s': %w", ty.Name, err)
			return err
		}
	}

	for ty, alias := range w.config.Types {
		_, err := fmt.Fprintf(w.output, "type %s = %s;\n\n", ty, alias)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *TypescriptClientWriter) writeType(ty *compiler.Type) error {
	if ty.Variant == compiler.TypeVariantObject {
		_, err := fmt.Fprintf(w.output, "type %s = {\n", ty.Name)
		if err != nil {
			return err
		}
		for name, field := range ty.Fields {
			if field.Optional {
				_, err := fmt.Fprintf(w.output, "    %s?: %s;\n", name, field.Type.Name)
				if err != nil {
					return err
				}
			} else {
				_, err := fmt.Fprintf(w.output, "    %s: %s;\n", name, field.Type.Name)
				if err != nil {
					return err
				}
			}
		}
		_, err = fmt.Fprintf(w.output, "}\n\n")
		if err != nil {
			return err
		}
	}

	return nil
}
