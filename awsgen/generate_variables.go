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

	// Get all writable attributes sorted by name
	var attrNames []string
	for name, attr := range schema.Block.Attributes {
		if attr.IsWritable() && !attr.Deprecated {
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

		// Set type
		varBody.SetAttributeValue("type", cty.StringVal(getTerraformType(attr)))

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

	// Generate variables for block types
	var blockNames []string
	for name := range schema.Block.BlockTypes {
		blockNames = append(blockNames, name)
	}
	sort.Strings(blockNames)

	for _, name := range blockNames {
		body.AppendNewline()
		
		blockType := schema.Block.BlockTypes[name]
		
		varBlock := body.AppendNewBlock("variable", []string{name})
		varBody := varBlock.Body()

		// Nested blocks are typically represented as list(object) or set(object)
		typeStr := "list(any)"
		if blockType.NestingMode == "set" {
			typeStr = "set(any)"
		} else if blockType.MaxItems == 1 {
			typeStr = "object(any)"
		}
		varBody.SetAttributeValue("type", cty.StringVal(typeStr))

		desc := blockType.Block.Description
		if desc == "" {
			desc = fmt.Sprintf("Configuration for %s", name)
		}
		varBody.SetAttributeValue("description", cty.StringVal(desc))

		// Default to empty for optional blocks
		if blockType.MinItems == 0 {
			if blockType.MaxItems == 1 {
				varBody.SetAttributeValue("default", cty.NullVal(cty.DynamicPseudoType))
			} else {
				varBody.SetAttributeValue("default", cty.ListValEmpty(cty.DynamicPseudoType))
			}
		}
	}

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
