#!/bin/bash

echo "Get Submarine"
curl -X GET -H 'Content-Type: application/json' http://localhost:8080/api/submarine

echo "Get artifact"
curl -X GET -H 'Content-Type: application/json' http://localhost:8080/api/artifact

echo "Get fish"
curl -X GET -H 'Content-Type: application/json' http://localhost:8080/api/fish

echo "Add fish"
curl -X POST -H 'Content-Type: application/json' -d '{"x": 99, "y": 10}' http://localhost:8080/api/fish/add

echo "Add fish"
curl -X POST -H 'Content-Type: application/json' -d '{"x": 45, "y": 10}' http://localhost:8080/api/fish/add

echo "Get fish"
curl -X GET -H 'Content-Type: application/json' http://localhost:8080/api/fish

echo "Update Artifact"
curl -X POST -H 'Content-Type: application/json' -d '{"x": 9, "y": 10}' http://localhost:8080/api/artifact/update

echo "Get Artifact"
curl -X GET -H 'Content-Type: application/json' http://localhost:8080/api/artifact

echo "Move submarine"
curl -X POST -H 'Content-Type: application/json' -d '{"x": -27, "y": 1}' http://localhost:8080/api/submarine/move

echo "Move submarine"
curl -X GET -H 'Content-Type: application/json' http://localhost:8080/api/submarine

curl -X GET -H 'Content-Type: application/json' http://tema2.local:8080/api/submarine

curl -X GET -H 'Content-Type: application/json' http://192.168.7.2:8080/api/submarine
echo "Move submarine"
curl -X POST -H 'Content-Type: application/json' -d '{"x": 10, "y": 10}' http://192.168.171.131:8080/api/submarine/move
