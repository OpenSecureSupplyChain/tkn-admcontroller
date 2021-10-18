# Tekton Admission Controller

tkn-admcontroller is an admission controller that checks and verifies sigstore style cryptographic signatures for
tekton pipeline / taskrun YAML files.

> :warning: Not ready for use yet!

tkn-admcontroller is still under active development, you're welcome to kick the tyres, but it's
advised not to use this until 1.0 is released.

## Install cert-manager
 
Install the cert-manager, by running the following command:

```bash
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.5.3/cert-manager.yaml
```

## Deploy Tekton

Install tekton, by running the following command:

```bash
kubectl apply --filename https://storage.googleapis.com/tekton-releases/pipeline/latest/release.yaml
```

## Deploy tekton admission controller

Once all services tekton and cert-manager are in a `Running` state (`kubectl get pods --all-namespaces`), we can
proceed to deploy the tekton admission controller.

Two approaches are possible here, you can use the existing image we have available, create an image yourself or build
with ko direct from your local source code (the latter being better for development workflows)

### Use the existing image

To use the local image, `demo/deployment.yaml' requires the following entry (this is already the default):

```yaml
spec:
  containers:
  - name: server
    image: ghcr.io/opensecuresupplychain/tkn-admission-controllers:0.0.2
```

### Build your own image from source

This will require setting the registry access credentials into `secrets.yaml`

Run the following commands:

   ```bash
   # replace <oci-registry> with your container registry url
   # replace <api-key> with your api key
   # replace <api-user> with your username
   # replace <email-address> with your email

   export oci_secret=$(kubectl create secret --dry-run=true -o yaml docker-registry registry-key --docker-server=<oci-registry> --docker-password=<api-key> --docker-username=<api-user> --docker-email=<email-address>)

   # assuming you're in the root directory
   echo $oci_secret > demo/secrets.yaml
   ```
Update the image reference in demo/deployment.yaml

```yaml
spec:
  containers:
  - name: server
    image: <your-image>

```
## Deploy tkn-admission-controller

Change directory to `demo` folder and run the `deploy.sh` script

   ```bash
   # change directory
   cd demo/
   # change permission for the deploy script
   chmod u+x deploy.sh
   # run the script
   ./deploy.sh
   ```

## Deploy using ko

Make sure you install ko and that it's in your `$PATH`.

Set up a local registry following the [steps outlined here](https://kind.sigs.k8s.io/docs/user/local-registry/).

Set the following environment variables

```bash
export KO_DOCKER_REPO="localhost:5000/mypipeline" 
```

Run `deploy.sh`, the script will sense ko is installed and deploy from `config/100-deployment.yaml`

## Sign and validate pipeline / taskrun manifest with cosign keys

1. Create a keypair with cosign

```bash
cosign generate-keys
```

2. Upload the public key as a secret

```bash
kubectl create secret generic cosign-pub --from-file=./cosign.pub`
```

3. Sign pipeline file with [k8s sigstore manifest tool](https://github.com/sigstore/k8s-manifest-sigstore)

```bash
kubectl-sigstore sign -k cosign.key -f my-manifest.yaml
```

* this will generate a signed file `my-manifest.yaml.signed`

4. Deploy the pipeline / taskrun

```bash
kubectl apply -f my-manifest.yaml.signed
```

Example

```bash
kubectl apply -f my-manifest.yaml.signed
pipeline.tekton.dev/tekton-pipeline configured

kubectl apply -f my-manifest.yaml
Error from server: error when creating "manifest.yaml": admission webhook "pipeline-validation.default.svc" denied the request: signature or message annotation not found


kubectl apply -f bad-manifest.yaml.signed
Error from server: error when creating "test-bad-sig.yaml": admission webhook "pipeline-validation.default.svc" denied the request: Signature validation failed

```
