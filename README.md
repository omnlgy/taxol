# "Microservices with Go" course project

This is the starter code for the "Microservices with Go" project.

## Project overview

In this project‑driven course, you’ll build the backend microservices system for a Uber‑style ride‑sharing app from the ground up—using Go, Docker, and Kubernetes.

By the end, you’ll have a fully deployed, horizontally scalable ride‑sharing system that’s ready for real traffic. Plus, you’ll walk away with reusable template for building future distributed projects—accelerating your path to become a lead engineer.


## Installation
The project requires a couple tools to run, most of which are part of many developer's toolchains.

- Docker
- Go
- Tilt
- A local Kubernetes cluster

### MacOS

1. Install Docker for Desktop from [Docker's official website](https://www.docker.com/products/docker-desktop/)

2. Install Minikube from [Minikube's official website](https://minikube.sigs.k8s.io/docs/)

3. Install Tilt from [Tilt's official website](https://tilt.dev/)

4. Install Go

5. Make sure [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl-macos/) is installed.

## Run

```bash
tilt up
```

## Monitor

```bash
kubectl get pods
```

or

```bash
minikube dashboard
```

