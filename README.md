# Download Manager for Pivotal OpsMgr

This application automates the discovery of new releases for installed tiles, and their upload to a Pivotal OpsMgr. Administrators can then move quickly to installation/upgrade.

It also allows discovery of new product tiles and their initial upload.

## Local deployment

Setup environment variables for configuration:

```bash
export OPSMGR_URL=https://10.20.30.40
export OPSMGR_SKIP_SSL_VERIFICATION=true
export OPSMGR_USERNAME=admin
export OPSMGR_PASSWORD=<password>
export PIVOTAL_NETWORK_TOKEN=<token>
```

Next launch the app:

```
go run main.go
```

The application will eventually run at http://localhost:3000

## Cloud Foundry deployment

```bash
appname=opsmgr-download-mgr
cf push $appname -m 128M -k 2G --no-start
cf set-env $appname OPSMGR_URL $OPSMGR_URL
cf set-env $appname OPSMGR_SKIP_SSL_VERIFICATION $OPSMGR_SKIP_SSL_VERIFICATION
cf set-env $appname OPSMGR_USERNAME $OPSMGR_USERNAME
cf set-env $appname OPSMGR_PASSWORD $OPSMGR_PASSWORD
cf set-env $appname PIVOTAL_NETWORK_TOKEN $PIVOTAL_NETWORK_TOKEN
cf start $appname
```
