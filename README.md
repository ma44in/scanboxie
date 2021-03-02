# Scanboxie

Execute commands based on barcode scans with the primary goal to play music on a raspberry pi.

TODO

## Build

    GOOS=linux GOARCH=arm GOARM=5 go build -o scanboxie-arm

## Deploy on Raspbery Pi with ansible

Prepare SSH Keys if not present.

```sh
ssh-keygen -b 4096
ssh-copy-id pi@raspberrypi
```

Start ssh-agent to prevent permanent password prompt.

```sh
eval `ssh-agent`
ssh-add
```

Execute Ansible Playbook.

```sh
cd ./setup-pi
ansible-playbook -i ./inventory ./playbook.yml
```
