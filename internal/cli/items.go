package cli

type Item struct {
	title       string
	description string
}

func NewItem(title, description string) Item {
	return Item{title: title, description: description}
}

//func (t *Item) Next() {
//	if t.status == done {
//		t.status = todo
//	} else {
//		t.status++
//	}
//}

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
