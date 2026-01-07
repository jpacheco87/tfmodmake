package awsgen

import (
	"sort"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/matt-FFFFFF/tfmodmake/awsschema"
	"github.com/matt-FFFFFF/tfmodmake/hclgen"
	"github.com/zclconf/go-cty/cty"
)

// generateOutputs generates outputs.tf with computed attributes
func generateOutputs(schema *awsschema.ResourceSchema, resourceType, outputDir string) error {
	file := hclwrite.NewEmptyFile()
	body := file.Body()

	// Get all computed attributes that users might want
	var computedAttrs []string
	for name, attr := range schema.Block.Attributes {
		// Output computed attributes (including those that are also writable)
		if attr.Computed && !attr.Deprecated {
			computedAttrs = append(computedAttrs, name)
		}
	}
	sort.Strings(computedAttrs)

	// Generate output for each computed attribute
	for i, name := range computedAttrs {
		if i > 0 {
			body.AppendNewline()
		}

		attr := schema.Block.Attributes[name]

		outputBlock := body.AppendNewBlock("output", []string{name})
		outputBody := outputBlock.Body()

		// Set value
		outputBody.SetAttributeRaw("value", hclgen.TokensForTraversal(resourceType, "this", name))

		// Set description
		desc := attr.Description
		if desc == "" {
			desc = "The " + name + " of the " + resourceType
		}
		outputBody.SetAttributeValue("description", cty.StringVal(desc))

		// Mark sensitive outputs
		if attr.Sensitive {
			outputBody.SetAttributeValue("sensitive", cty.True)
		}
	}

	return hclgen.WriteFileToDir(outputDir, "outputs.tf", file)
}
