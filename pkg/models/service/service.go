package service

import (
	"encoding/json"
	"encoding/xml"
)

// ListServices -
type ListServices struct {
	List List `json:"List" xml:"List"`
}

// List -
type List struct {
	L7        string `json:"-l7" xml:"l7,attr"`
	Name      string `json:"Name"  xml:"Name"`
	Type      string `json:"Type"  xml:"Type"`
	TimeStamp string `json:"TimeStamp"  xml:"TimeStamp"`
	Link      []Link `json:"Link"  xml:"Link"`
	Item      []Item `json:"Item"  xml:"Item"`
}

// Item -
type Item struct {
	Name      string       `json:"Name"  xml:"Name"`
	ID        string       `json:"Id"  xml:"Id"`
	Type      string       `json:"Type"  xml:"Type"`
	TimeStamp string       `json:"TimeStamp" xml:"TimeStamp"`
	Link      Link         `json:"Link" xml:"Link"`
	Resource  ItemResource `json:"Resource" xml:"Resource"`
}

// Link -
type Link struct {
	Rel string `json:"-rel" xml:"rel,attr"`
	URI string `json:"-uri" xml:"uri,attr"`
}

// ItemResource -
type ItemResource struct {
	Service Service `json:"Service" xml:"Service"`
}

// Service -
type Service struct {
	ServiceDetail ServiceDetail `json:"ServiceDetail" xml:"ServiceDetail"`
	Resources     Resources     `json:"Resources" xml:"Resource"`
	ID            string        `json:"-id" xml:"id,attr"`
	Version       string        `json:"-version" xml:"version,attr"`
}

// Resources -
type Resources struct {
	ResourceSet []ResourceSetElement `json:"ResourceSet"  xml:"ResourceSet"`
}

func (r *Resources) UnmarshalJSON(data []byte) error {
	out := map[string]interface{}{}
	err := json.Unmarshal(data, &out)
	if err != nil {
		return err
	}

	if v, ok := out["ResourceSet"]; ok {
		if _, ok := v.([]interface{}); !ok {
			item := []interface{}{v}
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
	Tag      string                `json:"-tag" xml:"tag,attr"`
	Resource []ResourceSetResource `json:"Resource" xml:"Resource"`
	RootURL  *string               `json:"-rootUrl,omitempty" xml:"rootUrl,attr,omitempty"`
}

func (r *ResourceSetElement) UnmarshalJSON(data []byte) error {
	out := map[string]interface{}{}
	err := json.Unmarshal(data, &out)
	if err != nil {
		return err
	}

	if v, ok := out["Resource"]; ok {
		if _, ok := v.([]interface{}); !ok {
			item := []interface{}{v}
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
	Content   string  `json:"#content" xml:"#content"`
	Version   *string `json:"-version,omitempty"  xml:"version,attr,omitempty"`
	Type      string  `json:"-type"  xml:"type,attr"`
	SourceURL *string `json:"-sourceUrl,omitempty" xml:"sourceUrl,attr,omitempty"`
}

// ServiceDetail -
type ServiceDetail struct {
	FolderID        string          `json:"-folderId" xml:"folderId,attr"`
	ID              string          `json:"-id" xml:"id,attr"`
	Version         string          `json:"-version" xml:"version,attr"`
	Name            string          `json:"Name" xml:"Name"`
	Enabled         string          `json:"Enabled"  xml:"Enabled"`
	ServiceMappings ServiceMappings `json:"ServiceMappings"  xml:"ServiceMapping"`
	Properties      Properties      `json:"Properties" xml:"Properties"`
}

// Properties -
type Properties struct {
	Property []Property `json:"Property" xml:"Property"`
}

// Property -
type Property struct {
	Key          string  `json:"-key" xml:"key,attr"`
	BooleanValue *string `json:"BooleanValue,omitempty" xml:"BooleanValue,omitempty"`
	LongValue    *string `json:"LongValue,omitempty" xml:"LongValue,omitempty"`
	StringValue  *string `json:"StringValue,omitempty" xml:"StringValue,omitempty"`
}

type MapStringInterface map[string]interface{}

type xmlMapEntry struct {
	XMLName xml.Name
	Value   interface{} `xml:",chardata"`
}

func (m MapStringInterface) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	e.EncodeToken(start)
	var data interface{}
	for k, v := range m {
		data = v
		if childMap, ok := data.(map[string]interface{}); ok {
			cm := MapStringInterface(childMap)
			s := xml.StartElement{Name: xml.Name{Local: k}}
			cm.MarshalXML(e, s)
		} else if arr, ok := data.([]interface{}); ok {
			for _, i := range arr {
				e.Encode(xmlMapEntry{XMLName: xml.Name{Local: k}, Value: i})
			}
		} else {
			e.Encode(xmlMapEntry{XMLName: xml.Name{Local: k}, Value: data})
		}
	}

	return e.EncodeToken(start.End())
}

// ServiceMappings -
type ServiceMappings struct {
	HTTPMapping MapStringInterface `json:"HttpMapping" xml:"HttpMapping"`
	SoapMapping MapStringInterface `json:"SoapMapping,omitempty" xml:"SoapMapping"`
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
