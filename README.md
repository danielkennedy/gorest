gorest
======
A simple HTTP server, written in Go, intended for Cloud Foundry deployment.

## Deploying to Cloud Foundry (public)

1. `git clone git@github.com:danielkennedy/gorest`
1. `cf push APPNAME -b https://github.com/michaljemala/cloudfoundry-buildpack-go.git`
