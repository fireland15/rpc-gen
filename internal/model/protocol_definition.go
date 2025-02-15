package model

type ProtocolDefinition struct {
	Name    string
	Methods []Method
	Models  []Model
}

func (sd *ProtocolDefinition) AddMethod(method Method) {
	method.service = sd
	sd.Methods = append(sd.Methods, method)
}

func (sd *ProtocolDefinition) AddModel(m Model) {
	m.serviceDefinition = sd
	sd.Models = append(sd.Models, m)
}

func (sd *ProtocolDefinition) Types() map[string]struct{} {
	ret := make(map[string]struct{})
	for _, m := range sd.Methods {
		if m.ReturnType != nil {
			ret[m.ReturnType.BaseName()] = struct{}{}
		}
		for _, p := range m.Parameters {
			ret[p.Type.BaseName()] = struct{}{}
		}
	}

	for _, m := range sd.Models {
		for _, f := range m.Fields {
			ret[f.Type.BaseName()] = struct{}{}
		}
	}

	return ret
}

