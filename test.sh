#!/usr/bin/env bash

randomString=$(tr -dc A-Za-z0-9 </dev/urandom | head -c 20 ; echo '')

echo "$randomString"