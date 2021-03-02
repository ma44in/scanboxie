#!/bin/sh

set -u -x -e

ansible-playbook -i ./inventory ./playbook.yml $*