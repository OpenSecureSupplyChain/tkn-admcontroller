# tkn-admcontroller



### Prerequisites and Testing on `minikube`
 
1. `cert-manager`: Install the cert-manager, by running the following command:

   ```bash
   kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.5.3/cert-manager.yaml
   ```

2. If you are not changing the code, you can use the sample public image `icr.io/gitsecure/sample-admission-controllers:0.0.2` for testing. If you are changing the code and wants to build a new image, you can follow these instructions:

a) Having the secrets to pull docker images in your cluster. When using the `us.icr.io` container registry, the secret obtained by running the following command is required is required.

   ```bash
   # replace :api-key with your api key
   # replace :email-address with your email

   export icr_secret=$(kubectl create secret --dry-run=true -o yaml docker-registry icr-registry-key --docker-server=us.icr.io --docker-password=<api-key> --docker-username=iamapikey --docker-email=<email-address>)

   # assuming you're in the root directory
   echo $icr_secret > demo/secrets.yaml
   ```
b) Update the image reference in demo/deployment.yaml 

3. Change directory to `demo` folder and run the following command.

   ```bash
   # change directory
   cd demo/
   # change permission for the deploy script
   chmod u+x deploy.sh
   # run the script
   ./deploy.sh
   ```

4. Test for the running resources. We're using the default namespace for demo purposes, assuming you're running minikube.

   ```bash
   # still in the `demo` folder
   kubectl get all
   # you can log for the admission-server pod

   # testing fo the pipelines
   # create the tasks
   kubectl apply -f pipelines/01_tasks.yaml

   # test for failing pipeline
   kubectl apply -f pipelines/02_fail_pipeline_creation.yaml
   # this fails with the following error
   # Error from server: error when creating "./pipelines/02_fail_pipeline_creation.yaml": admission webhook "pipeline-validation.default.svc" denied the request: sigstore sign annotation not found

   # now test for success pipeline
   kubectl apply -f pipelines/03_success_pipeline_creation.yaml

   ```

5. That is it for testing :)

## Pipeline

Install tekton

```bash
kubectl apply --filename https://storage.googleapis.com/tekton-releases/pipeline/latest/release.yaml
```