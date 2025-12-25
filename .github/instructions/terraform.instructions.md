---
applyTo: "terraform/**/*.tf,terraform/**/*.tfvars"
---

# Terraform Development Instructions

These instructions apply to all Terraform configurations for infrastructure as code.

## Project Terraform Structure

```
terraform/
├── main.tf                      # Main configuration
├── variables.tf                 # Input variables
├── outputs.tf                   # Output values
├── versions.tf                  # Provider versions (optional)
└── environments/
    ├── staging.tfvars           # Staging variables
    └── production.tfvars        # Production variables
```

## Provider Configuration

### Required Providers

```hcl
terraform {
  required_version = ">= 1.6.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.24"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.12"
    }
  }

  # Remote state with locking
  backend "s3" {
    bucket         = "terraform-state-bucket"
    key            = "app/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "terraform-locks"
  }
}
```

### Provider Configuration

```hcl
provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Project     = var.project_name
      Environment = var.environment
      ManagedBy   = "terraform"
    }
  }
}
```

## Variable Definitions

### variables.tf Pattern

```hcl
variable "aws_region" {
  description = "AWS region to deploy resources"
  type        = string
  default     = "us-east-1"
}

variable "environment" {
  description = "Environment name"
  type        = string

  validation {
    condition     = contains(["staging", "production"], var.environment)
    error_message = "Environment must be 'staging' or 'production'."
  }
}

variable "db_instance_class" {
  description = "RDS instance class"
  type        = string
  default     = "db.t3.micro"
}

variable "enable_monitoring" {
  description = "Enable enhanced monitoring"
  type        = bool
  default     = true
}

variable "tags" {
  description = "Additional tags"
  type        = map(string)
  default     = {}
}
```

### Environment Variables (tfvars)

```hcl
# environments/staging.tfvars
aws_region   = "us-east-1"
environment  = "staging"
project_name = "fullstack-app"

# EKS
eks_node_instance_types = ["t3.medium"]
eks_node_min_size       = 1
eks_node_max_size       = 3

# Database
db_instance_class    = "db.t3.micro"
db_allocated_storage = 20

# Redis
redis_node_type = "cache.t3.micro"
```

## Module Usage

### Using Community Modules

```hcl
module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "~> 5.0"

  name = "${var.project_name}-${var.environment}-vpc"
  cidr = "10.0.0.0/16"

  azs             = ["us-east-1a", "us-east-1b", "us-east-1c"]
  private_subnets = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
  public_subnets  = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]

  enable_nat_gateway = true
  single_nat_gateway = var.environment != "production"

  tags = local.tags
}

module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 19.0"

  cluster_name    = "${var.project_name}-${var.environment}-eks"
  cluster_version = "1.28"

  vpc_id     = module.vpc.vpc_id
  subnet_ids = module.vpc.private_subnets

  eks_managed_node_groups = {
    general = {
      instance_types = var.eks_node_instance_types
      min_size       = var.eks_node_min_size
      max_size       = var.eks_node_max_size
      desired_size   = var.eks_node_desired_size
    }
  }

  tags = local.tags
}
```

## Resource Patterns

### RDS Database

```hcl
module "rds" {
  source  = "terraform-aws-modules/rds/aws"
  version = "~> 6.0"

  identifier = "${local.name}-postgres"

  engine               = "postgres"
  engine_version       = "16.1"
  family               = "postgres16"
  instance_class       = var.db_instance_class

  allocated_storage     = var.db_allocated_storage
  max_allocated_storage = var.db_allocated_storage * 5

  db_name  = var.db_name
  username = var.db_username
  port     = 5432

  multi_az               = var.environment == "production"
  db_subnet_group_name   = module.vpc.database_subnet_group_name
  vpc_security_group_ids = [module.rds_sg.security_group_id]

  backup_retention_period = var.environment == "production" ? 30 : 7
  skip_final_snapshot     = var.environment != "production"

  storage_encrypted = true

  tags = local.tags
}
```

### Security Group

```hcl
module "rds_sg" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "~> 5.0"

  name        = "${local.name}-rds-sg"
  description = "Security group for RDS"
  vpc_id      = module.vpc.vpc_id

  ingress_with_source_security_group_id = [
    {
      from_port                = 5432
      to_port                  = 5432
      protocol                 = "tcp"
      source_security_group_id = module.eks.node_security_group_id
    }
  ]

  tags = local.tags
}
```

## Outputs

```hcl
output "vpc_id" {
  description = "VPC ID"
  value       = module.vpc.vpc_id
}

output "eks_cluster_endpoint" {
  description = "EKS cluster endpoint"
  value       = module.eks.cluster_endpoint
  sensitive   = true
}

output "rds_endpoint" {
  description = "RDS endpoint"
  value       = module.rds.db_instance_endpoint
  sensitive   = true
}

output "configure_kubectl" {
  description = "Configure kubectl command"
  value       = "aws eks update-kubeconfig --region ${var.aws_region} --name ${module.eks.cluster_name}"
}
```

## Local Values

```hcl
locals {
  name = "${var.project_name}-${var.environment}"

  tags = merge(
    {
      Project     = var.project_name
      Environment = var.environment
      ManagedBy   = "terraform"
    },
    var.tags
  )
}
```

## Common Commands

```bash
# Initialize
terraform init

# Format
terraform fmt -recursive

# Validate
terraform validate

# Plan (staging)
terraform plan -var-file=environments/staging.tfvars -out=tfplan

# Apply
terraform apply tfplan

# Show outputs
terraform output

# Destroy (staging only)
terraform destroy -var-file=environments/staging.tfvars

# Import existing resource
terraform import aws_s3_bucket.example my-bucket

# State management
terraform state list
terraform state show aws_s3_bucket.example
```

## Best Practices

### 1. State Management

- Use remote state (S3 + DynamoDB)
- Enable state locking
- Enable encryption
- Separate state per environment

### 2. Security

- Never commit secrets to tfvars
- Use AWS Secrets Manager or SSM Parameter Store
- Enable encryption for all storage
- Use least privilege IAM policies

### 3. Modules

- Use versioned modules
- Pin module versions
- Use community modules when available
- Create custom modules for reusability

### 4. Variables

- Always add descriptions
- Use validation blocks
- Provide sensible defaults
- Use type constraints

### 5. Naming

```hcl
# Good naming pattern
resource "aws_s3_bucket" "assets" {
  bucket = "${var.project_name}-${var.environment}-assets"
}

# Use locals for repeated names
locals {
  name_prefix = "${var.project_name}-${var.environment}"
}
```

## Checklist

- [ ] Remote state configured with locking
- [ ] Variables have descriptions and validation
- [ ] Sensitive outputs marked as sensitive
- [ ] Resources tagged consistently
- [ ] Security groups follow least privilege
- [ ] Encryption enabled for data at rest
- [ ] Separate tfvars per environment
- [ ] Module versions pinned
- [ ] terraform fmt applied
- [ ] terraform validate passes
