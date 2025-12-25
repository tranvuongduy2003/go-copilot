---
applyTo: "k8s/**/*"
---

# Kubernetes Development Instructions

These instructions apply to all Kubernetes manifests managed with Kustomize.

## Project Kubernetes Structure

```
k8s/
├── base/                            # Base manifests
│   ├── kustomization.yaml           # Base kustomization
│   ├── namespace.yaml               # Namespace definition
│   ├── ingress.yaml                 # Ingress + Network policies
│   ├── frontend/
│   │   ├── deployment.yaml          # Frontend deployment
│   │   └── service.yaml             # Service + HPA + PDB
│   └── backend/
│       ├── deployment.yaml          # Backend deployment
│       └── service.yaml             # Service + HPA + PDB
│
└── overlays/
    ├── staging/                     # Staging environment
    │   ├── kustomization.yaml       # Patches + images
    │   └── namespace.yaml           # Staging namespace
    └── production/                  # Production environment
        ├── kustomization.yaml       # Patches + images
        ├── namespace.yaml           # Namespace + quotas
        └── sealed-secrets.yaml      # Encrypted secrets
```

## Kustomize Usage

### Base Kustomization

```yaml
# k8s/base/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: app

commonLabels:
  app.kubernetes.io/name: fullstack-app
  app.kubernetes.io/managed-by: kustomize

resources:
  - namespace.yaml
  - frontend/deployment.yaml
  - frontend/service.yaml
  - backend/deployment.yaml
  - backend/service.yaml
  - ingress.yaml
```

### Overlay Kustomization

```yaml
# k8s/overlays/staging/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: staging
namePrefix: staging-

commonLabels:
  environment: staging

resources:
  - ../../base
  - namespace.yaml

patches:
  - patch: |-
      - op: replace
        path: /spec/replicas
        value: 1
    target:
      kind: Deployment

images:
  - name: frontend
    newName: ghcr.io/myorg/app/frontend
    newTag: staging
  - name: backend
    newName: ghcr.io/myorg/app/backend
    newTag: staging

configMapGenerator:
  - name: app-config
    behavior: merge
    literals:
      - API_URL=https://api.staging.example.com
```

## Deployment Patterns

### Standard Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: backend
  labels:
    app: backend
spec:
  replicas: 2
  selector:
    matchLabels:
      app: backend
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  template:
    metadata:
      labels:
        app: backend
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "8080"
    spec:
      serviceAccountName: backend
      securityContext:
        runAsNonRoot: true
        runAsUser: 1001
        fsGroup: 1001
      containers:
        - name: backend
          image: backend:latest
          imagePullPolicy: Always
          ports:
            - name: http
              containerPort: 8080
          env:
            - name: NODE_ENV
              value: "production"
            - name: DATABASE_URL
              valueFrom:
                secretKeyRef:
                  name: backend-secrets
                  key: DATABASE_URL
          resources:
            requests:
              cpu: "100m"
              memory: "256Mi"
            limits:
              cpu: "500m"
              memory: "512Mi"
          livenessProbe:
            httpGet:
              path: /health
              port: http
            initialDelaySeconds: 30
            periodSeconds: 10
          readinessProbe:
            httpGet:
              path: /health/ready
              port: http
            initialDelaySeconds: 5
            periodSeconds: 5
          securityContext:
            allowPrivilegeEscalation: false
            readOnlyRootFilesystem: true
            capabilities:
              drop: ["ALL"]
          volumeMounts:
            - name: tmp
              mountPath: /tmp
      volumes:
        - name: tmp
          emptyDir: {}
      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
            - weight: 100
              podAffinityTerm:
                labelSelector:
                  matchLabels:
                    app: backend
                topologyKey: kubernetes.io/hostname
```

## Service Patterns

### ClusterIP Service with HPA

```yaml
apiVersion: v1
kind: Service
metadata:
  name: backend
spec:
  type: ClusterIP
  ports:
    - port: 8080
      targetPort: http
      name: http
  selector:
    app: backend

---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: backend
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: backend
  minReplicas: 2
  maxReplicas: 10
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
  behavior:
    scaleDown:
      stabilizationWindowSeconds: 300
    scaleUp:
      stabilizationWindowSeconds: 0

---
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: backend
spec:
  minAvailable: 1
  selector:
    matchLabels:
      app: backend
```

## Ingress Patterns

### Nginx Ingress with TLS

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: app-ingress
  annotations:
    kubernetes.io/ingress.class: nginx
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    nginx.ingress.kubernetes.io/proxy-body-size: "50m"
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
    # Security headers
    nginx.ingress.kubernetes.io/configuration-snippet: |
      add_header X-Frame-Options "SAMEORIGIN" always;
      add_header X-Content-Type-Options "nosniff" always;
spec:
  ingressClassName: nginx
  tls:
    - hosts:
        - example.com
        - api.example.com
      secretName: app-tls
  rules:
    - host: example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: frontend
                port:
                  number: 80
    - host: api.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: backend
                port:
                  number: 8080
```

## Network Policies

```yaml
# Default deny all
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny
spec:
  podSelector: {}
  policyTypes:
    - Ingress
    - Egress

---
# Allow backend access
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-backend
spec:
  podSelector:
    matchLabels:
      app: backend
  policyTypes:
    - Ingress
    - Egress
  ingress:
    - from:
        - podSelector:
            matchLabels:
              app: frontend
      ports:
        - protocol: TCP
          port: 8080
  egress:
    - to:
        - podSelector:
            matchLabels:
              app: postgres
      ports:
        - protocol: TCP
          port: 5432
```

## Secrets Management

### Using Sealed Secrets

```bash
# Create secret
kubectl create secret generic backend-secrets \
  --from-literal=DATABASE_URL='postgresql://...' \
  --dry-run=client -o yaml | \
  kubeseal --format yaml > sealed-secret.yaml

# Apply sealed secret
kubectl apply -f sealed-secret.yaml
```

### Secret Reference in Deployment

```yaml
env:
  - name: DATABASE_URL
    valueFrom:
      secretKeyRef:
        name: backend-secrets
        key: DATABASE_URL
```

## Common Commands

```bash
# Preview manifests
kubectl kustomize k8s/overlays/staging

# Apply to cluster
kubectl apply -k k8s/overlays/staging

# Check status
kubectl get pods,svc,ingress -n staging

# View logs
kubectl logs -f deployment/backend -n staging

# Rollout status
kubectl rollout status deployment/backend -n staging

# Rollback
kubectl rollout undo deployment/backend -n staging

# Scale
kubectl scale deployment/backend --replicas=5 -n staging

# Port forward
kubectl port-forward svc/backend 8080:8080 -n staging
```

## Resource Guidelines

### Development/Staging

```yaml
resources:
  requests:
    cpu: "50m"
    memory: "128Mi"
  limits:
    cpu: "250m"
    memory: "256Mi"
```

### Production

```yaml
resources:
  requests:
    cpu: "200m"
    memory: "512Mi"
  limits:
    cpu: "1000m"
    memory: "1Gi"
```

## Security Checklist

- [ ] Pods run as non-root
- [ ] Read-only root filesystem
- [ ] Capabilities dropped
- [ ] Resource limits defined
- [ ] Network policies configured
- [ ] RBAC with least privilege
- [ ] Secrets encrypted (Sealed Secrets)
- [ ] Pod Security Standards enforced
- [ ] Ingress has security headers
- [ ] TLS enabled for all services
