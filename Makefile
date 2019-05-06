.PHONY: test data-test run image k8s-deploy

test:
	go test ./... -v

run:
	go run *.go

mod:
	go mod tidy
	go mod vendor

image: mod
	gcloud builds submit \
		--project ${GCP_PROJECT} \
		--tag gcr.io/${GCP_PROJECT}/stockersrc:0.1.4