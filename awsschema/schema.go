// Package awsschema provides functions to parse AWS Terraform provider schemas.
package awsschema

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// ProviderSchemas represents the top-level structure from terraform providers schema -json
type ProviderSchemas struct {
	FormatVersion   string                     `json:"format_version"`
	ProviderSchemas map[string]ProviderSchema  `json:"provider_schemas"`
}

// ProviderSchema represents a single provider's schema
type ProviderSchema struct {
	Provider          ResourceSchema            `json:"provider"`
	ResourceSchemas   map[string]ResourceSchema `json:"resource_schemas"`
	DataSourceSchemas map[string]ResourceSchema `json:"data_source_schemas"`
}

// ResourceSchema represents a resource or data source schema
type ResourceSchema struct {
	Version int64 `json:"version"`
	Block   Block `json:"block"`
}

// Block represents a schema block with attributes and nested blocks
type Block struct {
	Attributes      map[string]Attribute `json:"attributes,omitempty"`
	BlockTypes      map[string]BlockType `json:"block_types,omitempty"`
	Description     string               `json:"description,omitempty"`
	DescriptionKind string               `json:"description_kind,omitempty"`
	Deprecated      bool                 `json:"deprecated,omitempty"`
}

// Attribute represents a single attribute in a schema
type Attribute struct {
	Type            json.RawMessage `json:"type"`
	Description     string          `json:"description,omitempty"`
	DescriptionKind string          `json:"description_kind,omitempty"`
	Required        bool            `json:"required,omitempty"`
	Optional        bool            `json:"optional,omitempty"`
	Computed        bool            `json:"computed,omitempty"`
	Sensitive       bool            `json:"sensitive,omitempty"`
	Deprecated      bool            `json:"deprecated,omitempty"`
}

// BlockType represents a nested block type
type BlockType struct {
	NestingMode string `json:"nesting_mode"`
	Block       Block  `json:"block"`
	MinItems    int64  `json:"min_items,omitempty"`
	MaxItems    int64  `json:"max_items,omitempty"`
}

// LoadSchema loads the AWS provider schema from a Terraform working directory
// The directory must have terraform initialized with the AWS provider
func LoadSchema(workingDir string) (*ProviderSchemas, error) {
	cmd := exec.Command("terraform", "providers", "schema", "-json")
	cmd.Dir = workingDir
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("terraform providers schema failed: %s", string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("failed to run terraform providers schema: %w", err)
	}

	var schemas ProviderSchemas
	if err := json.Unmarshal(output, &schemas); err != nil {
		return nil, fmt.Errorf("failed to parse provider schema JSON: %w", err)
	}

	return &schemas, nil
}

// FindResource finds a resource schema by resource type (e.g., "aws_instance", "aws_s3_bucket")
func FindResource(schemas *ProviderSchemas, resourceType string) (*ResourceSchema, error) {
	// Look for the AWS provider (try common variations)
	awsProviderKeys := []string{
		"registry.terraform.io/hashicorp/aws",
		"hashicorp/aws",
		"aws",
	}

	var provider *ProviderSchema
	for _, key := range awsProviderKeys {
		if p, ok := schemas.ProviderSchemas[key]; ok {
			provider = &p
			break
		}
	}

	if provider == nil {
		return nil, fmt.Errorf("AWS provider not found in schema (tried keys: %v)", awsProviderKeys)
	}

	resource, ok := provider.ResourceSchemas[resourceType]
	if !ok {
		return nil, fmt.Errorf("resource type %s not found in AWS provider", resourceType)
	}

	return &resource, nil
}

// GetAllResourceTypes returns all available AWS resource types from the schema
func GetAllResourceTypes(schemas *ProviderSchemas) ([]string, error) {
	awsProviderKeys := []string{
		"registry.terraform.io/hashicorp/aws",
		"hashicorp/aws",
		"aws",
	}

	var provider *ProviderSchema
	for _, key := range awsProviderKeys {
		if p, ok := schemas.ProviderSchemas[key]; ok {
			provider = &p
			break
		}
	}

	if provider == nil {
		return nil, fmt.Errorf("AWS provider not found in schema")
	}

	types := make([]string, 0, len(provider.ResourceSchemas))
	for resourceType := range provider.ResourceSchemas {
		types = append(types, resourceType)
	}

	return types, nil
}

// IsWritable returns true if an attribute is writable (not computed-only)
func (a *Attribute) IsWritable() bool {
	// An attribute is writable if it's required or optional
	// Computed-only attributes (computed=true, required=false, optional=false) are not writable
	return a.Required || a.Optional
}

// IsRequired returns true if an attribute is required
func (a *Attribute) IsRequired() bool {
	return a.Required
}

// GetTypeString returns a simple string representation of the attribute type
func (a *Attribute) GetTypeString() string {
	var typeData interface{}
	if err := json.Unmarshal(a.Type, &typeData); err != nil {
		return "string" // fallback
	}

	switch t := typeData.(type) {
	case string:
		return t
	case []interface{}:
		if len(t) > 0 {
			if typeStr, ok := t[0].(string); ok {
				switch typeStr {
				case "list", "set":
					return typeStr
				case "map":
					return "map"
				case "object":
					return "object"
				}
			}
		}
		return "list" // fallback for complex types
	default:
		return "string" // fallback
	}
}

// GetService extracts the AWS service name from a resource type
// e.g., "aws_s3_bucket" -> "s3", "aws_instance" -> "ec2"
func GetService(resourceType string) string {
	parts := strings.Split(resourceType, "_")
	if len(parts) >= 2 && parts[0] == "aws" {
		return parts[1]
	}
	return ""
}

// GetResourceName extracts the resource name from a resource type
// e.g., "aws_s3_bucket" -> "bucket", "aws_instance" -> "instance"
func GetResourceName(resourceType string) string {
	parts := strings.Split(resourceType, "_")
	if len(parts) >= 2 && parts[0] == "aws" {
		return strings.Join(parts[2:], "_")
	}
	return resourceType
}
