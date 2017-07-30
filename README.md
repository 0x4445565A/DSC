# DSC SDEI Project

## Build Run Test
```
git clone git@github.com:0x4445565A/DSC.git
cd DSC

# Start Minikube
minikube start

# Switch to minikube scope
eval $(minikube docker-env)

# build dsc-logger container
# This also runs test
docker build -t dsc:v1 .

# build nginx-app this would normally be an individual project normally but for this example we'll leave it here
cd k8s/nginx
docker build -t nginx-app:v1 .

# build pods
cd ..
kubectl create -f dsc-project.yml

# View the pod is up and running
kubectl get pods

# Expose service
kubectl expose pod dsc-project --port=80 --name=dsc-project --type=NodePort

# Test endpoints
CLUSTER_URL=$(minikube service dsc-project --url)
curl $CLUSTER_URL
curl $CLUSTER_URL/300
curl $CLUSTER_URL/404
curl $CLUSTER_URL/500
curl $CLUSTER_URL/600

# View logs
kubectl logs dsc-project dsc-logger
```


## Clean Up
```
kubectl delete service dsc-project
kubectl delete pod dsc-project
minikube stop
```
