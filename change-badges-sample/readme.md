# Function to create badges based on build status (WIP)

This code was inspired by [the blog post](https://ljvmiranda921.github.io/notebook/2018/12/21/cloud-build-badge/)

This function

- Consumes build notifications via GCPPubSub source
- Processes GCPPubSub event 
- Creates the .svc file with the badge based on build status (SUCCESS, FAILURE) from the event

Function expects `CREDENTIALS` env variable passed through the secret to set `GOOGLE_APPLICATION_CREDENTIALS` env variable for GCP authentication and `BUCKET` to connect to selected GCP bucket 

Function connects to your bucket, reads contect of `/buid` folder expecting `failure.svg` and `success.svg` files and creates badges based on those files with the following naming convention: `/build/repoName-branchName.svg`

## Function Deploy

1. Create a valid secret ``` kubectl create secret generic badges --from-literal=BUCKET=yourbucketname --from-file=CREDENTIALS=credentials.json ``` credentials.json should be the name of the file with your service account credentials that has access to your bucket

2. Create folder `/build` in your bucket and upload two badges there (`failure.svg`, `success.svg`)

3. Deploy the function with buildtemplate and env-secret name ```tm deploy service go-lambda -f . --build-template https://raw.githubusercontent.com/triggermesh/knative-lambda-runtime/master/go-1.x/buildtemplate.yaml --env-secret badges --wait ```


4. Reference the badges in your README.md files! 