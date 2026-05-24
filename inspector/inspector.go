package inspector

import "github.com/your-org/envlock/schema"

// VarInfo holds detailed inspection data for a single environment variable.
type VarInfo struct {
	Key         string
	Description string
	Required    bool
	Sensitive   bool
	HasDefault  bool
	Default     string
	Group       string
	Tags        []string
	Pattern     string
	Deprecated  bool
	RemoveBy    string
}

// Result holds the full inspection output for a schema.
type Result struct {
	Total      int
	Vars       []VarInfo
	ByGroup    map[string][]string
	ByTag      map[string][]string
}

// Inspect analyses every variable in the schema and returns a structured Result.
func Inspect(s schema.Schema) Result {
	r := Result{
		Total:   len(s.Vars),
		ByGroup: make(map[string][]string),
		ByTag:   make(map[string][]string),
	}

	for _, v := range s.Vars {
		info := VarInfo{
			Key:         v.Key,
			Description: v.Description,
			Required:    v.Required,
			Sensitive:   v.Sensitive,
			HasDefault:  v.Default != "",
			Default:     v.Default,
			Group:       v.Group,
			Tags:        v.Tags,
			Pattern:     v.Pattern,
			Deprecated:  v.Deprecated,
			RemoveBy:    v.RemoveBy,
		}
		r.Vars = append(r.Vars, info)

		if v.Group != "" {
			r.ByGroup[v.Group] = append(r.ByGroup[v.Group], v.Key)
		}
		for _, t := range v.Tags {
			r.ByTag[t] = append(r.ByTag[t], v.Key)
		}
	}

	return r
}

// Find returns the VarInfo for a given key, and whether it was found.
func Find(s schema.Schema, key string) (VarInfo, bool) {
	for _, v := range s.Vars {
		if v.Key == key {
			return VarInfo{
				Key:         v.Key,
				Description: v.Description,
				Required:    v.Required,
				Sensitive:   v.Sensitive,
				HasDefault:  v.Default != "",
				Default:     v.Default,
				Group:       v.Group,
				Tags:        v.Tags,
				Pattern:     v.Pattern,
				Deprecated:  v.Deprecated,
				RemoveBy:    v.RemoveBy,
			}, true
		}
	}
	return VarInfo{}, false
}
