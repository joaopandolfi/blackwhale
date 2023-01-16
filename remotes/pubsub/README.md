# Pubsub abstraction
This is the pub/sub abstraction client to make more easy handle

## Setup
1. You need to setup the env GOOGLE_APPLICATION_CREDENTIALS with the json path
2. The json needs to contain the generated credentials on google cloud console

## Setup on k8s deployment
1. deployment.yaml
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  namespace: k8s-{ENVIRONMENT}
  name: test-service
  labels:
    name: test-service
spec:
  selector:
    matchLabels:
      name: test-service
  strategy:
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
    type: RollingUpdate
  template:
    metadata:
      labels:
        name: test-service
        app: test-service
    spec:
      volumes:
      - name: service-account-credentials-volume
        secret:
          secretName: GCP_CREDENTIAL_SECRET
          items:
          - key: sa_json
            path: sa_credentials.json
      containers:
        - image: gcr.io/plataform-xxxx/test-service:{VERSION}
          volumeMounts:
            - name: service-account-credentials-volume
              mountPath: /etc/gcp
              readOnly: true
          imagePullPolicy: Always
          name: test-service
          resources:
            requests:
              cpu: 100m
              memory: 100Mi
            limits:
              cpu: 100m
              memory: 256Mi
          env:
            - name: TZ
              value: America/Sao_Paulo
            - name: SNAKE_DEFAULT
              value: "false"
            - name: SYSTEM_ID
              value: "test"
            - name: BCRYPT_COST
              value: "12"
            - name: SERVER_PORT
              value: "5656"  
            - name: GOOGLE_APPLICATION_CREDENTIALS
              value: /etc/gcp/sa_credentials.json
          ports:
            - containerPort: 5656
              name: test-service
      restartPolicy: Always
---
apiVersion: v1
kind: Service
metadata:
  namespace: k8s-{ENVIRONMENT}
  name: test-service
  labels:
    app: test-service
spec:
  selector:
    name: test-service
  ports:
  - port: 80
    name: auth-ingress
    targetPort: 5656

```

2. The secret looks like
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: GCP_CREDENTIAL_SECRET
type: Opaque
data:
  sa_json: BASE64(JSON_CONTENT)

```