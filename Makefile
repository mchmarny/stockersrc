.PHONY: test data-test run image k8s-deploy


test:
	go test ./... -v

data-test:
	go test data/*.go -v

run:
	go run *.go

mod:
	go mod tidy
	go mod vendor

image: mod
	gcloud builds submit \
		--project ${GCP_PROJECT} \
		--tag gcr.io/${GCP_PROJECT}/stocker:0.1.3

docker-image:
	docker build -t stocker .

docker-run:
	docker run -it stocker \
		-v /Users/mchmarny/.gcp-keys:/Users/mchmarny/.gcp-keys \
		-e GOOGLE_APPLICATION_CREDENTIALS=${GOOGLE_APPLICATION_CREDENTIALS} \
		-e T_CONSUMER_KEY=${T_CONSUMER_KEY} \
		-e T_CONSUMER_SECRET=${T_CONSUMER_SECRET} \
		-e T_ACCESS_TOKEN=${T_ACCESS_TOKEN} \
		-e T_ACCESS_SECRET=${T_ACCESS_SECRET} \
		-e STOCKER_CONF_PATH=${STOCKER_CONF_PATH}

k8s-setup:
	kubectl create secret generic stocker-secrets \
		--from-literal=T_CONSUMER_KEY=${T_CONSUMER_KEY} \
		--from-literal=T_CONSUMER_SECRET=${T_CONSUMER_SECRET} \
		--from-literal=T_ACCESS_TOKEN=${T_ACCESS_TOKEN} \
		--from-literal=T_ACCESS_SECRET=${T_ACCESS_SECRET}

	kubectl create secret generic stocker-sa-key \
		--from-file=key.json=/Users/mchmarny/.gcp-keys/s9-demo-key.json

	kubectl create configmap stocker-conf \
		--from-literal=STOCKER_CONF_PATH=https://raw.githubusercontent.com/mchmarny/stocker/master/config/dow.yaml

k8s-deploy:
	kubectl apply -f config/stocker.yaml

k8s-cleanup:
	kubectl delete secret stocker-secrets --ignore-not-found=true
	kubectl delete secret stocker-sa-key --ignore-not-found=true
	kubectl delete configmap stocker-conf --ignore-not-found=true
	kubectl delete -f config/stocker.yaml