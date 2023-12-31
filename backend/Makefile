# Makefile

# Set project ID and image name
GCP_PROJECT_ID=growth-squad-396607
BOT_IMAGE=wallemon-bot
SERVICE_IMAGE=wallemon-service
LATEST_TAG=latest
TESTNET_TAG=testnet
GCR_HOSTNAME=asia.gcr.io

.PHONY: docker-compose push-bot push-service configure-docker deploy build-bot run-bot

run-bot:
	docker run -it --rm --network host --name $(BOT_IMAGE) $(GCR_HOSTNAME)/$(GCP_PROJECT_ID)/$(BOT_IMAGE):$(LATEST_TAG)
run-service:
	docker run -it --rm --network host --name $(SERVICE_IMAGE) $(GCR_HOSTNAME)/$(GCP_PROJECT_ID)/$(SERVICE_IMAGE):$(LATEST_TAG)
run-db:
	docker run -it --rm --name postgres -e POSTGRES_PASSWORD=user -e POSTGRES_USER=user -p 5432:5432 postgres:14.9-alpine


docker-compose:
	docker-compose -f docker-compose.yaml up

build-bot:
	@if [ "$(env)" = "local" ]; then \
		echo "Build image for running locally"; \
	else \
		echo "Build image for running on GCP"; \
	fi
	$(eval TAG=$(if $(filter $(testnet),true),$(TESTNET_TAG),$(LATEST_TAG)))
	docker build --build-arg ENV=$(env) -f dockerfile-bot -t $(BOT_IMAGE) .
	docker tag $(BOT_IMAGE) $(GCR_HOSTNAME)/$(GCP_PROJECT_ID)/$(BOT_IMAGE):$(TAG)

build-service:
	@if [ "$(env)" = "local" ]; then \
		echo "Build image for running locally"; \
	else \
		echo "Build image for running on GCP"; \
	fi
	$(eval TAG=$(if $(filter $(testnet),true),$(TESTNET_TAG),$(LATEST_TAG)))
	@echo "TAG is: $(TAG)"
	docker build --build-arg ENV=$(env) -f dockerfile-service -t $(SERVICE_IMAGE) .
	docker tag $(SERVICE_IMAGE) $(GCR_HOSTNAME)/$(GCP_PROJECT_ID)/$(SERVICE_IMAGE):$(TAG)


# Push the Docker image to Google Container Registry
push-service: configure-docker build-service
	@echo "Pushing Service Docker image to GCR..."
	$(eval TAG=$(if $(filter $(testnet),true),$(TESTNET_TAG),$(LATEST_TAG)))
	@echo "TAG is: $(TAG)"
	docker tag $(SERVICE_IMAGE) $(GCR_HOSTNAME)/$(GCP_PROJECT_ID)/$(SERVICE_IMAGE):$(TAG)
	docker push $(GCR_HOSTNAME)/$(GCP_PROJECT_ID)/$(SERVICE_IMAGE):$(TAG)

push-bot: configure-docker build-bot
	@echo "Pushing Bot Docker image to GCR..."
	$(eval TAG=$(if $(filter $(testnet),true),$(TESTNET_TAG),$(LATEST_TAG)))
	@echo "TAG is: $(TAG)"
	docker tag $(BOT_IMAGE) $(GCR_HOSTNAME)/$(GCP_PROJECT_ID)/$(BOT_IMAGE):$(TAG)
	docker push $(GCR_HOSTNAME)/$(GCP_PROJECT_ID)/$(BOT_IMAGE):$(TAG)

# Authenticate Docker to GCR
configure-docker:
	@echo "Configuring Docker for GCR..."
	gcloud auth configure-docker $(GCR_HOSTNAME)

# Deploy the Docker image to a GCP VM
deploy: push
	@echo "Deploying Docker image to GCP VM..."
	ssh $(VM_USER)@$(VM_IP) "\
		gcloud auth configure-docker $(GCR_HOSTNAME) && \
		docker pull $(GCR_HOSTNAME)/$(GCP_PROJECT_ID)/$(BOT_IMAGE):$(LATEST_TAG) && \
		docker stop $(BOT_IMAGE) || true && \
		docker rm $(BOT_IMAGE) || true && \
		docker run -d --name $(BOT_IMAGE) -p 8080:8080 $(GCR_HOSTNAME)/$(GCP_PROJECT_ID)/$(BOT_IMAGE):$(LATEST_TAG)"
