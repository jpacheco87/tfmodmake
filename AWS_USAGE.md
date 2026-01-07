# AWS Module Generation

tfmodmake can generate Terraform modules for AWS resources using the AWS provider schema.

## Prerequisites

1. Terraform installed (`>= 1.3`)
2. AWS provider initialized in a directory (for schema extraction)

## Quick Start

### Setup AWS Provider for Schema Extraction

Create a directory with AWS provider initialized:

```bash
mkdir ~/.tfmodmake-aws-schema
cd ~/.tfmodmake-aws-schema

cat > main.tf << 'PROVIDER'
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = "us-east-1"
  skip_credentials_validation = true
  skip_requesting_account_id = true
  skip_metadata_api_check = true
  access_key = "mock"
  secret_key = "mock"
}
PROVIDER

terraform init
```

### Generate AWS Modules

```bash
# Generate S3 bucket module
tfmodmake gen-aws --resource aws_s3_bucket --region us-east-1

# Generate Lambda function module
tfmodmake gen-aws --resource aws_lambda_function --region us-west-2

# Generate EC2 instance module
tfmodmake gen-aws --resource aws_instance --output ./modules/ec2

# Use custom schema directory
tfmodmake gen-aws --resource aws_vpc --terraform-dir ~/.tfmodmake-aws-schema
```

## Generated Files

The `gen-aws` command generates:

- `terraform.tf` - Provider requirements and configuration
- `variables.tf` - Input variables for all writable attributes
- `main.tf` - Resource definition using native AWS provider
- `outputs.tf` - All computed attributes (ARNs, IDs, etc.)

## Example Output

### S3 Bucket Module

```hcl
# terraform.tf
terraform {
  required_version = ">= 1.3"
  required_providers {
    aws = {
      source   = "hashicorp/aws"
      version  = "~> 5.0"
    }
  }
}

provider "aws" {
  region = "us-east-1"
}

# variables.tf
variable "bucket" {
  type        = string
  description = "The bucket for the aws_s3_bucket"
  default     = null
}

variable "force_destroy" {
  type        = bool
  description = "The force_destroy for the aws_s3_bucket"
  default     = null
}

variable "tags" {
  type        = map(string)
  description = "The tags for the aws_s3_bucket"
  default     = null
}

# main.tf
resource "aws_s3_bucket" "this" {
  bucket        = var.bucket
  force_destroy = var.force_destroy
  tags          = var.tags
}

# outputs.tf
output "arn" {
  value       = aws_s3_bucket.this.arn
  description = "The arn of the aws_s3_bucket"
}

output "bucket" {
  value       = aws_s3_bucket.this.bucket
  description = "The bucket of the aws_s3_bucket"
}

output "id" {
  value       = aws_s3_bucket.this.id
  description = "The id of the aws_s3_bucket"
}
```

## Supported Resources

Works with any AWS provider resource, including:

- `aws_s3_bucket`
- `aws_instance`
- `aws_lambda_function`
- `aws_vpc`
- `aws_security_group`
- `aws_rds_instance`
- `aws_iam_role`
- And 900+ other AWS resources

## CLI Options

```
tfmodmake gen-aws [options]

OPTIONS:
  --resource string       AWS resource type (e.g., aws_s3_bucket, aws_instance) [required]
  --region string         AWS region (optional, e.g., us-east-1)
  --output string         Output directory (default: current directory)
  --terraform-dir string  Directory with initialized Terraform AWS provider (default: /tmp/terraform-aws-test)
```

## Validation

All generated modules pass `terraform validate`:

```bash
cd my-generated-module
terraform init
terraform validate
# Success! The configuration is valid.
```

## Differences from Azure Generation

- **No nested `body` object**: AWS resources use flat attribute syntax
- **No `locals.tf`**: AWS provider doesn't need complex body reconstruction
- **Native resources**: Uses `aws_*` resources instead of `azapi_resource`
- **Simpler structure**: Attributes map directly to variables
- **No AVM interfaces**: AWS doesn't follow Azure Verified Modules patterns

## Tips

1. **Schema freshness**: Re-run `terraform init` in your schema directory periodically to get the latest provider schema
2. **Module customization**: Generated modules are starting points - customize as needed
3. **Complex nested blocks**: Some AWS resources have complex nested configurations - add these manually as needed
4. **Validation**: Always run `terraform validate` on generated modules before use
