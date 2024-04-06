#!/bin/bash

kubectl apply -f secret.yaml
kubectl apply -f redis.yaml
kubectl apply -f app.yaml
