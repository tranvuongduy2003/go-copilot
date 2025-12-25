# =============================================================================
# Staging Environment Variables
# =============================================================================

aws_region   = "us-east-1"
project_name = "fullstack-app"
environment  = "staging"

# Domain
domain_name = "staging.example.com"

# EKS Configuration
eks_cluster_version     = "1.28"
eks_node_instance_types = ["t3.medium"]
eks_node_min_size       = 1
eks_node_max_size       = 3
eks_node_desired_size   = 2

# Database Configuration
db_name              = "app_staging"
db_username          = "postgres"
db_instance_class    = "db.t3.micro"
db_allocated_storage = 20
db_backup_retention  = 7

# Redis Configuration
redis_node_type       = "cache.t3.micro"
redis_num_cache_nodes = 1

# Monitoring
enable_monitoring = true
enable_backup     = true

# Additional Tags
additional_tags = {
  CostCenter = "engineering"
  Team       = "platform"
}
