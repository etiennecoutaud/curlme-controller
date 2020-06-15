# CurlMe Controller

[![CurlMeController](https://github.com/etiennecoutaud/curlme-controller/workflows/CurlMeController/badge.svg)](https://github.com/etiennecoutaud/curlme-controller/actions)
[![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/etiennecoutaud/curlme-controller)](https://hub.docker.com/repository/docker/etiennecoutaud/curlme-controller)
[![codecov](https://codecov.io/gh/etiennecoutaud/curlme-controller/branch/master/graph/badge.svg)](https://codecov.io/gh/etiennecoutaud/curlme-controller)
[![Go Report Card](https://goreportcard.com/badge/github.com/etiennecoutaud/curlme-controller)](https://goreportcard.com/report/github.com/etiennecoutaud/curlme-controller)

CurlMe controller is a simple Kubernetes controller that watches `ConfigMap` with `x-k8s.io/curl-me-that` annotation.

## Overview

CurMe controller watches `ConfigMap` resources like below
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: example
  annotations:
    x-k8s.io/curl-me-that: joke=curl-a-joke.herokuapp.com
data: {}
```

`curl-a-joke.herokuapp.com` response request will be stored into the configmap datas under `joke` key.

If the request fails, a warning event will be raised

## Quickstart

To deploy the controller into your kubernetes cluster :

```bash
$ kubectl apply -f https://raw.githubusercontent.com/etiennecoutaud/curlme-controller/master/manifests/all-in-one.yaml
```

If you use `Prometheus operator` and `ServiceMonitor CRD` into your cluster. You can monitor your cluster with :
```bash
$ kubectl apply -f https://raw.githubusercontent.com/etiennecoutaud/curlme-controller/master/manifests/monitoring.yaml
```

## Run locally

To run the app locally make sure you have `go 1.14` install and `KUBECONFIG` environment var set.

```bash
$ make run-local
```

## Questions

[docs/questions](docs/questions.md)