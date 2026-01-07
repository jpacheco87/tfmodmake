package awsgen

import (
	"sort"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/matt-FFFFFF/tfmodmake/awsschema"
	"github.com/matt-FFFFFF/tfmodmake/hclgen"
)

// generateMain generates main.tf with the AWS resource block
func generateMain(schema *awsschema.ResourceSchema, resourceType, outputDir string) error {
	file := hclwrite.NewEmptyFile()
	body := file.Body()

	// Create resource block
	resourceBlock := body.AppendNewBlock("resource", []string{resourceType, "this"})
	resourceBody := resourceBlock.Body()

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

	// Set each attribute to its corresponding variable
	for _, name := range attrNames {
		resourceBody.SetAttributeRaw(name, hclgen.TokensForTraversal("var", name))
	}

	// Note: We're intentionally not generating dynamic blocks for nested configurations
	// as they can have complex nested requirements. Users can add these manually as needed
	// or we can generate them in a future enhancement with proper schema introspection.

	return hclgen.WriteFileToDir(outputDir, "main.tf", file)
}
