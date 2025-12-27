# Kubernetes Manifest Skill

Generate Kubernetes manifests for deploying applications.

## Usage

```
/project:skill:k8s <resource-type>
```

Types: `deployment`, `service`, `ingress`, `configmap`, `secret`, `full-stack`

## Namespace

**`k8s/namespace.yaml`**
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: app
  labels:
    app.kubernetes.io/name: app
    app.kubernetes.io/managed-by: kubectl
```

## Deployment

**`k8s/api/deployment.yaml`**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api
  namespace: app
  labels:
    app.kubernetes.io/name: api
    app.kubernetes.io/component: backend
spec:
  replicas: 3
  selector:
    matchLabels:
      app.kubernetes.io/name: api
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  template:
    metadata:
      labels:
        app.kubernetes.io/name: api
        app.kubernetes.io/component: backend
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "8080"
        prometheus.io/path: "/metrics"
    spec:
      serviceAccountName: api
      securityContext:
        runAsNonRoot: true
        runAsUser: 1000
        runAsGroup: 1000
        fsGroup: 1000
      containers:
        - name: api
          image: ghcr.io/yourorg/api:latest
          imagePullPolicy: Always
          ports:
            - name: http
              containerPort: 8080
              protocol: TCP
          env:
            - name: DATABASE_URL
              valueFrom:
                secretKeyRef:
                  name: api-secrets
                  key: database-url
            - name: REDIS_URL
              valueFrom:
                secretKeyRef:
                  name: api-secrets
                  key: redis-url
            - name: JWT_SECRET
              valueFrom:
                secretKeyRef:
                  name: api-secrets
                  key: jwt-secret
          envFrom:
            - configMapRef:
                name: api-config
          resources:
            requests:
              memory: "128Mi"
              cpu: "100m"
            limits:
              memory: "512Mi"
              cpu: "500m"
          livenessProbe:
            httpGet:
              path: /live
              port: http
            initialDelaySeconds: 10
            periodSeconds: 10
            timeoutSeconds: 5
            failureThreshold: 3
          readinessProbe:
            httpGet:
              path: /ready
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
                    app.kubernetes.io/name: api
                topologyKey: kubernetes.io/hostname
```

## Service

**`k8s/api/service.yaml`**
```yaml
apiVersion: v1
kind: Service
metadata:
  name: api
  namespace: app
  labels:
    app.kubernetes.io/name: api
    app.kubernetes.io/component: backend
spec:
  type: ClusterIP
  selector:
    app.kubernetes.io/name: api
  ports:
    - name: http
      port: 80
      targetPort: http
      protocol: TCP
```

## Ingress

**`k8s/api/ingress.yaml`**
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: api
  namespace: app
  labels:
    app.kubernetes.io/name: api
  annotations:
    kubernetes.io/ingress.class: nginx
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    nginx.ingress.kubernetes.io/proxy-body-size: "10m"
    nginx.ingress.kubernetes.io/rate-limit: "100"
    nginx.ingress.kubernetes.io/rate-limit-window: "1m"
spec:
  tls:
    - hosts:
        - api.example.com
      secretName: api-tls
  rules:
    - host: api.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: api
                port:
                  number: 80
```

## ConfigMap

**`k8s/api/configmap.yaml`**
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: api-config
  namespace: app
  labels:
    app.kubernetes.io/name: api
data:
  LOG_LEVEL: "info"
  LOG_FORMAT: "json"
  SERVER_PORT: "8080"
  SERVER_READ_TIMEOUT: "15s"
  SERVER_WRITE_TIMEOUT: "15s"
  CORS_ALLOWED_ORIGINS: "https://app.example.com"
```

## Secret

**`k8s/api/secret.yaml`**
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: api-secrets
  namespace: app
  labels:
    app.kubernetes.io/name: api
type: Opaque
stringData:
  database-url: "postgres://user:password@postgres:5432/app?sslmode=require"
  redis-url: "redis://:password@redis:6379"
  jwt-secret: "your-jwt-secret-here"
```

## ServiceAccount

**`k8s/api/serviceaccount.yaml`**
```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: api
  namespace: app
  labels:
    app.kubernetes.io/name: api
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: api
  namespace: app
rules:
  - apiGroups: [""]
    resources: ["configmaps", "secrets"]
    verbs: ["get", "list", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: api
  namespace: app
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: api
subjects:
  - kind: ServiceAccount
    name: api
    namespace: app
```

## HorizontalPodAutoscaler

**`k8s/api/hpa.yaml`**
```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: api
  namespace: app
  labels:
    app.kubernetes.io/name: api
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: api
  minReplicas: 3
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
```

## PodDisruptionBudget

**`k8s/api/pdb.yaml`**
```yaml
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: api
  namespace: app
  labels:
    app.kubernetes.io/name: api
spec:
  minAvailable: 2
  selector:
    matchLabels:
      app.kubernetes.io/name: api
```

## NetworkPolicy

**`k8s/api/networkpolicy.yaml`**
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: api
  namespace: app
  labels:
    app.kubernetes.io/name: api
spec:
  podSelector:
    matchLabels:
      app.kubernetes.io/name: api
  policyTypes:
    - Ingress
    - Egress
  ingress:
    - from:
        - namespaceSelector:
            matchLabels:
              kubernetes.io/metadata.name: ingress-nginx
      ports:
        - protocol: TCP
          port: 8080
  egress:
    - to:
        - namespaceSelector:
            matchLabels:
              kubernetes.io/metadata.name: app
        - podSelector:
            matchLabels:
              app.kubernetes.io/name: postgres
      ports:
        - protocol: TCP
          port: 5432
    - to:
        - namespaceSelector:
            matchLabels:
              kubernetes.io/metadata.name: app
        - podSelector:
            matchLabels:
              app.kubernetes.io/name: redis
      ports:
        - protocol: TCP
          port: 6379
    - to:
        - namespaceSelector: {}
          podSelector:
            matchLabels:
              k8s-app: kube-dns
      ports:
        - protocol: UDP
          port: 53
```

## Kustomization

**`k8s/kustomization.yaml`**
```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: app

resources:
  - namespace.yaml
  - api/serviceaccount.yaml
  - api/configmap.yaml
  - api/secret.yaml
  - api/deployment.yaml
  - api/service.yaml
  - api/ingress.yaml
  - api/hpa.yaml
  - api/pdb.yaml
  - api/networkpolicy.yaml

commonLabels:
  app.kubernetes.io/managed-by: kustomize
  app.kubernetes.io/part-of: app
```

## Commands

```bash
# Apply all manifests
kubectl apply -k k8s/

# Check deployment status
kubectl -n app get deployments

# Check pods
kubectl -n app get pods

# View logs
kubectl -n app logs -f deployment/api

# Scale deployment
kubectl -n app scale deployment/api --replicas=5

# Rollout status
kubectl -n app rollout status deployment/api

# Rollback
kubectl -n app rollout undo deployment/api
```

## Checklist

- [ ] Namespace created
- [ ] Deployment with proper resource limits
- [ ] Liveness and readiness probes
- [ ] Service for internal communication
- [ ] Ingress with TLS
- [ ] ConfigMap for non-sensitive config
- [ ] Secret for sensitive data
- [ ] ServiceAccount with minimal permissions
- [ ] HPA for autoscaling
- [ ] PDB for high availability
- [ ] NetworkPolicy for security
- [ ] Security context configured
