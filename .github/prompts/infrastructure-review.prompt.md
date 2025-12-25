---
description: Review infrastructure code for security, best practices, and cost optimization
---

# Infrastructure Review

Perform a comprehensive review of infrastructure code including Docker, Kubernetes, Terraform, and CI/CD configurations.

## Review Scope

**Components to Review**:
- [ ] Docker configurations
- [ ] Kubernetes manifests
- [ ] Terraform code
- [ ] CI/CD pipelines
- [ ] Monitoring setup

## Review Checklist

### Docker Review

#### Security
- [ ] Using specific base image versions (not `latest`)
- [ ] Running as non-root user
- [ ] No secrets in images or Dockerfiles
- [ ] Minimal base images (alpine preferred)
- [ ] Health checks defined
- [ ] .dockerignore configured properly

#### Performance
- [ ] Multi-stage builds used
- [ ] Layer order optimized for caching
- [ ] Unnecessary files excluded
- [ ] Image size reasonable

#### Best Practices
- [ ] Labels for metadata
- [ ] Single process per container
- [ ] Proper signal handling

### Kubernetes Review

#### Security
- [ ] Pods run as non-root
- [ ] Security contexts configured
- [ ] Network policies defined
- [ ] RBAC with least privilege
- [ ] Secrets management (Sealed Secrets/External Secrets)
- [ ] Pod Security Standards enforced

#### Reliability
- [ ] Resource requests and limits set
- [ ] Liveness probes configured
- [ ] Readiness probes configured
- [ ] Pod Disruption Budgets defined
- [ ] Horizontal Pod Autoscaler configured
- [ ] Anti-affinity rules for HA

#### Networking
- [ ] Services properly typed
- [ ] Ingress with TLS
- [ ] Security headers in ingress
- [ ] Rate limiting configured

### Terraform Review

#### Security
- [ ] Remote state with encryption
- [ ] State locking enabled
- [ ] No secrets in tfvars
- [ ] Security groups follow least privilege
- [ ] Encryption at rest enabled

#### Best Practices
- [ ] Provider versions pinned
- [ ] Module versions pinned
- [ ] Variables have descriptions
- [ ] Outputs documented
- [ ] Resources tagged consistently
- [ ] terraform fmt applied
- [ ] terraform validate passes

#### Cost Optimization
- [ ] Right-sized instances
- [ ] Auto-scaling configured
- [ ] Unused resources identified
- [ ] Reserved instances considered

### CI/CD Review

#### Security
- [ ] Secrets not hardcoded
- [ ] Security scanning enabled
- [ ] Dependency scanning enabled
- [ ] Container image scanning
- [ ] OIDC for cloud authentication

#### Reliability
- [ ] Tests run before deploy
- [ ] Deployments require approval
- [ ] Rollback capability exists
- [ ] Concurrency control configured
- [ ] Caching for performance

#### Best Practices
- [ ] Latest action versions used
- [ ] Path filters for efficiency
- [ ] Matrix builds where applicable
- [ ] Notifications configured

### Monitoring Review

#### Observability
- [ ] Metrics collected (Prometheus)
- [ ] Logs aggregated (Loki)
- [ ] Traces enabled (Tempo)
- [ ] Dashboards created

#### Alerting
- [ ] Critical alerts configured
- [ ] Alert routing defined
- [ ] Runbooks linked
- [ ] PagerDuty integration
- [ ] Slack notifications

## Review Output Template

```markdown
# Infrastructure Review Report

**Date**: {{date}}
**Reviewer**: {{reviewer}}
**Scope**: {{scope}}

## Summary

- **Critical Issues**: X
- **Warnings**: X
- **Recommendations**: X
- **Score**: X/10

## Critical Issues

### 1. [Issue Title]
- **Location**: `path/to/file`
- **Issue**: Description
- **Impact**: Security/Performance/Cost/Reliability
- **Recommendation**: How to fix
- **Priority**: P0/P1/P2

## Warnings

### 1. [Warning Title]
- **Location**: `path/to/file`
- **Issue**: Description
- **Recommendation**: How to fix

## Recommendations

### 1. [Recommendation Title]
- **Location**: `path/to/file`
- **Current State**: Description
- **Suggested Improvement**: How to improve
- **Benefit**: Expected benefit

## Security Findings

| Severity | Count | Description |
|----------|-------|-------------|
| Critical | X | Summary |
| High | X | Summary |
| Medium | X | Summary |
| Low | X | Summary |

## Cost Analysis

| Resource | Current | Recommended | Monthly Savings |
|----------|---------|-------------|-----------------|
| EC2 | t3.large | t3.medium | $XX |
| RDS | Multi-AZ | Single-AZ (staging) | $XX |

## Action Items

- [ ] P0: Fix critical security issue
- [ ] P1: Implement recommended changes
- [ ] P2: Review cost optimizations

## Next Steps

1. Address critical issues immediately
2. Schedule remediation for warnings
3. Plan implementation of recommendations
```

## Review Commands

```bash
# Terraform
terraform fmt -check -recursive
terraform validate
tflint
checkov -d terraform/

# Kubernetes
kubectl kustomize k8s/overlays/staging | kubeval
kube-score score k8s/base/*.yaml
kubesec scan k8s/base/*/deployment.yaml

# Docker
hadolint docker/Dockerfile.*
trivy config docker/
dockle ghcr.io/org/app:latest

# Security
gitleaks detect --source .
trivy fs --severity HIGH,CRITICAL .
```

## Output

Provide:
1. Review report following template
2. List of critical issues with fixes
3. Security findings summary
4. Cost optimization suggestions
5. Prioritized action items
