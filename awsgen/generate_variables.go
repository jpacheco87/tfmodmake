package awsgen

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/matt-FFFFFF/tfmodmake/awsschema"
	"github.com/matt-FFFFFF/tfmodmake/hclgen"
	"github.com/zclconf/go-cty/cty"
)

// generateVariables generates variables.tf for an AWS resource
func generateVariables(schema *awsschema.ResourceSchema, resourceType, outputDir string) error {
	file := hclwrite.NewEmptyFile()
	body := file.Body()

	// Get all truly writable attributes (required or optional, but not computed-only)
	// Skip: computed-only (computed=true, required=false, optional=false)
	// Skip: optional+computed that are typically managed by provider (id, tags_all, etc.)
	var attrNames []string
	for name, attr := range schema.Block.Attributes {
		if attr.Deprecated {
			continue
		}
		// Skip computed-only attributes
		if attr.Computed && !attr.Optional && !attr.Required {
			continue
		}
		// Skip common provider-managed computed attributes even if optional
		if attr.Computed && attr.Optional {
			// Common AWS provider-managed attributes
			if name == "id" || name == "arn" || name == "tags_all" {
				continue
			}
		}
		if attr.IsWritable() {
			attrNames = append(attrNames, name)
		}
	}
	sort.Strings(attrNames)

	// Generate variable for each writable attribute
	for i, name := range attrNames {
		if i > 0 {
			body.AppendNewline()
		}

		attr := schema.Block.Attributes[name]
		
		varBlock := body.AppendNewBlock("variable", []string{name})
		varBody := varBlock.Body()

		// Set type (use raw tokens for unquoted type)
		varBody.SetAttributeRaw("type", hclwrite.TokensForIdentifier(getTerraformType(attr)))

		// Set description
		if attr.Description != "" {
			desc := strings.TrimSpace(attr.Description)
			if len(desc) < 80 {
				varBody.SetAttributeValue("description", cty.StringVal(desc))
			} else {
				varBody.SetAttributeRaw("description", hclgen.TokensForHeredoc(desc))
			}
		} else {
			varBody.SetAttributeValue("description", cty.StringVal(fmt.Sprintf("The %s for the %s", name, resourceType)))
		}

		// Set default to null for optional attributes
		if !attr.Required {
			varBody.SetAttributeValue("default", cty.NullVal(cty.DynamicPseudoType))
		}

		// Mark sensitive attributes
		if attr.Sensitive {
			varBody.SetAttributeValue("sensitive", cty.True)
		}
	}

	// Note: We're intentionally not generating variables for block types
	// as they have complex nested structures that are better handled manually
	// Users can add these as needed

	return hclgen.WriteFileToDir(outputDir, "variables.tf", file)
}

// getTerraformType converts AWS schema type to Terraform type string
func getTerraformType(attr awsschema.Attribute) string {
	typeStr := attr.GetTypeString()
	
	switch typeStr {
	case "string":
		return "string"
	case "bool":
		return "bool"
	case "number":
		return "number"
	case "list":
		return "list(any)"
	case "set":
		return "set(any)"
	case "map":
		return "map(string)"
	case "object":
		return "any"
	default:
		return "any"
	}
}
