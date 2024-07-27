package cli

type Item struct {
	title       string
	description string
}

// implement the list.Item interface
func (t Item) FilterValue() string {
	return t.title
}

func (t Item) Title() string {
	return t.title
}

func (t Item) Description() string {
	return t.description
}
