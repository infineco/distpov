
Distributed POVRAY on google cloudrun
=====================================

The purpose of this demo is to show how simple it is to build & deploy a "simple" distributed process on [Cloud Run](https://cloud.google.com/run/). For this demo we will use [POV-Ray](https://github.com/POV-Ray/povray.git) - a raytracer - and without modifying one single byte of the original project, we will manage to run it on Cloud Run.

[![Run on Google Cloud](https://storage.googleapis.com/cloudrun/button.svg)](https://console.cloud.google.com/cloudshell/editor?shellonly=true&cloudshell_image=gcr.io/cloudrun/button&cloudshell_git_repo=https://github.com/infineco/distpov.git)

### Prerequisite :   
- [Locally installed docker](https://docs.docker.com/install/), not mandatory, but better for testing
- [CGP account](https://console.cloud.google.com/) 


### Setup :
1. [Select](https://console.cloud.google.com/projectselector2/home/dashboard) or [create](https://console.cloud.google.com/projectcreate) a GCP project.
2. Make sure that billing is enabled for your Google Cloud Platform project. 
[Learn](https://cloud.google.com/billing/docs/how-to/modify-project) how to enable billing. 
3. [Install and initialize](https://cloud.google.com/sdk/docs/) the Cloud SDK.
4. Install the gcloud beta component: `gcloud components install beta`  
5. Update components: `gcloud components update`


### Prepare GCloud
```
gcloud auth login  
PROJECT=povcloud
gcloud projects create $PROJECT --set-as-default
gcloud alpha billing accounts list
gcloud alpha billing projects link $PROJECT --billing-account 0X0X0X-0X0X0X-0X0X0X
```
### Build
```
docker build -t $PROJECT 
```

### Run Locally
```
docker run -t -i -p 8080:8080 $PROJECT
```

### Push the image to google image repository
```
PROJECT_ID=$(gcloud projects list | grep povcloud | awk '{print $1}')
docker tag povcloud eu.gcr.io/$PROJECT_ID/povcloud
docker push eu.gcr.io/$PROJECT_ID/pov
```

### Deploy to Cloud Run
deploy on [cloud run](https://console.cloud.google.com/run) using the interface. 
- choose the image you just uploaded 
- choose cloud run and select the nearest location 
- give a name to your service
- authorize non authenticated calls
- set 1 request per container 
- set 256 MB memory
- click on create and wait until your the service is created

or 
```
gcloud beta run deploy --image gcr.io/$PROJECT_ID/povcloud --platform managed
```

You will be prompted for the service name: press Enter to accept the default name.   
You will be prompted for region: select the region of your choice, for example us-central1.  
You will be prompted to allow unauthenticated invocations: respond y .  
Then wait a few moments until the deployment is complete. On success, the command line displays the service URL.  
  
Visit your deployed container by opening the service URL in a web browser.  

### Build a wrapper
