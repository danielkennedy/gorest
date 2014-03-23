gorest
======
A simple HTTP server, written in Go, intended for Cloud Foundry deployment.

## Deploying to Cloud Foundry (public)

Note: I'm using my own MongoLab account for testing. Set the MONGO_URL appropriately for your environment.

1. `git clone git@github.com:danielkennedy/gorest`
1. `cf push <APPNAME> -b https://github.com/michaljemala/cloudfoundry-buildpack-go.git`
1. `cf set-env gorest MONGO_URL <MONGOLAB>`
1. `cf restart <APPNAME>`
