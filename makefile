SERVICE_NAME=code-sandbox
SERVICE_PORT=6758

run:
	@echo "running docker image for $(SERVICE_NAME)"
	@docker run -d --privileged --name coder -p $(SERVICE_PORT):$(SERVICE_PORT)  $(SERVICE_NAME)

build:
	@echo "building docker image for $(SERVICE_NAME)"
	@docker build --build-arg=SERVICE_PORT=$(SERVICE_PORT) \
	-t $(SERVICE_NAME) .
