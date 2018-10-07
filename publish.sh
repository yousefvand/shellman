#!/usr/bin/env bash

vsce login Remisa
vsce publish -p $1
