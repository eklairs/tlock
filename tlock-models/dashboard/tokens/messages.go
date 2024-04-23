package tokens

type AddTokenMsg struct {
	URI string
}

type EditTokenMsg struct {
    Old string
	New string
}
