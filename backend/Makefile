# Makefile

# Set project ID and image name
GCP_PROJECT_ID=growth-squad-396607
IMAGE_NAME=wallemon
IMAGE_TAG=latest
GCR_HOSTNAME=asia.gcr.io

# Create a Docker image
.PHONY: build
build:
	@echo "Building Docker image..."
	docker build -t $(IMAGE_NAME) .

# Push the Docker image to Google Container Registry
.PHONY: push
push: configure-docker build
	@echo "Pushing Docker image to GCR..."
	docker tag $(IMAGE_NAME) $(GCR_HOSTNAME)/$(GCP_PROJECT_ID)/$(IMAGE_NAME):$(IMAGE_TAG)
	docker push $(GCR_HOSTNAME)/$(GCP_PROJECT_ID)/$(IMAGE_NAME):$(IMAGE_TAG)

# Authenticate Docker to GCR
.PHONY: configure-docker
configure-docker:
	@echo "Configuring Docker for GCR..."
	gcloud auth configure-docker $(GCR_HOSTNAME)

# Deploy the Docker image to a GCP VM
.PHONY: deploy
deploy: push
	@echo "Deploying Docker image to GCP VM..."
	ssh $(VM_USER)@$(VM_IP) "\
		gcloud auth configure-docker $(GCR_HOSTNAME) && \
		docker pull $(GCR_HOSTNAME)/$(GCP_PROJECT_ID)/$(IMAGE_NAME):$(IMAGE_TAG) && \
		docker stop $(IMAGE_NAME) || true && \
		docker rm $(IMAGE_NAME) || true && \
		docker run -d --name $(IMAGE_NAME) -p 8080:8080 $(GCR_HOSTNAME)/$(GCP_PROJECT_ID)/$(IMAGE_NAME):$(IMAGE_TAG)"