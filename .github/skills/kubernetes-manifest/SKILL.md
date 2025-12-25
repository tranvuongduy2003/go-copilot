---
name: kubernetes-manifest
description: Create Kubernetes manifests with Kustomize for multi-environment deployments
---

# Kubernetes Manifest Skill

This skill guides you through creating Kubernetes manifests using Kustomize for production-ready, multi-environment deployments.

## Tech Stack

This project uses:
- **Backend**: Go 1.25+ with Chi router (Clean Architecture + DDD + CQRS)
- **Frontend**: React 19 with pnpm (Next.js standalone output)
- **Database**: PostgreSQL 16
- **Cache**: Redis
- **Orchestration**: Kubernetes with Kustomize
- **Container Registry**: ghcr.io

## When to Use This Skill

- Deploying new services to Kubernetes
- Creating multi-environment configurations
- Setting up autoscaling and high availability
- Configuring network policies and security

## Project Structure

```
k8s/
├── base/                            # Base manifests
│   ├── kustomization.yaml
│   ├── namespace.yaml
│   ├── ingress.yaml
│   ├── frontend/
│   │   ├── deployment.yaml
│   │   └── service.yaml
│   └── backend/
│       ├── deployment.yaml
│       └── service.yaml
└── overlays/
    ├── staging/
    │   ├── kustomization.yaml
    │   └── namespace.yaml
    └── production/
        ├── kustomization.yaml
        ├── namespace.yaml
        └── sealed-secrets.yaml
```

## Templates

### Template 1: Deployment

```yaml
# k8s/base/<service>/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: SERVICE_NAME
  labels:
    app: SERVICE_NAME
    tier: TIER
spec:
  replicas: 2
  selector:
    matchLabels:
      app: SERVICE_NAME
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  template:
    metadata:
      labels:
        app: SERVICE_NAME
        tier: TIER
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "PORT"
        prometheus.io/path: "/metrics"
    spec:
      serviceAccountName: SERVICE_NAME
      securityContext:
        runAsNonRoot: true
        runAsUser: 1001
        fsGroup: 1001
      terminationGracePeriodSeconds: 30
      containers:
        - name: SERVICE_NAME
          image: SERVICE_NAME:latest
          imagePullPolicy: Always
          ports:
            - name: http
              containerPort: PORT
              protocol: TCP
          env:
            - name: NODE_ENV
              value: "production"
            - name: PORT
              value: "PORT"
          envFrom:
            - configMapRef:
                name: SERVICE_NAME-config
            - secretRef:
                name: SERVICE_NAME-secrets
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
            timeoutSeconds: 5
            failureThreshold: 3
          readinessProbe:
            httpGet:
              path: /health/ready
              port: http
            initialDelaySeconds: 5
            periodSeconds: 5
            timeoutSeconds: 3
            failureThreshold: 3
          securityContext:
            allowPrivilegeEscalation: false
            readOnlyRootFilesystem: true
            capabilities:
              drop:
                - ALL
          volumeMounts:
            - name: tmp
              mountPath: /tmp
            - name: cache
              mountPath: /app/.cache
      volumes:
        - name: tmp
          emptyDir: {}
        - name: cache
          emptyDir: {}
      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
            - weight: 100
              podAffinityTerm:
                labelSelector:
                  matchLabels:
                    app: SERVICE_NAME
                topologyKey: kubernetes.io/hostname
      topologySpreadConstraints:
        - maxSkew: 1
          topologyKey: topology.kubernetes.io/zone
          whenUnsatisfiable: ScheduleAnyway
          labelSelector:
            matchLabels:
              app: SERVICE_NAME
```

### Template 2: Service with HPA and PDB

```yaml
# k8s/base/<service>/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: SERVICE_NAME
  labels:
    app: SERVICE_NAME
spec:
  type: ClusterIP
  ports:
    - port: PORT
      targetPort: http
      protocol: TCP
      name: http
  selector:
    app: SERVICE_NAME

---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: SERVICE_NAME
  labels:
    app: SERVICE_NAME

---
apiVersion: v1
kind: ConfigMap
metadata:
  name: SERVICE_NAME-config
  labels:
    app: SERVICE_NAME
data:
  LOG_LEVEL: "info"
  LOG_FORMAT: "json"

---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: SERVICE_NAME
  labels:
    app: SERVICE_NAME
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: SERVICE_NAME
  minReplicas: 2
  maxReplicas: 10
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
    - type: Resource
      resource:
        name: memory
        target:
          type: Utilization
          averageUtilization: 80
  behavior:
    scaleDown:
      stabilizationWindowSeconds: 300
      policies:
        - type: Percent
          value: 10
          periodSeconds: 60
    scaleUp:
      stabilizationWindowSeconds: 0
      policies:
        - type: Percent
          value: 100
          periodSeconds: 15
        - type: Pods
          value: 4
          periodSeconds: 15
      selectPolicy: Max

---
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: SERVICE_NAME
  labels:
    app: SERVICE_NAME
spec:
  minAvailable: 1
  selector:
    matchLabels:
      app: SERVICE_NAME
```

### Template 3: Ingress with TLS

```yaml
# k8s/base/ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: app-ingress
  annotations:
    kubernetes.io/ingress.class: nginx
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    nginx.ingress.kubernetes.io/proxy-body-size: "50m"
    nginx.ingress.kubernetes.io/proxy-read-timeout: "300"
    nginx.ingress.kubernetes.io/limit-rps: "100"
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
    nginx.ingress.kubernetes.io/configuration-snippet: |
      add_header X-Frame-Options "SAMEORIGIN" always;
      add_header X-Content-Type-Options "nosniff" always;
      add_header X-XSS-Protection "1; mode=block" always;
spec:
  ingressClassName: nginx
  tls:
    - hosts:
        - example.com
        - api.example.com
      secretName: app-tls-secret
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
                  number: 3000
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

### Template 4: Network Policy

```yaml
# k8s/base/network-policy.yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny-all
spec:
  podSelector: {}
  policyTypes:
    - Ingress
    - Egress

---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-SERVICE_NAME
spec:
  podSelector:
    matchLabels:
      app: SERVICE_NAME
  policyTypes:
    - Ingress
    - Egress
  ingress:
    - from:
        - namespaceSelector:
            matchLabels:
              name: ingress-nginx
        - podSelector:
            matchLabels:
              app: frontend
      ports:
        - protocol: TCP
          port: PORT
  egress:
    - to:
        - podSelector:
            matchLabels:
              app: postgres
      ports:
        - protocol: TCP
          port: 5432
    - to:
        - namespaceSelector: {}
          podSelector:
            matchLabels:
              k8s-app: kube-dns
      ports:
        - protocol: UDP
          port: 53
```

### Template 5: Base Kustomization

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

### Template 6: Staging Overlay

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
  # Reduce replicas
  - patch: |-
      - op: replace
        path: /spec/replicas
        value: 1
    target:
      kind: Deployment

  # Reduce resources
  - patch: |-
      - op: replace
        path: /spec/template/spec/containers/0/resources/requests/cpu
        value: "50m"
      - op: replace
        path: /spec/template/spec/containers/0/resources/requests/memory
        value: "128Mi"
      - op: replace
        path: /spec/template/spec/containers/0/resources/limits/cpu
        value: "250m"
      - op: replace
        path: /spec/template/spec/containers/0/resources/limits/memory
        value: "256Mi"
    target:
      kind: Deployment

  # Reduce HPA
  - patch: |-
      - op: replace
        path: /spec/minReplicas
        value: 1
      - op: replace
        path: /spec/maxReplicas
        value: 3
    target:
      kind: HorizontalPodAutoscaler

images:
  - name: frontend
    newName: ghcr.io/myorg/app/frontend
    newTag: staging
  - name: backend
    newName: ghcr.io/myorg/app/backend
    newTag: staging

configMapGenerator:
  - name: frontend-config
    behavior: merge
    literals:
      - API_URL=https://api.staging.example.com
  - name: backend-config
    behavior: merge
    literals:
      - LOG_LEVEL=debug
```

### Template 7: Production Overlay

```yaml
# k8s/overlays/production/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: production
namePrefix: prod-

commonLabels:
  environment: production

resources:
  - ../../base
  - namespace.yaml
  - sealed-secrets.yaml

patches:
  # Production replicas
  - patch: |-
      - op: replace
        path: /spec/replicas
        value: 3
    target:
      kind: Deployment
      name: backend

  # Production resources
  - patch: |-
      - op: replace
        path: /spec/template/spec/containers/0/resources/requests/cpu
        value: "200m"
      - op: replace
        path: /spec/template/spec/containers/0/resources/requests/memory
        value: "512Mi"
      - op: replace
        path: /spec/template/spec/containers/0/resources/limits/cpu
        value: "1000m"
      - op: replace
        path: /spec/template/spec/containers/0/resources/limits/memory
        value: "1Gi"
    target:
      kind: Deployment
      name: backend

  # Production HPA
  - patch: |-
      - op: replace
        path: /spec/minReplicas
        value: 3
      - op: replace
        path: /spec/maxReplicas
        value: 20
    target:
      kind: HorizontalPodAutoscaler
      name: backend

images:
  - name: frontend
    newName: ghcr.io/myorg/app/frontend
    newTag: latest
  - name: backend
    newName: ghcr.io/myorg/app/backend
    newTag: latest
```

## Commands

```bash
# Preview manifests
kubectl kustomize k8s/overlays/staging

# Apply to cluster
kubectl apply -k k8s/overlays/staging

# Diff changes
kubectl diff -k k8s/overlays/staging

# Delete resources
kubectl delete -k k8s/overlays/staging

# Check status
kubectl get pods,svc,ingress -n staging

# Rollout status
kubectl rollout status deployment/backend -n staging

# Rollback
kubectl rollout undo deployment/backend -n staging
```

## Checklist

- [ ] Deployment has rolling update strategy
- [ ] Resource requests and limits set
- [ ] Liveness and readiness probes configured
- [ ] Security context with non-root user
- [ ] HPA configured for autoscaling
- [ ] PDB configured for availability
- [ ] Network policies restrict traffic
- [ ] Secrets managed securely
- [ ] ConfigMaps for configuration
- [ ] Ingress with TLS and security headers
