# Dialogflow CX Call Transfer Router
The provided webhook provides call transfer routing based on a provided entity name, routing the caller to the that entity's phone number.  

# Usage
## Requirements
- Make sure you have sufficient IAM permissions to: 
  1. Edit Dialogflow CX Agent definitions (make and assign webhooks)
  2. Create and Store container images in the Container Registry
  3. Define and deploy Cloud Run services.
  4. The Cloudbuild (cloudbuild.googleapis.com) must be enabled and configured to define and create Cloud Run Services 
- It is recommended to host the Cloud Run service that hosts the webhook in the same project where the Dialogflow Virtual Agent exists.  
- Grant the Dialogflow Service Agent (found under IAM & Admin > IAM, you must turn Enable Google Provided role grants to see this service account) the Cloud Run invoker role: roles/run.invoker.  

## Instructions
Start by cloning a copy of this repository and switching directories:

```shell
git clone https://github.com/YvanJAquino/cx-call-transfer-router.git
cd cx-call-transfer-router.git
```

Review the provided cloudbuild.yaml's Cloud Run configuration (see step id:`gcloud-run-deploy-cx-call-transfer-router`).  Once reviewed, run `gcloud builds submit` from Cloud Shell.  This will create and store the container within Google Cloud's container registry and then create the Cloud Run service that hosts the Webhook.  

Once the Cloud Run service is ready, copy the provided URL, and append the Handler's path (`/route`) to it: 

- `https://auto-generated-spam.run.app/route`

This is the fully qualified URL to use as the webhook URL inside of the virtual agent.  

```yaml
steps:
- id: docker-build-push-cx-call-transfer-router
  waitFor: ['-']
  name: gcr.io/cloud-builders/docker
  dir: service
  entrypoint: bash
  args:
    - -c
    - |
      docker build -t gcr.io/$PROJECT_ID/${_SERVICE} . &&
      docker push gcr.io/$PROJECT_ID/${_SERVICE}

- id: gcloud-run-deploy-cx-call-transfer-router
  waitFor: ['docker-build-push-cx-wi-pubdef']
  name: gcr.io/google.com/cloudsdktool/cloud-sdk
  entrypoint: bash
  # REVIEW THE FOLLOWING SETTINGS !
  args: 
    - -c
    - |
      gcloud run deploy ${_SERVICE} \
        --project $PROJECT_ID \
        --image gcr.io/$PROJECT_ID/${_SERVICE} \
        --set-env-vars PROJECT_ID=$PROJECT_ID \
        --set-env-vars COLLECTION=${_COLLECTION} \
        --set-env-vars DOCUMENT=${_DOCUMENT} \
        --timeout 5m \
        --region ${_REGION} \
        --no-cpu-throttling \
        --min-instances 0 \
        --max-instances 5 \
        --allow-unauthenticated

substitutions:
  _SERVICE: cx-call-transfer-router
  _REGION: us-central1
  _COLLECTION: lawyers
  _DOCUMENT: lines
```

# As-Is Disclaimer
Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.