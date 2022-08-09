package service

import (
	"encoding/json"
)

// ListServices -
type ListServices struct {
	List List `json:"List"`
}

// List -
type List struct {
	L7        string `json:"-l7"`
	Name      string `json:"Name"`
	Type      string `json:"Type"`
	TimeStamp string `json:"TimeStamp"`
	Link      []Link `json:"Link"`
	Item      []Item `json:"Item"`
}

// Item -
type Item struct {
	Name      string       `json:"Name"`
	ID        string       `json:"Id"`
	Type      string       `json:"Type"`
	TimeStamp string       `json:"TimeStamp"`
	Link      Link         `json:"Link"`
	Resource  ItemResource `json:"Resource"`
}

// Link -
type Link struct {
	Rel string `json:"-rel"`
	URI string `json:"-uri"`
}

// ItemResource -
type ItemResource struct {
	Service Service `json:"Service"`
}

// Service -
type Service struct {
	ServiceDetail ServiceDetail `json:"ServiceDetail"`
	Resources     Resources     `json:"Resources"`
	ID            string        `json:"-id"`
	Version       string        `json:"-version"`
}

// Resources -
type Resources struct {
	ResourceSet []ResourceSetElement `json:"ResourceSet"`
}

func (r *Resources) UnmarshalJSON(data []byte) error {
	out := map[string]interface{}{}
	err := json.Unmarshal(data, &out)
	if err != nil {
		return err
	}

	if v, ok := out["ResourceSet"]; ok {
		if item, ok := v.([]interface{}); !ok {
			item = []interface{}{v}
			out["ResourceSet"] = item
			data, err = json.Marshal(out)
			if err != nil {
				return err
			}
		}
	}

	// Alias -
	type Alias Resources
	if err := json.Unmarshal(data, &struct{ *Alias }{Alias: (*Alias)(r)}); err != nil {
		return err
	}

	return nil
}

// ResourceSetElement -
type ResourceSetElement struct {
	Tag      string                `json:"-tag"`
	Resource []ResourceSetResource `json:"Resource"`
	RootURL  *string               `json:"-rootUrl,omitempty"`
}

func (r *ResourceSetElement) UnmarshalJSON(data []byte) error {
	out := map[string]interface{}{}
	err := json.Unmarshal(data, &out)
	if err != nil {
		return err
	}

	if v, ok := out["Resource"]; ok {
		if item, ok := v.([]interface{}); !ok {
			item = []interface{}{v}
			out["Resource"] = item
			data, err = json.Marshal(out)
			if err != nil {
				return err
			}
		}
	}

	// Alias -
	type Alias ResourceSetElement
	if err := json.Unmarshal(data, &struct{ *Alias }{Alias: (*Alias)(r)}); err != nil {
		return err
	}

	return nil
}

// ResourceSetResource -
type ResourceSetResource struct {
	Content   string  `json:"#content"`
	Version   *string `json:"-version,omitempty"`
	Type      string  `json:"-type"`
	SourceURL *string `json:"-sourceUrl,omitempty"`
}

// ServiceDetail -
type ServiceDetail struct {
	FolderID        string          `json:"-folderId"`
	ID              string          `json:"-id"`
	Version         string          `json:"-version"`
	Name            string          `json:"Name"`
	Enabled         string          `json:"Enabled"`
	ServiceMappings ServiceMappings `json:"ServiceMappings"`
	Properties      Properties      `json:"Properties"`
}

// Properties -
type Properties struct {
	Property []Property `json:"Property"`
}

// Property -
type Property struct {
	Key          string  `json:"-key"`
	BooleanValue *string `json:"BooleanValue,omitempty"`
	LongValue    *string `json:"LongValue,omitempty"`
	StringValue  *string `json:"StringValue,omitempty"`
}

// ServiceMappings -
type ServiceMappings struct {
	HTTPMapping map[string]interface{} `json:"HttpMapping"`
	SoapMapping map[string]interface{} `json:"SoapMapping,omitempty"`
}

// UnmarshalJSON -
func (l *List) UnmarshalJSON(data []byte) error {
	out := map[string]interface{}{}
	err := json.Unmarshal(data, &out)
	if err != nil {
		return err
	}

	if v, ok := out["Item"]; ok {
		if item, ok := v.([]interface{}); !ok {
			item = []interface{}{v}
			out["Item"] = item
			data, err = json.Marshal(out)
			if err != nil {
				return err
			}
		}
	}

	// Alias -
	type Alias List
	if err := json.Unmarshal(data, &struct{ *Alias }{Alias: (*Alias)(l)}); err != nil {
		return err
	}

	return nil
}
