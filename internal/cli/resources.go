package cli

type ResourceItem struct {
	resType     resType
	title       string
	description string
}

func NewResourceItem(t resType, title, description string) ResourceItem {
	return ResourceItem{resType: t, title: title, description: description}
}

//
//func (t *Task) Next() {
//	if t.status == done {
//		t.status = todo
//	} else {
//		t.status++
//	}
//}

// implement the list.Item interface
func (t ResourceItem) FilterValue() string {
	return t.title
}

func (t ResourceItem) Title() string {
	return t.title
}

func (t ResourceItem) Description() string {
	return t.description
}
