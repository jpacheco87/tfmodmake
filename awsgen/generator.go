// Package awsgen provides functions to generate Terraform files for AWS resources.
package awsgen

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/matt-FFFFFF/tfmodmake/awsschema"
)

// GenerateOptions contains options for generating AWS Terraform modules
type GenerateOptions struct {
	ResourceType string // e.g., "aws_s3_bucket"
	OutputDir    string // Directory to write files to
	Region       string // Optional AWS region
}

// Generate generates a complete Terraform module for an AWS resource
func Generate(schema *awsschema.ResourceSchema, opts GenerateOptions) error {
	if err := generateTerraform(opts.OutputDir, opts.Region); err != nil {
		return fmt.Errorf("generating terraform.tf: %w", err)
	}

	if err := generateVariables(schema, opts.ResourceType, opts.OutputDir); err != nil {
		return fmt.Errorf("generating variables.tf: %w", err)
	}

	if err := generateMain(schema, opts.ResourceType, opts.OutputDir); err != nil {
		return fmt.Errorf("generating main.tf: %w", err)
	}

	if err := generateOutputs(schema, opts.ResourceType, opts.OutputDir); err != nil {
		return fmt.Errorf("generating outputs.tf: %w", err)
	}

	return nil
}

// writeFileToDir writes content to a file in the specified directory
func writeFileToDir(dir, filename string, content []byte) error {
	path := filepath.Join(dir, filename)
	return os.WriteFile(path, content, 0644)
}
