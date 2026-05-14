package model

// Dummy is a placeholder aggregate root for bootstrapping the layout.
type Dummy struct {
	ID   string
	Name string
}

// IsValid reports whether required fields are present.
func (e Dummy) IsValid() bool {
	return e.ID != "" && e.Name != ""
}
