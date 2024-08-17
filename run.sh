#!/bin/bash

set -e

dockerd -p /var/run/docker.pid &
./rce
