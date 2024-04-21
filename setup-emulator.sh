# stop on the first sign of failure
set -e

# create and configure a new gcloud configuration
gcloud config configurations create emulator-todos --no-activate || true
gcloud config configurations activate emulator-todos
gcloud config set auth/disable_credentials true --configuration=emulator-todos
gcloud config set project project-todos --configuration=emulator-todos
yes Y | gcloud config set api_endpoint_overrides/spanner http://localhost:9020/ --configuration=emulator-todos

# create a Spanner instance and database
gcloud spanner instances create instance-todos --config=emulator-config --description="Local Todo Instance" --nodes=1 --configuration=emulator-todos
gcloud spanner databases create database-todos --instance=instance-todos --configuration=emulator-todos

# update the database schema using DDL statements from a file
gcloud spanner databases ddl update database-todos --instance=instance-todos --ddl-file="./createDataBase.ddl" --configuration=emulator-todos

echo "Configuration 'emulator-todos' set up complete."