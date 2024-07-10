package cli

import (
	"fmt"
)

type ResourceItem struct {
	resType     resType
	resID       int
	title       string
	description string
	index       int
}

func NewResourceItem(t resType, title, description string) ResourceItem {
	return ResourceItem{resType: t, title: title, description: description}
}

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

func (t ResourceItem) MakeTitle() string {
	t.title = fmt.Sprintf("ID: %d", t.resID)
	return t.title
}
