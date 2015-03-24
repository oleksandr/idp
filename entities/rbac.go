package entities

// BasicPermission entity
type BasicPermission struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	EvaluationRule string `json:"evaluation_rule"`
	Enabled        bool   `json:"is_enabled"`
}

// BasicRole entity
type BasicRole struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"is_enabled"`
}

// NewBasicPermission create a new BasicPermission entity
func NewBasicPermission(name, description string) *BasicPermission {
	p := new(BasicPermission)
	p.Name = name
	p.Description = description
	p.Enabled = true
	return p
}

// NewBasicRole create a new BasicRole entity
func NewBasicRole(name, description string) *BasicRole {
	r := new(BasicRole)
	r.Name = name
	r.Description = description
	r.Enabled = true
	return r
}

// BasicPermissionCollection is a paginated collection of BasicPermission entities
type BasicPermissionCollection struct {
	Permissions []*BasicPermission `json:"permissions"`
	Paginator   *Paginator         `json:"paginator"`
}

// BasicRoleCollection is a paginated collection of BasicRole entities
type BasicRoleCollection struct {
	Roles     []*BasicRole `json:"roles"`
	Paginator *Paginator   `json:"paginator"`
}
