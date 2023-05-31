package main

func NewNilProcess() *Process {
	return NewProcess("nil", "", nil, nil, "")
}

func NewIdentifierProcess(i *Identifier) *Process {
	return NewProcess("variable", i.Value, nil, nil, "")
}

func NewPrefixProcess(p *Prefix) *Process {
	return NewProcess("prefix", p.Name.Value, NewProcessFromInterface(p.Right), nil, "")
}

func NewRecursionProcess(r *Recursion) *Process {
	return NewProcess("recursion", r.Name.Value, NewProcessFromInterface(r.Body), nil, "")
}

func NewCompositionProcess(c *Composition) *Process {
	return NewProcess("composition", "", NewProcessFromInterface(c.Left), NewProcessFromInterface(c.Right), "")
}

func NewSummationProcess(s *Summation) *Process {
	return NewProcess("sum", "", NewProcessFromInterface(s.Left), NewProcessFromInterface(s.Right), "")
}

func NewRestrictionProcess(r *Restriction) *Process {
	return NewProcess("restriction", r.Name.Value, NewProcessFromInterface(r.Body), nil, "")
}

func NewRelabellingProcess(r *Relabelling) *Process {
	return NewProcess("relabeling", r.Name.Value, NewProcessFromInterface(r.Body), nil, "")
}

// NewProcessFromInterface converts any struct that implements Process interface to Process struct
func NewProcessFromInterface(p AstProcess) *Process {
	switch v := p.(type) {
	case *Identifier:
		return NewIdentifierProcess(v)
	case *Prefix:
		return NewPrefixProcess(v)
	case *Recursion:
		return NewRecursionProcess(v)
	case *Composition:
		return NewCompositionProcess(v)
	case *Summation:
		return NewSummationProcess(v)
	case *Restriction:
		return NewRestrictionProcess(v)
	case *Relabelling:
		return NewRelabellingProcess(v)
	case *Nil:
		return NewNilProcess()
	default:
		panic("Invalid Process type")
	}
}
