# About this repo

This is a proof of concept for a Golang gRPC API server, running in Google Kubernetes Engine, exposed through Cloud Endpoints with HTTP REST annotation, and auto-documented using Swagger UI. Follow the steps below and you will be able to deploy the configuration using Google Container Builder in your own Google Cloud Platform project.

# Instructions

```bash
# Clone the repo and enter it
git clone https://github.com/MichaelMarkieta/todo-gke-endpoints-grpc-rest-golang.git && cd todo-gke-endpoints-grpc-rest-golang

# Set up environmental variables
PROJECT_ID=Your-Project-ID
PROJECT_NUMBER=$(gcloud projects describe $(gcloud config get-value core/project) --format=value\(projectNumber\)) 
ZONE=Your-Compute-Zone
CLUSTER=Your-Cluster-Name

# Configure gcloud defaults
gcloud config set project $PROJECT_ID
gcloud config set compute/zone $ZONE

# Enable APIs
gcloud services enable container.googleapis.com
gcloud services enable cloudbuild.googleapis.com
gcloud services enable servicecontrol.googleapis.com
gcloud services enable servicemanagement.googleapis.com
gcloud services enable endpoints.googleapis.com

# Assign additional roles to the Google Container Builder service account
gcloud projects add-iam-policy-binding \
    $(gcloud config get-value core/project) \
    --member serviceAccount:$PROJECT_NUMBER@cloudbuild.gserviceaccount.com \
    --role roles/editor
gcloud projects add-iam-policy-binding \
    $(gcloud config get-value core/project) \
    --member serviceAccount:$PROJECT_NUMBER@cloudbuild.gserviceaccount.com \
    --role roles/container.developer
gcloud projects add-iam-policy-binding \
    $(gcloud config get-value core/project) \
    --member serviceAccount:$PROJECT_NUMBER@cloudbuild.gserviceaccount.com \
    --role roles/servicemanagement.admin

# Create a Google Kubernetes Engine Cluster
gcloud container clusters create $CLUSTER

# Deploy with Google Container Builder
gcloud container builds submit --config=cloudbuild.yaml --substitutions _ZONE=$ZONE,_CLUSTER=$CLUSTER .
```

# Swagger UI

```
HTTP_INGRESS_IP=$(kubectl get ingress ingress-guestbook \
    --output=jsonpath={.status.loadBalancer.ingress[0].ip})
echo http://$HTTP_INGRESS_IP/swagger/
```
