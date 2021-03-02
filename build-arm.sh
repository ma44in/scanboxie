#!/bin/sh

set -x -e -u

GOOS=linux GOARCH=arm GOARM=5 go build -o scanboxie-arm

cd ./setup-pi-ansible
ansible-playbook -i ./inventory ./playbook.yml --tags copy-scanboxie-binary
