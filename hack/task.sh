#!/bin/bash

host="${HOST:-localhost:8080}"

curl -X POST ${host}/tasks -d \
'{
    "name":"Windup",
    "locator": "windup",
    "addon": "windup",
    "data": {
      "application": 3,
      "path": "/etc"
    }
}' | jq -M .
