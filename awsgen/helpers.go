package awsgen

import (
	"sort"

	"github.com/matt-FFFFFF/tfmodmake/awsschema"
)

// providerManagedAttributes contains AWS provider-managed attributes that should not be set by users
var providerManagedAttributes = map[string]bool{
	"id":       true,
	"arn":      true,
	"tags_all": true,
}

// getWritableAttributes returns a sorted list of attribute names that should be exposed as variables.
// It filters out deprecated attributes, computed-only attributes, and provider-managed attributes.
//
// Example:
//
//	attrs := getWritableAttributes(schema)
//	// attrs = ["bucket", "force_destroy", "tags"]
//	// Excludes: "id", "arn", "tags_all", deprecated fields, computed-only fields
func getWritableAttributes(schema *awsschema.ResourceSchema) []string {
	attrNames := make([]string, 0, len(schema.Block.Attributes))
	for name, attr := range schema.Block.Attributes {
		if shouldIncludeAttribute(name, attr) {
			attrNames = append(attrNames, name)
		}
	}
	sort.Strings(attrNames)
	return attrNames
}

// shouldIncludeAttribute determines if an attribute should be included in generated variables and resources.
// Returns false if the attribute is:
//   - Deprecated
//   - Computed-only (computed=true, required=false, optional=false)
//   - Provider-managed (id, arn, tags_all)
//   - Not writable (according to IsWritable method)
func shouldIncludeAttribute(name string, attr awsschema.Attribute) bool {
	if attr.Deprecated {
		return false
	}

	// Skip computed-only attributes (read-only)
	if attr.Computed && !attr.Optional && !attr.Required {
		return false
	}

	// Skip provider-managed attributes even if they are optional+computed
	if attr.Computed && attr.Optional {
		if providerManagedAttributes[name] {
			return false
		}
	}

	return attr.IsWritable()
}
