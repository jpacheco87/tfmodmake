package awsgen

import (
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

	// Get writable attributes using helper function
	attrNames := getWritableAttributes(schema)

	// Set each attribute to its corresponding variable
	for _, name := range attrNames {
		resourceBody.SetAttributeRaw(name, hclgen.TokensForTraversal("var", name))
	}

	// Note: We're intentionally not generating dynamic blocks for nested configurations
	// as they can have complex nested requirements. Users can add these manually as needed
	// or we can generate them in a future enhancement with proper schema introspection.

	return hclgen.WriteFileToDir(outputDir, "main.tf", file)
}
