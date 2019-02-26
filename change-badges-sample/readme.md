# Function to create badges based on build status (WIP)

Function expects `CREDENTIALS` env variable to create creds.json file that will be used for authentication. 

## First create a secret
kubectl create secret generic badges --from-literal=ENV=test --from-literal=BUCKET=yourbucketname --from-file=CREDENTIALS

CREDENTIALS should be the name of the file with your GOOGLE_APPLICATION_CREDENTIALS

## Second deploy a function with buildtemplate and env-secret name 

tm deploy service go-lambda -f . --build-template https://raw.githubusercontent.com/triggermesh/knative-lambda-runtime/master/go-1.x/buildtemplate.yaml --env-secret badges --wait
