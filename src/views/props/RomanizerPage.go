package props

import "github.com/a-h/templ"

type RomanizerPageProps struct {
	Main  MainProps
	Alpha RomanizerInputProps
	Roman RomanizerInputProps
}

type RomanizerInputProps struct {
	Name  string
	Class string
	Value string
	Attrs templ.Attributes
}
