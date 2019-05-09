# stockersrc

Twitter data source

## Setup

Creating PubSub topic which will be used to publish tweet stream

```shell
gcloud pubsub topics create stocker-source
```

Create BigQuery company table

```shell
bq query --use_legacy_sql=false "
  CREATE OR REPLACE TABLE stocker.company (
    symbol STRING NOT NULL,
    aliases STRING NOT NULL
)"
```

Load companies into the company table

> Edit `bin/stocks.csv` before loading if you want to change the stocks to track. Just remember, tracking too many companies will lead to twitter API throttling

```shell
bq --location=US load --source_format=CSV stocker.company ./bin/stocks.csv
```

Publish your own image of this twitter source

```shell
go mod tidy
go mod vendor

gcloud builds submit \
	--project ${GCP_PROJECT} \
	--tag gcr.io/${GCP_PROJECT}/stockersrc:0.1.4
```

## Deploy

> Note, you will need Twitter API consumer and access keys along with their secrets


Create source vm

```shell
gcloud compute instances create-with-container stockersrc-vm \
       --container-image "gcr.io/s9-demo/stockersrc:0.1.4" \
       --machine-type n1-standard-2 \
       --zone us-west1-c \
       --image-family=cos-stable \
       --image-project=cos-cloud \
       --maintenance-policy MIGRATE \
       --container-restart-policy=always \
       --scopes "cloud-platform" \
       --container-privileged \
       --container-env="GCP_PROJECT=${GCP_PROJECT},T_CONSUMER_KEY=${T_CONSUMER_KEY},T_CONSUMER_SECRET=${T_CONSUMER_SECRET},T_ACCESS_TOKEN=${T_ACCESS_TOKEN},T_ACCESS_SECRET=${T_ACCESS_SECRET}"
```

To trace the logs from above deployed container, first capture the VM's ID

```shell
INSTANCE_ID=$(gcloud compute instances describe stockersrc-vm --zone us-west1-c --format="value(id)")
```

```shell
gcloud logging read "resource.type=gce_instance AND \
    logName=projects/${GCP_PROJECT}/logs/cos_containers AND \
    resource.labels.instance_id=${INSTANCE_ID}"
```

## Cleanup

Delete content load dataflow job

```shell
gcloud dataflow jobs cancel stocker-stocker-source-content --region=us-central1
```

Delete source vm

```shell
gcloud compute instances delete stockersrc-vm --zone=us-west1-c
```

Delete source topics

```shell
gcloud beta pubsub topics delete stocker-source
```
