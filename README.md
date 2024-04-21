# Setup local environment

## 1. Start emulator using gCloud

1. Make sure Docker OR OrbStack is installed on your system
2. Install the [gcloud CLI](https://cloud.google.com/sdk/docs/install)
3. Update gcloud to get the latest version

```
gcloud components update
```

4. Start emulator

```
gcloud emulators spanner start
```

By default, the emulator hosts two local endpoints: localhost:9010 for gRPC requests and localhost:9020 for REST
requests.

Remember that closing the emulator will erase all data in it (e.g. database, tables).

## 2. Setup emulator configuration on gCloud CLI

To use the emulator with gcloud, you must disable authentication and override the endpoint.

While spanner emulator is running, run the following commands in new console:

```
gcloud config configurations create emulator-todos
gcloud config set auth/disable_credentials true
gcloud config set project project-todos
gcloud config set api_endpoint_overrides/spanner http://localhost:9020/
```

#### extra:

To switch between the emulator and default configuration, run:
`
gcloud config configurations activate [emulator | default]
`
To check all created configurations:
`
gcloud config configurations list
`

### Create on the emulator: Instance, DataBase and Tables

1. Create instance

```
gcloud spanner instances create instance-todos --config=emulator-config --description="Local Todo Instance" --nodes=1
```

2. Create database

```
gcloud spanner databases create database-todos --instance=instance-todos
```

3. Create tables - make sure you are in root dir of this project

```
gcloud spanner databases ddl update database-todos --instance=instance-todos --ddl-file="./createDataBase.ddl"
```

Remember that closing the emulator will erase the setup, and you will need to do it again after emulator restart

Instead of running all those commands for setting up the spanner emulator. You can just run the `setup-emulator.sh`
file.

#### extra:

`gcloud spanner instances list` - check all instances on this configuration
(if haven't created the instance yet, the output should be: Listed 0 items.)

`gcloud spanner databases list --instance=instance-todos` - check all databases on provided instance ("nicespanner")

`gcloud spanner databases execute-sql database-todos --instance=instance-todos --sql="SELECT table_name FROM information_schema.tables WHERE table_schema=''";` -
select tables name to verify if setup is successful

## 4. Start Local app and connect with local spanner

After having the spanner setup and running, you can run the backend server on your local machine.

1. Make sure [Go](https://go.dev/doc/install) is installed on your system.

2. Set variable so later when you run local backend server you use spanner emulator

```
export SPANNER_EMULATOR_HOST=localhost:9010
```

3. Install air (if you haven't done it earlier). `air` will automatically rebuild project after change.

```
go install github.com/cosmtrek/air@latest
```

3. Use this command to start local backend server.
```
air
```

#### extra:

To unset variable you can use: `export SPANNER_EMULATOR_HOST=""`. If you run the server it will try to connect with
remote spanner if `firebase-adminsdk.json` exists with correct keys.

## Summary

After following the instructions above you should have 2 running programs:

1. gCloud spanner emulator - this is your local Database - it will be cleaned after stopping the program.
2. Go backend server - API layer between client with database