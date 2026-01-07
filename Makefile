# Default output directory for generated modules
OUTPUT_DIR ?= ./output

# AWS provider schema directory (needs terraform with AWS provider initialized)
AWS_SCHEMA_DIR ?= /tmp/terraform-aws-test

# Build the CLI tool
.PHONY: build
build:
	go build -o tfmodmake ./cmd/tfmodmake

# Install the CLI tool
.PHONY: install
install:
	go install ./cmd/tfmodmake

# Run all tests
.PHONY: test
test:
	go test -count=1 ./...

# Run example tests
.PHONY: test-examples
test-examples:
	./scripts/test_examples.sh

# Setup AWS provider schema directory for AWS module generation
.PHONY: setup-aws-schema
setup-aws-schema:
	@echo "Setting up AWS provider schema directory..."
	@mkdir -p $(AWS_SCHEMA_DIR)
	@cd $(AWS_SCHEMA_DIR) && \
	if [ ! -f main.tf ]; then \
		echo 'terraform {' > main.tf; \
		echo '  required_providers {' >> main.tf; \
		echo '    aws = {' >> main.tf; \
		echo '      source  = "hashicorp/aws"' >> main.tf; \
		echo '      version = "~> 5.0"' >> main.tf; \
		echo '    }' >> main.tf; \
		echo '  }' >> main.tf; \
		echo '}' >> main.tf; \
		echo '' >> main.tf; \
		echo 'provider "aws" {' >> main.tf; \
		echo '  region = "us-east-1"' >> main.tf; \
		echo '  skip_credentials_validation = true' >> main.tf; \
		echo '  skip_requesting_account_id = true' >> main.tf; \
		echo '  skip_metadata_api_check = true' >> main.tf; \
		echo '  access_key = "mock"' >> main.tf; \
		echo '  secret_key = "mock"' >> main.tf; \
		echo '}' >> main.tf; \
		terraform -chdir=$(AWS_SCHEMA_DIR) init; \
	fi
	@echo "✓ AWS provider schema directory ready at $(AWS_SCHEMA_DIR)"

# Generate AWS S3 bucket module
.PHONY: gen-aws-s3
gen-aws-s3: build
	@mkdir -p $(OUTPUT_DIR)/aws-s3-bucket
	./tfmodmake gen-aws --resource aws_s3_bucket --region us-east-1 --output $(OUTPUT_DIR)/aws-s3-bucket --terraform-dir $(AWS_SCHEMA_DIR)
	@echo "✓ Generated AWS S3 bucket module in $(OUTPUT_DIR)/aws-s3-bucket"

# Generate AWS Lambda function module
.PHONY: gen-aws-lambda
gen-aws-lambda: build
	@mkdir -p $(OUTPUT_DIR)/aws-lambda-function
	./tfmodmake gen-aws --resource aws_lambda_function --region us-east-1 --output $(OUTPUT_DIR)/aws-lambda-function --terraform-dir $(AWS_SCHEMA_DIR)
	@echo "✓ Generated AWS Lambda function module in $(OUTPUT_DIR)/aws-lambda-function"

# Generate AWS EC2 instance module
.PHONY: gen-aws-ec2
gen-aws-ec2: build
	@mkdir -p $(OUTPUT_DIR)/aws-instance
	./tfmodmake gen-aws --resource aws_instance --region us-east-1 --output $(OUTPUT_DIR)/aws-instance --terraform-dir $(AWS_SCHEMA_DIR)
	@echo "✓ Generated AWS EC2 instance module in $(OUTPUT_DIR)/aws-instance"

# Generate AWS VPC module
.PHONY: gen-aws-vpc
gen-aws-vpc: build
	@mkdir -p $(OUTPUT_DIR)/aws-vpc
	./tfmodmake gen-aws --resource aws_vpc --region us-east-1 --output $(OUTPUT_DIR)/aws-vpc --terraform-dir $(AWS_SCHEMA_DIR)
	@echo "✓ Generated AWS VPC module in $(OUTPUT_DIR)/aws-vpc"

# Generate AWS RDS instance module
.PHONY: gen-aws-rds
gen-aws-rds: build
	@mkdir -p $(OUTPUT_DIR)/aws-rds-instance
	./tfmodmake gen-aws --resource aws_db_instance --region us-east-1 --output $(OUTPUT_DIR)/aws-rds-instance --terraform-dir $(AWS_SCHEMA_DIR)
	@echo "✓ Generated AWS RDS instance module in $(OUTPUT_DIR)/aws-rds-instance"

# Generate AWS IAM role module
.PHONY: gen-aws-iam-role
gen-aws-iam-role: build
	@mkdir -p $(OUTPUT_DIR)/aws-iam-role
	./tfmodmake gen-aws --resource aws_iam_role --region us-east-1 --output $(OUTPUT_DIR)/aws-iam-role --terraform-dir $(AWS_SCHEMA_DIR)
	@echo "✓ Generated AWS IAM role module in $(OUTPUT_DIR)/aws-iam-role"

# Generate AWS Security Group module
.PHONY: gen-aws-sg
gen-aws-sg: build
	@mkdir -p $(OUTPUT_DIR)/aws-security-group
	./tfmodmake gen-aws --resource aws_security_group --region us-east-1 --output $(OUTPUT_DIR)/aws-security-group --terraform-dir $(AWS_SCHEMA_DIR)
	@echo "✓ Generated AWS Security Group module in $(OUTPUT_DIR)/aws-security-group"

# Generate all common AWS modules
.PHONY: gen-aws-all
gen-aws-all: gen-aws-s3 gen-aws-lambda gen-aws-ec2 gen-aws-vpc gen-aws-rds gen-aws-iam-role gen-aws-sg
	@echo "✓ Generated all common AWS modules"

# Generate custom AWS resource (usage: make gen-aws RESOURCE=aws_dynamodb_table REGION=us-west-2)
.PHONY: gen-aws
gen-aws: build
	@if [ -z "$(RESOURCE)" ]; then \
		echo "Error: RESOURCE is required. Usage: make gen-aws RESOURCE=aws_dynamodb_table"; \
		exit 1; \
	fi
	@mkdir -p $(OUTPUT_DIR)/$(RESOURCE)
	./tfmodmake gen-aws --resource $(RESOURCE) --region $(REGION) --output $(OUTPUT_DIR)/$(RESOURCE) --terraform-dir $(AWS_SCHEMA_DIR)
	@echo "✓ Generated $(RESOURCE) module in $(OUTPUT_DIR)/$(RESOURCE)"

# Clean generated output
.PHONY: clean
clean:
	rm -rf $(OUTPUT_DIR)
	rm -f tfmodmake
	@echo "✓ Cleaned generated output and binary"

# Help target
.PHONY: help
help:
	@echo "tfmodmake - Terraform Module Generator"
	@echo ""
	@echo "Build & Test Commands:"
	@echo "  make build              Build the CLI tool"
	@echo "  make install            Install the CLI tool"
	@echo "  make test               Run all tests"
	@echo "  make test-examples      Run example tests"
	@echo ""
	@echo "AWS Setup:"
	@echo "  make setup-aws-schema   Setup AWS provider schema directory (required once)"
	@echo ""
	@echo "AWS Module Generation (Common Resources):"
	@echo "  make gen-aws-s3         Generate AWS S3 bucket module"
	@echo "  make gen-aws-lambda     Generate AWS Lambda function module"
	@echo "  make gen-aws-ec2        Generate AWS EC2 instance module"
	@echo "  make gen-aws-vpc        Generate AWS VPC module"
	@echo "  make gen-aws-rds        Generate AWS RDS instance module"
	@echo "  make gen-aws-iam-role   Generate AWS IAM role module"
	@echo "  make gen-aws-sg         Generate AWS Security Group module"
	@echo "  make gen-aws-all        Generate all common AWS modules"
	@echo ""
	@echo "AWS Module Generation (Custom):"
	@echo "  make gen-aws RESOURCE=aws_dynamodb_table [REGION=us-east-1]"
	@echo "                          Generate custom AWS resource module"
	@echo ""
	@echo "Configuration:"
	@echo "  OUTPUT_DIR=./output     Change output directory (default: ./output)"
	@echo "  AWS_SCHEMA_DIR=/path    Change AWS schema directory (default: /tmp/terraform-aws-test)"
	@echo "  REGION=us-west-2        Change AWS region (default: us-east-1)"
	@echo ""
	@echo "Utility:"
	@echo "  make clean              Clean generated output and binary"
	@echo "  make help               Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make gen-aws-s3"
	@echo "  make gen-aws RESOURCE=aws_dynamodb_table REGION=us-west-2"
	@echo "  OUTPUT_DIR=./modules make gen-aws-lambda"
