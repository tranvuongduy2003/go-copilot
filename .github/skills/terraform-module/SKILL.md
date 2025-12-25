---
name: terraform-module
description: Create Terraform configurations and modules for AWS infrastructure
---

# Terraform Module Skill

This skill guides you through creating Terraform configurations and reusable modules for AWS infrastructure.

## Tech Stack

This project uses:
- **Backend**: Go 1.25+ with Chi router
- **Frontend**: React 19 with pnpm
- **Database**: PostgreSQL 16 (RDS)
- **Cache**: Redis (ElastiCache)
- **Container Orchestration**: Kubernetes (EKS)
- **CDN**: CloudFront + S3
- **Cloud Provider**: AWS

## When to Use This Skill

- Creating new cloud infrastructure
- Setting up multi-environment configurations
- Creating reusable Terraform modules
- Managing infrastructure state

## Project Structure

```
terraform/
├── main.tf                      # Main configuration
├── variables.tf                 # Input variables
├── outputs.tf                   # Output values
├── versions.tf                  # Provider versions
├── locals.tf                    # Local values
└── environments/
    ├── staging.tfvars           # Staging variables
    └── production.tfvars        # Production variables
```

## Templates

### Template 1: Provider Configuration

```hcl
# versions.tf
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
    random = {
      source  = "hashicorp/random"
      version = "~> 3.6"
    }
  }

  backend "s3" {
    bucket         = "terraform-state-bucket"
    key            = "app/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "terraform-locks"
  }
}

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

### Template 2: Variables

```hcl
# variables.tf
variable "aws_region" {
  description = "AWS region to deploy resources"
  type        = string
  default     = "us-east-1"
}

variable "project_name" {
  description = "Name of the project"
  type        = string
}

variable "environment" {
  description = "Environment (staging, production)"
  type        = string

  validation {
    condition     = contains(["staging", "production"], var.environment)
    error_message = "Environment must be 'staging' or 'production'."
  }
}

variable "vpc_cidr" {
  description = "CIDR block for VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "enable_nat_gateway" {
  description = "Enable NAT Gateway for private subnets"
  type        = bool
  default     = true
}

variable "db_instance_class" {
  description = "RDS instance class"
  type        = string
  default     = "db.t3.micro"
}

variable "tags" {
  description = "Additional tags"
  type        = map(string)
  default     = {}
}
```

### Template 3: Local Values

```hcl
# locals.tf
locals {
  name = "${var.project_name}-${var.environment}"

  azs = slice(data.aws_availability_zones.available.names, 0, 3)

  tags = merge(
    {
      Project     = var.project_name
      Environment = var.environment
      ManagedBy   = "terraform"
    },
    var.tags
  )
}

data "aws_availability_zones" "available" {
  state = "available"
}

data "aws_caller_identity" "current" {}
```

### Template 4: VPC Module

```hcl
# main.tf - VPC
module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "~> 5.0"

  name = "${local.name}-vpc"
  cidr = var.vpc_cidr

  azs             = local.azs
  private_subnets = [for k, v in local.azs : cidrsubnet(var.vpc_cidr, 4, k)]
  public_subnets  = [for k, v in local.azs : cidrsubnet(var.vpc_cidr, 8, k + 48)]

  enable_nat_gateway     = var.enable_nat_gateway
  single_nat_gateway     = var.environment != "production"
  one_nat_gateway_per_az = var.environment == "production"

  enable_dns_hostnames = true
  enable_dns_support   = true

  # VPC Flow Logs
  enable_flow_log                      = true
  create_flow_log_cloudwatch_iam_role  = true
  create_flow_log_cloudwatch_log_group = true

  public_subnet_tags = {
    "kubernetes.io/role/elb" = 1
  }

  private_subnet_tags = {
    "kubernetes.io/role/internal-elb" = 1
  }

  tags = local.tags
}
```

### Template 5: EKS Cluster

```hcl
# main.tf - EKS
module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 19.0"

  cluster_name    = "${local.name}-eks"
  cluster_version = "1.28"

  cluster_endpoint_public_access  = true
  cluster_endpoint_private_access = true

  vpc_id     = module.vpc.vpc_id
  subnet_ids = module.vpc.private_subnets

  enable_cluster_creator_admin_permissions = true

  eks_managed_node_groups = {
    general = {
      name           = "general"
      instance_types = var.environment == "production" ? ["t3.large"] : ["t3.medium"]

      min_size     = var.environment == "production" ? 3 : 1
      max_size     = var.environment == "production" ? 10 : 3
      desired_size = var.environment == "production" ? 3 : 2

      disk_size = 50

      labels = {
        role = "general"
      }
    }
  }

  cluster_addons = {
    coredns            = { most_recent = true }
    kube-proxy         = { most_recent = true }
    vpc-cni            = { most_recent = true }
    aws-ebs-csi-driver = { most_recent = true }
  }

  tags = local.tags
}
```

### Template 6: RDS Database

```hcl
# main.tf - RDS
module "rds" {
  source  = "terraform-aws-modules/rds/aws"
  version = "~> 6.0"

  identifier = "${local.name}-postgres"

  engine               = "postgres"
  engine_version       = "16.1"
  family               = "postgres16"
  major_engine_version = "16"
  instance_class       = var.db_instance_class

  allocated_storage     = var.environment == "production" ? 100 : 20
  max_allocated_storage = var.environment == "production" ? 500 : 100

  db_name  = var.db_name
  username = var.db_username
  port     = 5432

  multi_az               = var.environment == "production"
  db_subnet_group_name   = module.vpc.database_subnet_group_name
  vpc_security_group_ids = [module.rds_sg.security_group_id]

  backup_retention_period = var.environment == "production" ? 30 : 7
  skip_final_snapshot     = var.environment != "production"

  performance_insights_enabled = var.environment == "production"
  storage_encrypted            = true

  tags = local.tags
}

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

### Template 7: Outputs

```hcl
# outputs.tf
output "vpc_id" {
  description = "VPC ID"
  value       = module.vpc.vpc_id
}

output "eks_cluster_name" {
  description = "EKS cluster name"
  value       = module.eks.cluster_name
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

### Template 8: Environment tfvars

```hcl
# environments/staging.tfvars
aws_region   = "us-east-1"
project_name = "fullstack-app"
environment  = "staging"

vpc_cidr           = "10.0.0.0/16"
enable_nat_gateway = true

db_instance_class = "db.t3.micro"
db_name           = "app_staging"
db_username       = "postgres"

tags = {
  CostCenter = "engineering"
}
```

```hcl
# environments/production.tfvars
aws_region   = "us-east-1"
project_name = "fullstack-app"
environment  = "production"

vpc_cidr           = "10.0.0.0/16"
enable_nat_gateway = true

db_instance_class = "db.t3.medium"
db_name           = "app_production"
db_username       = "postgres"

tags = {
  CostCenter = "engineering"
  Compliance = "soc2"
}
```

## Commands

```bash
# Initialize
terraform init

# Format
terraform fmt -recursive

# Validate
terraform validate

# Plan
terraform plan -var-file=environments/staging.tfvars -out=tfplan

# Apply
terraform apply tfplan

# Destroy (staging only!)
terraform destroy -var-file=environments/staging.tfvars

# Show state
terraform state list
terraform state show <resource>

# Output
terraform output
terraform output -json
```

## Checklist

- [ ] Remote state with encryption and locking
- [ ] Provider versions pinned
- [ ] Variables have descriptions and validation
- [ ] Outputs documented
- [ ] Resources tagged consistently
- [ ] Security groups follow least privilege
- [ ] Encryption enabled for data at rest
- [ ] terraform fmt applied
- [ ] terraform validate passes
