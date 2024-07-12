export ARTIFACTORY_AWS_ACCESS_KEY_ID :=
export ARTIFACTORY_AWS_SECRET_ACCESS_KEY :=
export ARTIFACTORY_AWS_REGION :=
export ARTIFACTORY_S3_BUCKET :=
export ARTIFACTORY_MAX_SIZE :=

all: build run

build:
	go build -o ./bin/upload-service .
	chmod +x ./bin/upload-service

run:
	./bin/upload-service

docker-build:
	docker build -t artifactory-upload:latest .

deploy:
	kubectl apply -k ./k8s