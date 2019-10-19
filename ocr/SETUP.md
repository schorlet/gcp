
## Setup

Install the gcloud beta component:

```sh
gcloud components install beta
```

Set default configuration:

```sh
GCP_PROJECT_ID=[YOUR_PROJECT_ID]
GCP_REGION_ID=[YOUR_REGION_ID]

gcloud config set project $GCP_PROJECT_ID
gcloud config set run/platform managed
gcloud config set run/region $GCP_REGION_ID
gcloud config set functions/region $GCP_REGION_ID
```

Enable the Google APIs:

```sh
gcloud services enable \
	cloudbuild.googleapis.com \
	cloudfunctions.googleapis.com \
	container.googleapis.com \
	containerregistry.googleapis.com \
	pubsub.googleapis.com \
	run.googleapis.com \
	storage.googleapis.com \
	vision.googleapis.com
```



### Service Accounts

The application use the princple of least-privilege allocation and each components (the webapp and the two functions) runs under a specified service account that allows the minimal permissions to perform the intented work.

In addition, the webapp uses a specific service account (`upload-sa`) to create a signed URL to allow the upload of an image directly to a specific storage bucket. As the [SignedURL](https://godoc.org/cloud.google.com/go/storage#SignedURL) needs a service account private key, this key is stored in a private bucket only accessible to the webapp service account.


#### The webapp service account

```sh
WEBAPP_USER=webapp-sa
WEBAPP_MEMBER=${WEBAPP_USER}@${GCP_PROJECT_ID}.iam.gserviceaccount.com

gcloud iam service-accounts create $WEBAPP_USER \
	--display-name="webapp runtime account"

```


#### The upload service account

```sh
UPLOAD_USER=upload-sa
UPLOAD_MEMBER=${UPLOAD_USER}@${GCP_PROJECT_ID}.iam.gserviceaccount.com

gcloud iam service-accounts create $UPLOAD_USER \
	--display-name="upload service account"
```


Generate the key file:

```sh
UPLOAD_FILE=upload-sa.json

gcloud iam service-accounts keys create $UPLOAD_FILE \
	--iam-account $UPLOAD_MEMBER \
	--key-file-type json

```

#### The process-image service account:

```sh
PROCESS_IMAGE_USER=process-image-sa
PROCESS_IMAGE_MEMBER=${PROCESS_IMAGE_USER}@${GCP_PROJECT_ID}.iam.gserviceaccount.com

gcloud iam service-accounts create $PROCESS_IMAGE_USER \
	--display-name="process-image function runtime account"

```


#### The process-text service account:

```sh
PROCESS_TEXT_USER=process-text-sa
PROCESS_TEXT_MEMBER=${PROCESS_TEXT_USER}@${GCP_PROJECT_ID}.iam.gserviceaccount.com

gcloud iam service-accounts create $PROCESS_TEXT_USER \
	--display-name="process-text function runtime account"
```






### Upload Bucket

This is the bucket where clients upload images.

```sh
UPLOAD_BUCKET=${GCP_PROJECT_ID}_upload
gsutil mb -c standard -l $GCP_REGION_ID gs://$UPLOAD_BUCKET
```

#### Grant permissions

Grant the `upload-sa` service account the `storage.objectCreator` role on the upload bucket:

```sh
gsutil iam ch \
	serviceAccount:${UPLOAD_MEMBER}:roles/storage.objectCreator \
	gs://$UPLOAD_BUCKET
```

Grant the `process-text-sa` service account the `storage.objectCreator` role on the upload bucket:

```sh
gsutil iam ch \
	serviceAccount:${PROCESS_TEXT_MEMBER}:roles/storage.objectCreator \
	gs://$UPLOAD_BUCKET
```



Grant the `process-image-sa` service account the `storage.legacyObjectReader` role on the upload bucket:

```sh
gsutil iam ch \
	serviceAccount:${PROCESS_IMAGE_MEMBER}:roles/storage.legacyObjectReader \
	gs://$UPLOAD_BUCKET
```

Grant the `webapp-sa` service account the `storage.legacyObjectReader` role on the upload bucket:

```sh
gsutil iam ch \
	serviceAccount:${WEBAPP_MEMBER}:roles/storage.legacyObjectReader \
	gs://$UPLOAD_BUCKET
```



#### Set CORS configuration

```json
tee cors.json <<EOF
[
  {
    "method": ["PUT"],
    "origin": ["http://localhost:8080"],
    "responseHeader": ["content-type", "x-goog-content-length-range", "x-goog-if-generation-match", "x-goog-storage-class"],
    "maxAgeSeconds": 3600
  },
  {
    "method": ["PUT"],
    "origin": ["https://webapp-[RANDOM_NUMBER]-ew.a.run.app"],
    "responseHeader": ["content-type", "x-goog-content-length-range", "x-goog-if-generation-match", "x-goog-storage-class"],
    "maxAgeSeconds": 3600
  }
]
EOF
```



### Creds Bucket

This is the private bucket where the service account key is stored.

```sh
CREDS_BUCKET=${GCP_PROJECT_ID}_creds
gsutil mb -c standard -l $GCP_REGION_ID gs://$CREDS_BUCKET
```

Revoke the default bucket permissions:

```sh
gsutil defacl set private gs://$CREDS_BUCKET
gsutil acl set private gs://$CREDS_BUCKET
```
>Only the bucket owner and explicitly granted users can access objects inside.


Copy the `upload-sa` service account key file:

```sh
gsutil -h 'Content-Type: application/json' \
	cp $UPLOAD_FILE gs://$CREDS_BUCKET/$UPLOAD_FILE
```


Grant the `webapp-sa` service account the `storage.legacyObjectReader` role on the `upload-sa` service account key file:

```sh
gsutil iam ch \
	serviceAccount:${WEBAPP_MEMBER}:roles/storage.legacyObjectReader \
	gs://$CREDS_BUCKET/$UPLOAD_FILE
```






### PubSub

Create a topic:

```sh
TEXT_TOPIC=text_topic
gcloud pubsub topics create $TEXT_TOPIC \
	--message-storage-policy-allowed-regions=$GCP_REGION_ID
```

Create the role `text_topic.publisher`:

```sh
gcloud iam roles create ${TEXT_TOPIC}.publisher \
	--project $GCP_PROJECT_ID \
	--description "Allow to publish on the $TEXT_TOPIC" \
	--permissions "pubsub.topics.publish,pubsub.topics.get"
```


Grant the `process-image-sa` service account the `text_topic.publisher` role on the topic:

```sh
gcloud pubsub topics add-iam-policy-binding $TEXT_TOPIC \
	--member serviceAccount:${PROCESS_IMAGE_MEMBER} \
	--role projects/${GCP_PROJECT_ID}/roles/${TEXT_TOPIC}.publisher
```




### Cloud Functions

Deploy the `process-image` function:

```sh
gcloud functions deploy process-image \
	--entry-point ProcessImage \
	--max-instances 1 \
	--memory 128MB \
	--runtime go111 \
	--service-account $PROCESS_IMAGE_MEMBER \
	--set-env-vars "TEXT_TOPIC=$TEXT_TOPIC" \
	--source process/image \
	--timeout 30s \
	--trigger-resource $UPLOAD_BUCKET \
	--trigger-event google.storage.object.finalize
```

Deploy the `process-text` function:

```sh
gcloud functions deploy process-text \
	--entry-point ProcessText \
	--max-instances 1 \
	--memory 128MB \
	--runtime go111 \
	--service-account $PROCESS_TEXT_MEMBER \
	--source process/text \
	--timeout 30s \
	--trigger-topic $TEXT_TOPIC
```






### Cloud Run


Build the Docker image using Cloud Build;

```sh
gcloud builds submit -t eu.gcr.io/${GCP_PROJECT_ID}/webapp webapp
```

Deploy the image:

```sh
gcloud beta run deploy webapp \
	--allow-unauthenticated \
	--image eu.gcr.io/${GCP_PROJECT_ID}/webapp \
	--memory 256M \
	--platform managed \
	--region $GCP_REGION_ID \
	--service-account $WEBAPP_MEMBER \
	--set-env-vars "UPLOAD_BUCKET=${UPLOAD_BUCKET}" \
	--set-env-vars "UPLOAD_CREDS=${CREDS_BUCKET}/${UPLOAD_FILE}" \
	--timeout 10s
```

To get the URL of the service:

```sh
WEBAPP_URL=$(gcloud beta run services describe webapp \
	--region $GCP_REGION_ID \
	--platform managed \
	--format "value(status.address.hostname)")
```

Use this URL to customize the CORS configuration of the upload bucket.

