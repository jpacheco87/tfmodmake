# Makefile Usage Guide

This guide explains how to use the Makefile to easily generate Terraform modules.

## Quick Start

### 1. Setup (One-time)

First, setup the AWS provider schema directory:

```bash
make setup-aws-schema
```

This creates a directory with an initialized Terraform AWS provider that's used for schema extraction.

### 2. Generate AWS Modules

Generate common AWS resource modules:

```bash
# Generate AWS S3 bucket module
make gen-aws-s3

# Generate AWS Lambda function module
make gen-aws-lambda

# Generate AWS EC2 instance module
make gen-aws-ec2

# Generate all common AWS modules
make gen-aws-all
```

### 3. Generate Custom AWS Resources

Generate any AWS resource type:

```bash
# DynamoDB table
make gen-aws RESOURCE=aws_dynamodb_table

# ECS cluster with custom region
make gen-aws RESOURCE=aws_ecs_cluster REGION=us-west-2

# CloudFront distribution with custom output directory
OUTPUT_DIR=./my-modules make gen-aws RESOURCE=aws_cloudfront_distribution
```

## Available Commands

### Build & Test

```bash
make build              # Build the CLI tool
make install            # Install the CLI tool globally
make test               # Run all tests
make test-examples      # Run example tests
```

### AWS Setup

```bash
make setup-aws-schema   # Setup AWS provider schema directory (required once)
```

### AWS Module Generation (Common Resources)

```bash
make gen-aws-s3         # Generate AWS S3 bucket module
make gen-aws-lambda     # Generate AWS Lambda function module
make gen-aws-ec2        # Generate AWS EC2 instance module
make gen-aws-vpc        # Generate AWS VPC module
make gen-aws-rds        # Generate AWS RDS instance module
make gen-aws-iam-role   # Generate AWS IAM role module
make gen-aws-sg         # Generate AWS Security Group module
make gen-aws-all        # Generate all common AWS modules
```

### AWS Module Generation (Custom)

```bash
make gen-aws RESOURCE=<aws_resource_type> [REGION=<region>]
```

**Examples:**

```bash
# DynamoDB Table
make gen-aws RESOURCE=aws_dynamodb_table

# SQS Queue in us-west-2
make gen-aws RESOURCE=aws_sqs_queue REGION=us-west-2

# API Gateway REST API
make gen-aws RESOURCE=aws_api_gateway_rest_api

# EKS Cluster
make gen-aws RESOURCE=aws_eks_cluster
```

### Utility Commands

```bash
make clean              # Clean generated output and binary
make help               # Show help message with all commands
```

## Configuration Options

You can customize the behavior using environment variables:

### Output Directory

```bash
# Default output directory is ./output
OUTPUT_DIR=./my-modules make gen-aws-s3

# Generate multiple modules in custom directory
OUTPUT_DIR=./terraform-modules make gen-aws-all
```

### AWS Schema Directory

```bash
# Use a custom directory for AWS provider schema
AWS_SCHEMA_DIR=~/.tfmodmake-aws make gen-aws-s3
```

### AWS Region

```bash
# Specify a custom AWS region
make gen-aws RESOURCE=aws_s3_bucket REGION=eu-west-1
```

## Generated Output Structure

By default, modules are generated in the `./output` directory:

```
output/
├── aws-s3-bucket/
│   ├── terraform.tf      # Provider requirements
│   ├── variables.tf      # Input variables
│   ├── main.tf           # Resource definition
│   └── outputs.tf        # Computed outputs
├── aws-lambda-function/
│   └── ...
└── aws-instance/
    └── ...
```

Each module contains:
- `terraform.tf` - Provider requirements (AWS ~> 5.0)
- `variables.tf` - Typed input variables for all resource attributes
- `main.tf` - Native AWS resource block
- `outputs.tf` - All computed attributes (ARNs, IDs, etc.)

## Examples

### Generate a Single Module

```bash
# Build and generate S3 bucket module
make gen-aws-s3

# The module will be in ./output/aws-s3-bucket/
cd output/aws-s3-bucket
terraform init
terraform validate
```

### Generate Multiple Modules

```bash
# Generate all common AWS modules
make gen-aws-all

# This creates:
# - output/aws-s3-bucket/
# - output/aws-lambda-function/
# - output/aws-instance/
# - output/aws-vpc/
# - output/aws-rds-instance/
# - output/aws-iam-role/
# - output/aws-security-group/
```

### Generate Custom Resources

```bash
# Generate a DynamoDB table module
make gen-aws RESOURCE=aws_dynamodb_table

# Generate an SNS topic module in us-west-2
make gen-aws RESOURCE=aws_sns_topic REGION=us-west-2

# Generate a CloudWatch log group with custom output directory
OUTPUT_DIR=./modules make gen-aws RESOURCE=aws_cloudwatch_log_group
```

### Clean Up

```bash
# Remove all generated modules and binary
make clean

# Then regenerate specific modules
make gen-aws-s3
make gen-aws-lambda
```

## Workflow Example

Complete workflow for generating and using AWS modules:

```bash
# 1. Setup (first time only)
make setup-aws-schema

# 2. Generate modules you need
make gen-aws-s3
make gen-aws RESOURCE=aws_dynamodb_table

# 3. Verify generated modules
cd output/aws-s3-bucket
terraform init
terraform validate

# 4. Customize as needed
# Edit variables.tf or main.tf to add specific configurations

# 5. Use in your Terraform code
# module "s3_bucket" {
#   source = "./output/aws-s3-bucket"
#   bucket = "my-unique-bucket-name"
#   tags   = { Environment = "dev" }
# }

# 6. Clean up when done
cd ../..
make clean
```

## Tips

1. **Schema Freshness**: Periodically re-run `make setup-aws-schema` to get the latest AWS provider schema
2. **Version Control**: The `output/` directory is in `.gitignore` - generated modules are meant as starting points
3. **Customization**: Generated modules are templates - customize them for your specific needs
4. **Documentation**: Each variable includes a description from the AWS provider schema
5. **Validation**: Always run `terraform validate` on generated modules before use

## Troubleshooting

### "terraform: not found"

Install Terraform:
```bash
wget https://releases.hashicorp.com/terraform/1.10.5/terraform_1.10.5_linux_amd64.zip
unzip terraform_1.10.5_linux_amd64.zip
sudo mv terraform /usr/local/bin/
```

### "failed to load AWS provider schema"

Run the setup command:
```bash
make setup-aws-schema
```

### Invalid resource type

Verify the resource type exists in the AWS provider documentation:
https://registry.terraform.io/providers/hashicorp/aws/latest/docs

## Support for 900+ AWS Resources

The tool works with any AWS provider resource type. Common examples:

**Compute:**
- `aws_instance`, `aws_launch_template`, `aws_autoscaling_group`
- `aws_lambda_function`, `aws_lambda_layer_version`
- `aws_ecs_cluster`, `aws_ecs_service`, `aws_ecs_task_definition`

**Storage:**
- `aws_s3_bucket`, `aws_s3_bucket_policy`
- `aws_ebs_volume`, `aws_efs_file_system`

**Database:**
- `aws_db_instance`, `aws_db_cluster`
- `aws_dynamodb_table`, `aws_elasticache_cluster`

**Networking:**
- `aws_vpc`, `aws_subnet`, `aws_route_table`
- `aws_security_group`, `aws_network_interface`
- `aws_lb`, `aws_lb_target_group`

**And many more...**

For a complete list, see: https://registry.terraform.io/providers/hashicorp/aws/latest/docs
