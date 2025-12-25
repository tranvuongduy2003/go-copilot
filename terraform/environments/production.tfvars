# =============================================================================
# Production Environment Variables
# =============================================================================

aws_region   = "us-east-1"
project_name = "fullstack-app"
environment  = "production"

# Domain
domain_name = "example.com"

# EKS Configuration
eks_cluster_version     = "1.28"
eks_node_instance_types = ["t3.large", "t3.xlarge"]
eks_node_min_size       = 3
eks_node_max_size       = 20
eks_node_desired_size   = 5

# Database Configuration
db_name              = "app_production"
db_username          = "postgres"
db_instance_class    = "db.t3.medium"
db_allocated_storage = 100
db_backup_retention  = 30

# Redis Configuration
redis_node_type       = "cache.t3.medium"
redis_num_cache_nodes = 2

# Monitoring
enable_monitoring = true
enable_backup     = true

# Additional Tags
additional_tags = {
  CostCenter  = "engineering"
  Team        = "platform"
  Compliance  = "soc2"
  DataClass   = "confidential"
}
