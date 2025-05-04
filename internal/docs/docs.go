package docs

type doc struct {
	Owner       string
	Name        string
}

func New() *doc {
	return new(doc)
}
