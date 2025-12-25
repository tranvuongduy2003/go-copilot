---
description: Create Kubernetes manifests for deploying a service with best practices
---

# Deploy to Kubernetes

Create Kubernetes manifests using Kustomize for deploying a service with production best practices.

## Service Details

**Service Name**: {{serviceName}}

**Service Type**: {{serviceType}}
- [ ] Web frontend
- [ ] API backend
- [ ] Background worker
- [ ] Database
- [ ] Cache

**Replicas**: {{replicas}}

**Port**: {{port}}

**Health Endpoint**: {{healthEndpoint}}

## Requirements

### Deployment

- [ ] Rolling update strategy
- [ ] Resource requests and limits
- [ ] Liveness and readiness probes
- [ ] Security context (non-root)
- [ ] Environment variables from ConfigMap/Secret

### Service

- [ ] ClusterIP service
- [ ] Port mapping
- [ ] Health check annotations

### Autoscaling

- [ ] Horizontal Pod Autoscaler
- [ ] CPU-based scaling
- [ ] Memory-based scaling
- [ ] Scale behavior configuration

### High Availability

- [ ] Pod Disruption Budget
- [ ] Pod anti-affinity
- [ ] Topology spread constraints

## Implementation

### 1. Deployment

Location: `k8s/base/{{serviceName}}/deployment.yaml`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{serviceName}}
  labels:
    app: {{serviceName}}
spec:
  replicas: {{replicas}}
  selector:
    matchLabels:
      app: {{serviceName}}
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  template:
    metadata:
      labels:
        app: {{serviceName}}
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "{{port}}"
    spec:
      serviceAccountName: {{serviceName}}
      securityContext:
        runAsNonRoot: true
        runAsUser: 1001
        fsGroup: 1001
      containers:
        - name: {{serviceName}}
          image: {{serviceName}}:latest
          ports:
            - name: http
              containerPort: {{port}}
          env:
            - name: NODE_ENV
              value: "production"
          envFrom:
            - configMapRef:
                name: {{serviceName}}-config
            - secretRef:
                name: {{serviceName}}-secrets
          resources:
            requests:
              cpu: "100m"
              memory: "256Mi"
            limits:
              cpu: "500m"
              memory: "512Mi"
          livenessProbe:
            httpGet:
              path: {{healthEndpoint}}
              port: http
            initialDelaySeconds: 30
            periodSeconds: 10
            timeoutSeconds: 5
            failureThreshold: 3
          readinessProbe:
            httpGet:
              path: {{healthEndpoint}}
              port: http
            initialDelaySeconds: 5
            periodSeconds: 5
            timeoutSeconds: 3
            failureThreshold: 3
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
                    app: {{serviceName}}
                topologyKey: kubernetes.io/hostname
```

### 2. Service

Location: `k8s/base/{{serviceName}}/service.yaml`

```yaml
apiVersion: v1
kind: Service
metadata:
  name: {{serviceName}}
spec:
  type: ClusterIP
  ports:
    - port: {{port}}
      targetPort: http
      name: http
  selector:
    app: {{serviceName}}

---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{serviceName}}

---
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{serviceName}}-config
data:
  LOG_LEVEL: "info"

---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: {{serviceName}}
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: {{serviceName}}
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
    scaleUp:
      stabilizationWindowSeconds: 0

---
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: {{serviceName}}
spec:
  minAvailable: 1
  selector:
    matchLabels:
      app: {{serviceName}}
```

### 3. Kustomization

Location: `k8s/base/{{serviceName}}/kustomization.yaml`

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
  - deployment.yaml
  - service.yaml
```

### 4. Add to Base Kustomization

Update: `k8s/base/kustomization.yaml`

```yaml
resources:
  - {{serviceName}}/deployment.yaml
  - {{serviceName}}/service.yaml
```

### 5. Environment Overlay (Staging)

Location: `k8s/overlays/staging/{{serviceName}}-patch.yaml`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{serviceName}}
spec:
  replicas: 1
  template:
    spec:
      containers:
        - name: {{serviceName}}
          resources:
            requests:
              cpu: "50m"
              memory: "128Mi"
            limits:
              cpu: "250m"
              memory: "256Mi"
```

## Deployment Commands

```bash
# Preview manifests
kubectl kustomize k8s/overlays/staging

# Apply to staging
kubectl apply -k k8s/overlays/staging

# Check status
kubectl get pods -l app={{serviceName}} -n staging
kubectl describe deployment {{serviceName}} -n staging

# View logs
kubectl logs -f deployment/{{serviceName}} -n staging

# Port forward for testing
kubectl port-forward svc/{{serviceName}} {{port}}:{{port}} -n staging

# Scale manually
kubectl scale deployment {{serviceName}} --replicas=3 -n staging

# Rollback
kubectl rollout undo deployment/{{serviceName}} -n staging
```

## Validation

After implementation:

1. Apply to staging namespace
2. Verify pods are running
3. Check health endpoints
4. Test autoscaling
5. Verify PDB works during rollout

## Output

Provide:
1. Deployment manifest
2. Service manifest
3. ConfigMap/Secret templates
4. HPA configuration
5. PDB configuration
6. Kustomization files
7. Deployment commands
