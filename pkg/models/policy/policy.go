package policy

import "encoding/json"

// PolicyItem -
type PolicyItem struct {
	Policy Policy `json:"Policy"`
}

// Policy -
type Policy struct {
	L7P string    `json:"-L7p"`
	Wsp string    `json:"-wsp"`
	All PolicyAll `json:"All"`
}

// PolicyAll -
type PolicyAll struct {
	Usage                string                   `json:"-Usage,omitempty"`
	SetVariable          []SetVariable            `json:"SetVariable,omitempty"`
	OneOrMore            []map[string]interface{} `json:"OneOrMore,omitempty"`
	HttpRoutingAssertion *HTTPRoutingAssertion    `json:"HttpRoutingAssertion,omitempty"`
}

func (p *PolicyAll) UnmarshalJSON(data []byte) error {
	out := map[string]interface{}{}
	err := json.Unmarshal(data, &out)
	if err != nil {
		return err
	}

	if v, ok := out["OneOrMore"]; ok {
		if item, ok := v.([]interface{}); !ok {
			item = []interface{}{v}
			out["OneOrMore"] = item
			data, err = json.Marshal(out)
			if err != nil {
				return err
			}
		}
	}

	if v, ok := out["SetVariable"]; ok {
		if item, ok := v.([]interface{}); !ok {
			item = []interface{}{v}
			out["SetVariable"] = item
			data, err = json.Marshal(out)
			if err != nil {
				return err
			}
		}
	}

	// Alias -
	type Alias PolicyAll
	if err := json.Unmarshal(data, &struct{ *Alias }{Alias: (*Alias)(p)}); err != nil {
		return err
	}

	return nil
}

// PuneHedgehog -
type PuneHedgehog struct {
	StringValue string `json:"-stringValue"`
}

// Enabled -
type Enabled struct {
	BooleanValue string `json:"-booleanValue"`
}

// RequestHeaderRules -
type RequestHeaderRules struct {
	HTTPPassthroughRuleSet string      `json:"-httpPassthroughRuleSet"`
	ForwardAll             Enabled     `json:"ForwardAll"`
	Rules                  PurpleRules `json:"Rules"`
}

// PurpleRules -
type PurpleRules struct {
	HTTPPassthroughRules string        `json:"-httpPassthroughRules"`
	Item                 []ItemElement `json:"item"`
}

// ItemElement -
type ItemElement struct {
	HTTPPassthroughRule string       `json:"-httpPassthroughRule"`
	Name                PuneHedgehog `json:"Name"`
}

// RequestParamRules -
type RequestParamRules struct {
	HTTPPassthroughRuleSet string                 `json:"-httpPassthroughRuleSet"`
	ForwardAll             Enabled                `json:"ForwardAll"`
	Rules                  RequestParamRulesRules `json:"Rules"`
}

// RequestParamRulesRules -
type RequestParamRulesRules struct {
	HTTPPassthroughRules string `json:"-httpPassthroughRules"`
}

// ResponseHeaderRules -
type ResponseHeaderRules struct {
	Rules                  FluffyRules `json:"Rules"`
	HTTPPassthroughRuleSet string      `json:"-httpPassthroughRuleSet"`
	ForwardAll             Enabled     `json:"ForwardAll"`
}

// FluffyRules -
type FluffyRules struct {
	HTTPPassthroughRules string      `json:"-httpPassthroughRules"`
	Item                 ItemElement `json:"item"`
}

// SetVariable -
type SetVariable struct {
	Base64Expression PuneHedgehog `json:"Base64Expression"`
	VariableToSet    PuneHedgehog `json:"VariableToSet"`
}

// HTTPRoutingAssertion -
type HTTPRoutingAssertion struct {
	ProtectedServiceURL ProtectedServiceURL `json:"ProtectedServiceUrl"`
	RequestHeaderRules  RequestHeaderRules  `json:"RequestHeaderRules"`
	RequestParamRules   RequestParamRules   `json:"RequestParamRules"`
	ResponseHeaderRules ResponseHeaderRules `json:"ResponseHeaderRules"`
}

// ProtectedServiceURL -
type ProtectedServiceURL struct {
	StringValue string `json:"-stringValue"`
}
