package awsgen

import (
	"sort"

	"github.com/hashicorp/hcl/v2/hclsyntax"
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

	// Get all writable attributes sorted by name
	var attrNames []string
	for name, attr := range schema.Block.Attributes {
		if attr.IsWritable() && !attr.Deprecated {
			attrNames = append(attrNames, name)
		}
	}
	sort.Strings(attrNames)

	// Set each attribute to its corresponding variable
	for _, name := range attrNames {
		resourceBody.SetAttributeRaw(name, hclgen.TokensForTraversal("var", name))
	}

	// Add dynamic blocks for nested block types
	var blockNames []string
	for name := range schema.Block.BlockTypes {
		blockNames = append(blockNames, name)
	}
	sort.Strings(blockNames)

	for _, name := range blockNames {
		blockType := schema.Block.BlockTypes[name]
		
		// Create dynamic block
		dynBlock := resourceBody.AppendNewBlock("dynamic", []string{name})
		dynBody := dynBlock.Body()

		// for_each determines when to create the block
		if blockType.MaxItems == 1 {
			// Single block: for_each on non-null variable
			dynBody.SetAttributeRaw("for_each", hclwrite.Tokens{
				{Type: hclsyntax.TokenIdent, Bytes: []byte("var")},
				{Type: hclsyntax.TokenDot, Bytes: []byte(".")},
				{Type: hclsyntax.TokenIdent, Bytes: []byte(name)},
				{Type: hclsyntax.TokenEqualOp, Bytes: []byte(" != ")},
				{Type: hclsyntax.TokenIdent, Bytes: []byte("null")},
				{Type: hclsyntax.TokenQuestion, Bytes: []byte(" ? ")},
				{Type: hclsyntax.TokenOBrack, Bytes: []byte("[")},
				{Type: hclsyntax.TokenIdent, Bytes: []byte("var")},
				{Type: hclsyntax.TokenDot, Bytes: []byte(".")},
				{Type: hclsyntax.TokenIdent, Bytes: []byte(name)},
				{Type: hclsyntax.TokenCBrack, Bytes: []byte("]")},
				{Type: hclsyntax.TokenColon, Bytes: []byte(" : ")},
				{Type: hclsyntax.TokenOBrack, Bytes: []byte("[")},
				{Type: hclsyntax.TokenCBrack, Bytes: []byte("]")},
			})
		} else {
			// Multiple blocks: for_each on variable (list/set)
			dynBody.SetAttributeRaw("for_each", hclgen.TokensForTraversal("var", name))
		}

		// content block - just pass through all attributes from the iterator
		contentBlock := dynBody.AppendNewBlock("content", nil)
		contentBody := contentBlock.Body()

		// For simplicity, set content as the iterator value
		// In a real implementation, we'd map each attribute in the block
		var attrNamesInBlock []string
		for attrName := range blockType.Block.Attributes {
			attrNamesInBlock = append(attrNamesInBlock, attrName)
		}
		sort.Strings(attrNamesInBlock)

		for _, attrName := range attrNamesInBlock {
			contentBody.SetAttributeRaw(attrName, hclgen.TokensForTraversal(name, "value", attrName))
		}
	}

	return hclgen.WriteFileToDir(outputDir, "main.tf", file)
}
