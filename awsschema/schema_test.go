package awsschema

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetService(t *testing.T) {
	tests := []struct {
		name         string
		resourceType string
		expected     string
	}{
		{"S3 Bucket", "aws_s3_bucket", "s3"},
		{"EC2 Instance", "aws_instance", "instance"},
		{"Lambda Function", "aws_lambda_function", "lambda"},
		{"VPC", "aws_vpc", "vpc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetService(tt.resourceType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetResourceName(t *testing.T) {
	tests := []struct {
		name         string
		resourceType string
		expected     string
	}{
		{"S3 Bucket", "aws_s3_bucket", "bucket"},
		{"EC2 Instance", "aws_instance", ""},
		{"Lambda Function", "aws_lambda_function", "function"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetResourceName(tt.resourceType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAttributeIsWritable(t *testing.T) {
	tests := []struct {
		name     string
		attr     Attribute
		expected bool
	}{
		{
			name: "Required attribute is writable",
			attr: Attribute{
				Required: true,
				Optional: false,
				Computed: false,
			},
			expected: true,
		},
		{
			name: "Optional attribute is writable",
			attr: Attribute{
				Required: false,
				Optional: true,
				Computed: false,
			},
			expected: true,
		},
		{
			name: "Computed-only attribute is not writable",
			attr: Attribute{
				Required: false,
				Optional: false,
				Computed: true,
			},
			expected: false,
		},
		{
			name: "Optional+Computed attribute is writable",
			attr: Attribute{
				Required: false,
				Optional: true,
				Computed: true,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.attr.IsWritable()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAttributeGetTypeString(t *testing.T) {
	tests := []struct {
		name     string
		typeJSON string
		expected string
	}{
		{
			name:     "Simple string type",
			typeJSON: `"string"`,
			expected: "string",
		},
		{
			name:     "Simple bool type",
			typeJSON: `"bool"`,
			expected: "bool",
		},
		{
			name:     "List type",
			typeJSON: `["list", "string"]`,
			expected: "list",
		},
		{
			name:     "Map type",
			typeJSON: `["map", "string"]`,
			expected: "map",
		},
		{
			name:     "Set type",
			typeJSON: `["set", "string"]`,
			expected: "set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attr := Attribute{
				Type: []byte(tt.typeJSON),
			}
			result := attr.GetTypeString()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestLoadSchemaIntegration is an integration test that requires terraform to be installed
// and a terraform working directory with AWS provider initialized
func TestLoadSchemaIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// This test uses /tmp/terraform-aws-test if it exists
	schemas, err := LoadSchema("/tmp/terraform-aws-test")
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
		return
	}

	require.NotNil(t, schemas)
	assert.NotEmpty(t, schemas.ProviderSchemas)

	// Try to find aws_s3_bucket resource
	resource, err := FindResource(schemas, "aws_s3_bucket")
	require.NoError(t, err)
	require.NotNil(t, resource)

	// Check that it has expected attributes
	assert.Contains(t, resource.Block.Attributes, "bucket")
	assert.Contains(t, resource.Block.Attributes, "tags")

	// Check bucket attribute properties
	bucketAttr := resource.Block.Attributes["bucket"]
	assert.True(t, bucketAttr.Optional || bucketAttr.Required)
}

func TestFindResourceErrors(t *testing.T) {
	schemas := &ProviderSchemas{
		ProviderSchemas: map[string]ProviderSchema{
			"registry.terraform.io/hashicorp/aws": {
				ResourceSchemas: map[string]ResourceSchema{
					"aws_s3_bucket": {},
				},
			},
		},
	}

	t.Run("Resource not found", func(t *testing.T) {
		_, err := FindResource(schemas, "aws_nonexistent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("Provider not found", func(t *testing.T) {
		emptySchemas := &ProviderSchemas{
			ProviderSchemas: map[string]ProviderSchema{},
		}
		_, err := FindResource(emptySchemas, "aws_s3_bucket")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "AWS provider not found")
	})
}
