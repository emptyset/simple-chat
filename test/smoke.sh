#!/usr/bin/env bash

# simple user create
curl -X POST localhost:8080/user -H "Username: emptyset" -H "Password: password"
curl -X POST localhost:8080/user -H "Username: dave" -H "Password: arbitrary"
curl -X POST "localhost:8080/message?s=1&r=2&t=text" -d "howdy"
curl -X POST "localhost:8080/message?s=1&r=1&t=text" -d "my life is made of patterns"
curl -X POST "localhost:8080/message?s=2&r=1&t=text" -d "it's gonna be OK"

# should return two messages
curl -X GET "localhost:8080/message?s=1&r=2&c=5&o=0"
