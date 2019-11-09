# cli-progress-animation

Example of progress animations and colour for a Golang CLI

## ASCII Cinema:

[![asciicast](https://asciinema.org/a/sJNDeN6VmI815UWNLvlohak57.svg)](https://asciinema.org/a/sJNDeN6VmI815UWNLvlohak57)

## Status:

This is a sample which could be extended or make more generic as part of a task runner.

It combines WaitGroup with github.com/morikuni/aec for colour.

I wrote this around 2 years ago when wanting to generate ASCII animations for the [OpenFaaS CLI](https://github.com/openfaas/faas-cli/) - it's not made it into the CLI yet, but if you're interested in contributing, then please let me know :-)

## Other work

Checkout:

* [My tech blog alexellis.io](https://blog.alexellis.io)

* [k3sup](https://github.com/alexellis/k3sup) - Installer for k3s and helm charts with strongly-typed parameters - works on any K8s cluster including RPi/ARM

* [inlets](https://inlets.dev) - HTTPS tunnels to your local network from anywhere

* [OpenFaaS](https://www.openfaas.com) - Serverless-style Functions & Microservices Made Simple on Kubernetes

## Testing

Clone the repo into your GOPATH and `go build`

```sh
git clone https://github.com/alexellis/cli-progress-animation
```
