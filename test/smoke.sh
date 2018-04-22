#!/usr/bin/env bash

# simple user create
curl -X POST localhost:8080/user -H "Username: emptyset" -H "Password: password"
