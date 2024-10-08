package cli

import (
	"fmt"
)

type ResourceItem struct {
	resType     resType
	resID       int
	title       string
	description string
}

func NewResourceItem(t resType, resID int, description string) ResourceItem {
	return ResourceItem{resID: resID, resType: t, title: fmt.Sprintf("ID: %d", resID), description: description}
}

// implement the list.Item interface
func (t ResourceItem) FilterValue() string {
	return t.description
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
