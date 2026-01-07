package awsgen

import (
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/matt-FFFFFF/tfmodmake/hclgen"
	"github.com/zclconf/go-cty/cty"
)

// generateTerraform generates terraform.tf with AWS provider requirements
func generateTerraform(outputDir, region string) error {
	file := hclwrite.NewEmptyFile()
	body := file.Body()

	// terraform block
	tfBlock := body.AppendNewBlock("terraform", nil)
	tfBody := tfBlock.Body()

	// required_version
	tfBody.SetAttributeValue("required_version", cty.StringVal(">= 1.3"))

	// required_providers block
	providersBlock := tfBody.AppendNewBlock("required_providers", nil)
	providersBody := providersBlock.Body()

	// aws provider
	awsProviderTokens := hclwrite.Tokens{
		{Type: hclsyntax.TokenOBrace, Bytes: []byte("{")},
		{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")},
		{Type: hclsyntax.TokenIdent, Bytes: []byte("source")},
		{Type: hclsyntax.TokenEqual, Bytes: []byte(" = ")},
		{Type: hclsyntax.TokenQuotedLit, Bytes: []byte(`"hashicorp/aws"`)},
		{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")},
		{Type: hclsyntax.TokenIdent, Bytes: []byte("version")},
		{Type: hclsyntax.TokenEqual, Bytes: []byte(" = ")},
		{Type: hclsyntax.TokenQuotedLit, Bytes: []byte(`"~> 5.0"`)},
		{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")},
		{Type: hclsyntax.TokenCBrace, Bytes: []byte("}")},
	}
	providersBody.SetAttributeRaw("aws", awsProviderTokens)

	body.AppendNewline()

	// provider block (optional, with region if specified)
	if region != "" {
		providerBlock := body.AppendNewBlock("provider", []string{"aws"})
		providerBody := providerBlock.Body()
		providerBody.SetAttributeValue("region", cty.StringVal(region))
	}

	return hclgen.WriteFileToDir(outputDir, "terraform.tf", file)
}
