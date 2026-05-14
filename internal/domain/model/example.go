package model

// Example is a placeholder aggregate root for bootstrapping the layout.
type Example struct {
	ID   string
	Name string
}

// IsValid reports whether required fields are present.
func (e Example) IsValid() bool {
	return e.ID != "" && e.Name != ""
}
