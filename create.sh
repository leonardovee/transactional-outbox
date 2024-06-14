#!/bin/bash

curl --header "Content-Type: application/json" --request POST --data '{"Command": "create", "Total": 100}' http://localhost:8080/orders/command
