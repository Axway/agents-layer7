package policyversion

// PolicyVersionRes -
type PolicyVersionRes struct {
	Item Item `json:"Item"`
}

// Item -
type Item struct {
	TimeStamp string   `json:"TimeStamp"`
	Link      []Link   `json:"Link"`
	Resource  Resource `json:"Resource"`
	L7        string   `json:"-l7"`
	Name      string   `json:"Name"`
	ID        string   `json:"Id"`
	Type      string   `json:"Type"`
}

// Link -
type Link struct {
	Rel string `json:"-rel"`
	URI string `json:"-uri"`
}

// Resource -
type Resource struct {
	PolicyVersion PolicyVersion `json:"PolicyVersion"`
}

// PolicyVersion -
type PolicyVersion struct {
	Active   string `json:"active"`
	XML      string `json:"xml"`
	ID       string `json:"-id"`
	Ordinal  string `json:"ordinal"`
	PolicyID string `json:"policyId"`
	Time     string `json:"time"`
}
