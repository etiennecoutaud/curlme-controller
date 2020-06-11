# CurlMe Controller

![CurlMeController](https://github.com/etiennecoutaud/curlme-controller/workflows/CurlMeController/badge.svg)
![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/etiennecoutaud/curlme-controller)
[![codecov](https://codecov.io/gh/etiennecoutaud/curlme-controller/branch/master/graph/badge.svg)](https://codecov.io/gh/etiennecoutaud/curlme-controller)
[![Go Report Card](https://goreportcard.com/badge/github.com/etiennecoutaud/curlme-controller)](https://goreportcard.com/report/github.com/etiennecoutaud/curlme-controller)

CurlMe controller is a simple Kubernetes controller that watch `ConfigMap` with `x-k8s.io/curl-me-that` annotation.

## Overview

CurMe controller watch `ConfigMap` ressource like below
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: example
  annotations:
    x-k8s.io/curl-me-that: joke=curl-a-joke.herokuapp.com
data: {}
```

`curl-a-joke.herokuapp.com` will be request and response body store into the configmap under `joke` key.

If request fail, a warning event will be raised

## Quickstart

To deploy the controller into your kubernetes cluster :

```bash
$ kubectl apply -f https://raw.githubusercontent.com/etiennecoutaud/curlme-controller/master/manifests/all-in-one.yaml
```

## Run locally

To run the app localy make sure you have `go 1.13` install and `KUBECONFIG` environment var set.

```bash
$ make run-local
```