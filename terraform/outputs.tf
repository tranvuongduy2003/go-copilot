# =============================================================================
# Terraform Outputs
# =============================================================================

# =============================================================================
# VPC Outputs
# =============================================================================

output "vpc_id" {
  description = "The ID of the VPC"
  value       = module.vpc.vpc_id
}

output "vpc_cidr" {
  description = "The CIDR block of the VPC"
  value       = module.vpc.vpc_cidr_block
}

output "private_subnets" {
  description = "List of private subnet IDs"
  value       = module.vpc.private_subnets
}

output "public_subnets" {
  description = "List of public subnet IDs"
  value       = module.vpc.public_subnets
}

# =============================================================================
# EKS Outputs
# =============================================================================

output "eks_cluster_name" {
  description = "Name of the EKS cluster"
  value       = module.eks.cluster_name
}

output "eks_cluster_endpoint" {
  description = "Endpoint for the EKS cluster"
  value       = module.eks.cluster_endpoint
  sensitive   = true
}

output "eks_cluster_ca_certificate" {
  description = "Base64 encoded CA certificate for the EKS cluster"
  value       = module.eks.cluster_certificate_authority_data
  sensitive   = true
}

output "eks_cluster_oidc_issuer_url" {
  description = "OIDC issuer URL for the EKS cluster"
  value       = module.eks.cluster_oidc_issuer_url
}

output "eks_configure_kubectl" {
  description = "Command to configure kubectl"
  value       = "aws eks update-kubeconfig --region ${var.aws_region} --name ${module.eks.cluster_name}"
}

# =============================================================================
# RDS Outputs
# =============================================================================

output "rds_endpoint" {
  description = "Endpoint of the RDS instance"
  value       = module.rds.db_instance_endpoint
  sensitive   = true
}

output "rds_port" {
  description = "Port of the RDS instance"
  value       = module.rds.db_instance_port
}

output "rds_database_name" {
  description = "Name of the database"
  value       = module.rds.db_instance_name
}

# =============================================================================
# ElastiCache Outputs
# =============================================================================

output "redis_endpoint" {
  description = "Endpoint of the Redis cluster"
  value       = module.elasticache.cluster_cache_nodes[0].address
  sensitive   = true
}

output "redis_port" {
  description = "Port of the Redis cluster"
  value       = 6379
}

# =============================================================================
# S3 Outputs
# =============================================================================

output "s3_assets_bucket" {
  description = "Name of the S3 assets bucket"
  value       = module.s3_assets.s3_bucket_id
}

output "s3_assets_bucket_arn" {
  description = "ARN of the S3 assets bucket"
  value       = module.s3_assets.s3_bucket_arn
}

# =============================================================================
# CloudFront Outputs
# =============================================================================

output "cloudfront_distribution_id" {
  description = "ID of the CloudFront distribution"
  value       = module.cloudfront.cloudfront_distribution_id
}

output "cloudfront_domain_name" {
  description = "Domain name of the CloudFront distribution"
  value       = module.cloudfront.cloudfront_distribution_domain_name
}

# =============================================================================
# ALB Outputs
# =============================================================================

output "alb_dns_name" {
  description = "DNS name of the ALB"
  value       = module.alb.dns_name
}

output "alb_zone_id" {
  description = "Zone ID of the ALB"
  value       = module.alb.zone_id
}

# =============================================================================
# Secrets Manager Outputs
# =============================================================================

output "secrets_manager_secret_arn" {
  description = "ARN of the Secrets Manager secret"
  value       = aws_secretsmanager_secret.app_secrets.arn
}

# =============================================================================
# Connection Strings (for reference)
# =============================================================================

output "connection_info" {
  description = "Connection information for the application"
  value = {
    database_host = split(":", module.rds.db_instance_endpoint)[0]
    database_port = module.rds.db_instance_port
    redis_host    = module.elasticache.cluster_cache_nodes[0].address
    redis_port    = 6379
  }
  sensitive = true
}
