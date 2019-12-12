#!/bin/bash

while true; do sleep 1h; curl --header "Content-Type: application/json" --request POST --data '{"username":"admin","password":"admin"}' http://127.0.0.1:8080/cookiescleaner; done
