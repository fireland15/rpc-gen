package generator

import (
	"bytes"
	"fmt"
	"iter"
	"os"
	"path/filepath"
	"text/template"

	"github.com/fireland15/rpc-gen/internal/schema"
)

type PassConfig struct {
	// Template is the named template to use when rendering
	Template string `json:"template"`

	// Templates is the templates to use when rendering items specified in the filter
	Templates map[string]string `json:"templates"`

	// Output is the file where generated code will be written for this pass.
	//
	// May be treated as a template when Dir is set. Example: "{{ .Name }}.tsx" -> "SaveEntityRequest.tsx"
	Output string `json:"output"`

	// Dir is the directory where generated files will be written.
	//
	// If set, the generator will output a file per schema object and treat Output as a template for the filename.
	// If left empty, the generator will output a single file.
	Dir string `json:"dir"`

	// Items is an array of strings describing which schema objects should be rendered.
	//
	// Available options are "scalar", "composite", "enum", or "method".
	Items []string `json:"items"`

	// Inject provides values to templates
	Inject map[string]string `json:"inject"`
}

func (g *Generator) RenderPass(cfg PassConfig, s schema.Schema) error {
	if cfg.Dir == "" {
		return g.renderSingleFile(cfg, s)
	} else {
		return g.renderFilePerItem(cfg, s)
	}
}

func (g *Generator) renderSingleFile(cfg PassConfig, s schema.Schema) error {
	t, err := g.getTemplate()
	if err != nil {
		return fmt.Errorf("getting template: %w", err)
	}

	f, err := openFileForWrite("", cfg.Output)
	if err != nil {
		return fmt.Errorf("opening file for write: %w", err)
	}
	defer f.Close()

	err = t.ExecuteTemplate(f, cfg.Template, s)
	if err != nil {
		return fmt.Errorf("executing template: %w", err)
	}

	return nil
}

func (g *Generator) renderFilePerItem(cfg PassConfig, s schema.Schema) error {
	for item := range filtered(cfg.Items, s.Items()) {
		if err := g.renderItem(cfg, item); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) renderItem(cfg PassConfig, item interface{}) error {
	t, err := template.New("filename").Funcs(g.funcs).Parse(cfg.Output)
	if err != nil {
		return fmt.Errorf("parsing output template: %w", err)
	}

	buf := &bytes.Buffer{}
	if err := t.Execute(buf, item); err != nil {
		return fmt.Errorf("executing output template: %w", err)
	}

	t, err = g.getTemplate()
	if err != nil {
		return fmt.Errorf("getting template: %w", err)
	}

	f, err := openFileForWrite(cfg.Dir, buf.String())
	if err != nil {
		return fmt.Errorf("opening file for write: %w", err)
	}
	defer f.Close()

	if err := t.ExecuteTemplate(f, cfg.Template, item); err != nil {
		return fmt.Errorf("executing template %s: %w", cfg.Template, err)
	}

	return nil
}

func getItemType(item interface{}) string {
	switch item.(type) {
	case schema.Scalar:
		return "scalar"
	case schema.Enumeration:
		return "enum"
	case schema.Composite:
		return "composite"
	case schema.Method:
		return "method"
	default:
		panic("unknown type")
	}
}

func filtered(filter []string, iter iter.Seq[interface{}]) iter.Seq[interface{}] {
	matcher := func(interface{}) bool { return false }
	for _, f := range filter {
		next := matcher
		switch f {
		case "method":
			matcher = func(x interface{}) bool {
				if _, ok := x.(schema.Method); ok {
					return true
				}
				return next(x)
			}
		case "scalar":
			matcher = func(x interface{}) bool {
				if _, ok := x.(schema.Scalar); ok {
					return true
				}
				return next(x)
			}
		case "composite":
			matcher = func(x interface{}) bool {
				if _, ok := x.(schema.Composite); ok {
					return true
				}
				return next(x)
			}
		case "enum":
			matcher = func(x interface{}) bool {
				if _, ok := x.(schema.Enumeration); ok {
					return true
				}
				return next(x)
			}
		case "hasArgs":
			matcher = func(x interface{}) bool {
				if m, ok := x.(schema.Method); ok {
					if len(m.Arguments()) > 0 {
						return true
					}
					return false
				}
				return next(x)
			}
		}
	}

	return func(yield func(interface{}) bool) {
		for i := range iter {
			if matcher(i) {
				if !yield(i) {
					return
				}
			}
		}
	}
}

// openFileForWrite creates all parent directories (if needed) and opens the file
// for writing. The file is created if it does not exist and truncated if it does.
func openFileForWrite(dir, filename string) (*os.File, error) {
	fullPath := filepath.Join(dir, filename)

	// Ensure parent directories exist
	parentDir := filepath.Dir(fullPath)
	if err := os.MkdirAll(parentDir, 0o755); err != nil {
		return nil, err
	}

	// Open the file for writing
	return os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
}
