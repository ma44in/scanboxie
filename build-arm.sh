#!/bin/sh

set -x -e -u

GOOS=linux GOARCH=arm GOARM=5 go build -o scanboxie-arm

set +x
echo "Deploy with:"
echo "  ansible-playbook -i ./setup-pi-ansible/inventory ./setup-pi-ansible/playbook.yml --tags copy-scanboxie-binary"
