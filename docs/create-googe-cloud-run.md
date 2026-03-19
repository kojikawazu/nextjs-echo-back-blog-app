

```bash

cd terraform-gcp
terraform init
terraform plan
terraform apply

cd ../
docker build --platform linux/amd64 -t asia-northeast1-docker.pkg.dev/portal-projects-449214/echo-blog-app-repo/echo-blog-app:latest .
docker push asia-northeast1-docker.pkg.dev/portal-projects-449214/echo-blog-app-repo/echo-blog-app:latest

cd terraform-gcp/
terraform plan
terraform apply

terraform import google_service_account.cloud_run_sa \
  projects/portal-projects-449214/serviceAccounts/cloud-run-sa@portal-projects-449214.iam.gserviceaccount.com
```