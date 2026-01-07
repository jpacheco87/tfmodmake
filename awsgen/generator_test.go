package awsgen

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/matt-FFFFFF/tfmodmake/awsschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateAWSS3Bucket(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Load AWS provider schema
	schemas, err := awsschema.LoadSchema("/tmp/terraform-aws-test")
	if err != nil {
		t.Skipf("Skipping test: %v", err)
		return
	}

	// Find aws_s3_bucket resource
	schema, err := awsschema.FindResource(schemas, "aws_s3_bucket")
	require.NoError(t, err)

	// Create temporary output directory
	tmpDir, err := os.MkdirTemp("", "tfmodmake-aws-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Generate module
	opts := GenerateOptions{
		ResourceType: "aws_s3_bucket",
		OutputDir:    tmpDir,
		Region:       "us-east-1",
	}

	err = Generate(schema, opts)
	require.NoError(t, err)

	// Verify files were created
	expectedFiles := []string{
		"terraform.tf",
		"variables.tf",
		"main.tf",
		"outputs.tf",
	}

	for _, filename := range expectedFiles {
		path := filepath.Join(tmpDir, filename)
		assert.FileExists(t, path, "Expected file %s to exist", filename)

		// Read and verify file is not empty
		content, err := os.ReadFile(path)
		require.NoError(t, err)
		assert.NotEmpty(t, content, "Expected file %s to have content", filename)
	}

	// Verify terraform.tf has AWS provider
	terraformContent, err := os.ReadFile(filepath.Join(tmpDir, "terraform.tf"))
	require.NoError(t, err)
	assert.Contains(t, string(terraformContent), "hashicorp/aws")
	assert.Contains(t, string(terraformContent), "us-east-1")

	// Verify variables.tf has some expected variables
	variablesContent, err := os.ReadFile(filepath.Join(tmpDir, "variables.tf"))
	require.NoError(t, err)
	assert.Contains(t, string(variablesContent), "variable")

	// Verify main.tf has resource block
	mainContent, err := os.ReadFile(filepath.Join(tmpDir, "main.tf"))
	require.NoError(t, err)
	assert.Contains(t, string(mainContent), "resource \"aws_s3_bucket\" \"this\"")

	// Verify outputs.tf has outputs
	outputsContent, err := os.ReadFile(filepath.Join(tmpDir, "outputs.tf"))
	require.NoError(t, err)
	assert.Contains(t, string(outputsContent), "output")
}

func TestGenerateAWSInstance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Load AWS provider schema
	schemas, err := awsschema.LoadSchema("/tmp/terraform-aws-test")
	if err != nil {
		t.Skipf("Skipping test: %v", err)
		return
	}

	// Find aws_instance resource
	schema, err := awsschema.FindResource(schemas, "aws_instance")
	require.NoError(t, err)

	// Create temporary output directory
	tmpDir, err := os.MkdirTemp("", "tfmodmake-aws-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Generate module
	opts := GenerateOptions{
		ResourceType: "aws_instance",
		OutputDir:    tmpDir,
	}

	err = Generate(schema, opts)
	require.NoError(t, err)

	// Verify main.tf has resource block with expected attributes
	mainContent, err := os.ReadFile(filepath.Join(tmpDir, "main.tf"))
	require.NoError(t, err)
	assert.Contains(t, string(mainContent), "resource \"aws_instance\" \"this\"")
	assert.Contains(t, string(mainContent), "ami")
	assert.Contains(t, string(mainContent), "instance_type")
}
