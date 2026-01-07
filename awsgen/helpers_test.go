package awsgen

import (
	"testing"

	"github.com/matt-FFFFFF/tfmodmake/awsschema"
	"github.com/stretchr/testify/assert"
)

func TestShouldIncludeAttribute(t *testing.T) {
	tests := []struct {
		name     string
		attrName string
		attr     awsschema.Attribute
		expected bool
	}{
		{
			name:     "include required attribute",
			attrName: "bucket",
			attr: awsschema.Attribute{
				Required: true,
			},
			expected: true,
		},
		{
			name:     "include optional attribute",
			attrName: "force_destroy",
			attr: awsschema.Attribute{
				Optional: true,
			},
			expected: true,
		},
		{
			name:     "exclude deprecated attribute",
			attrName: "old_field",
			attr: awsschema.Attribute{
				Optional:   true,
				Deprecated: true,
			},
			expected: false,
		},
		{
			name:     "exclude computed-only attribute",
			attrName: "created_at",
			attr: awsschema.Attribute{
				Computed: true,
			},
			expected: false,
		},
		{
			name:     "exclude provider-managed id",
			attrName: "id",
			attr: awsschema.Attribute{
				Optional: true,
				Computed: true,
			},
			expected: false,
		},
		{
			name:     "exclude provider-managed arn",
			attrName: "arn",
			attr: awsschema.Attribute{
				Optional: true,
				Computed: true,
			},
			expected: false,
		},
		{
			name:     "exclude provider-managed tags_all",
			attrName: "tags_all",
			attr: awsschema.Attribute{
				Optional: true,
				Computed: true,
			},
			expected: false,
		},
		{
			name:     "include optional+computed non-provider-managed",
			attrName: "bucket_prefix",
			attr: awsschema.Attribute{
				Optional: true,
				Computed: true,
			},
			expected: true,
		},
		{
			name:     "include required+computed attribute",
			attrName: "function_name",
			attr: awsschema.Attribute{
				Required: true,
				Computed: true,
			},
			expected: true,
		},
		{
			name:     "exclude deprecated required attribute",
			attrName: "deprecated_required",
			attr: awsschema.Attribute{
				Required:   true,
				Deprecated: true,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shouldIncludeAttribute(tt.attrName, tt.attr)
			assert.Equal(t, tt.expected, result, "Attribute: %s", tt.attrName)
		})
	}
}

func TestGetWritableAttributes(t *testing.T) {
	schema := &awsschema.ResourceSchema{
		Block: awsschema.Block{
			Attributes: map[string]awsschema.Attribute{
				"bucket": {
					Required: true,
				},
				"force_destroy": {
					Optional: true,
				},
				"tags": {
					Optional: true,
				},
				"id": {
					Optional: true,
					Computed: true,
				},
				"arn": {
					Optional: true,
					Computed: true,
				},
				"tags_all": {
					Optional: true,
					Computed: true,
				},
				"created_at": {
					Computed: true,
				},
				"deprecated_field": {
					Optional:   true,
					Deprecated: true,
				},
				"bucket_domain_name": {
					Computed: true,
				},
			},
		},
	}

	result := getWritableAttributes(schema)

	// Should include: bucket, force_destroy, tags
	// Should exclude: id, arn, tags_all (provider-managed)
	// Should exclude: created_at, bucket_domain_name (computed-only)
	// Should exclude: deprecated_field (deprecated)
	expected := []string{"bucket", "force_destroy", "tags"}

	assert.Equal(t, expected, result)
}

func TestGetWritableAttributesEmpty(t *testing.T) {
	schema := &awsschema.ResourceSchema{
		Block: awsschema.Block{
			Attributes: map[string]awsschema.Attribute{},
		},
	}

	result := getWritableAttributes(schema)
	assert.Empty(t, result)
}

func TestGetWritableAttributesSorted(t *testing.T) {
	schema := &awsschema.ResourceSchema{
		Block: awsschema.Block{
			Attributes: map[string]awsschema.Attribute{
				"zebra": {
					Optional: true,
				},
				"apple": {
					Required: true,
				},
				"mango": {
					Optional: true,
				},
			},
		},
	}

	result := getWritableAttributes(schema)

	// Should be sorted alphabetically
	expected := []string{"apple", "mango", "zebra"}
	assert.Equal(t, expected, result)
}

func TestProviderManagedAttributes(t *testing.T) {
	// Test that all provider-managed attributes are defined in the map
	managedAttrs := []string{"id", "arn", "tags_all"}

	for _, attr := range managedAttrs {
		assert.True(t, providerManagedAttributes[attr], "Attribute %s should be in providerManagedAttributes", attr)
	}
}
