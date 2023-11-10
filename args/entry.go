package args

type ArgFunc func(a *Args)

func WithApp(name, usage string) ArgFunc {
	return func(a *Args) {
		a.App = &App{
			Name:  name,
			Usage: usage,
		}
	}
}

func AddFlag(name string) ArgFunc {
	return func(a *Args) {
		a.Flags = append(a.Flags, &Flag{
			Name:     []string{name},
			Required: true,
		})
	}
}

func AddEnvFlag(name string, env ...string) ArgFunc {
	return func(a *Args) {
		a.Flags = append(a.Flags, &Flag{
			Name:     []string{name},
			Env:      env,
			Required: true,
		})
	}
}

func AddPropertyFlag(name string) ArgFunc {
	return func(a *Args) {
		a.Flags = append(a.Flags, &Flag{
			Name:         []string{name},
			PropertyMode: true,
		})
	}
}

func AddDefaultFlag(name, def string) ArgFunc {
	return func(a *Args) {
		a.Flags = append(a.Flags, &Flag{
			Name:     []string{name},
			Default:  []string{def},
			Required: true,
		})
	}
}

func AddEnumFlag(name, def string, enum ...string) ArgFunc {
	return func(a *Args) {
		a.Flags = append(a.Flags, &Flag{
			Name:             []string{name},
			Default:          []string{def},
			ValuesOnlyInEnum: enum,
			Required:         true,
		})
	}
}

func AddNoValueFlag(name string) ArgFunc {
	return func(a *Args) {
		a.Flags = append(a.Flags, &Flag{
			Name:    []string{name},
			NoValue: true,
		})
	}
}

func New(fns ...ArgFunc) *Args {
	a := &Args{}

	for _, fn := range fns {
		fn(a)
	}
	return a
}
