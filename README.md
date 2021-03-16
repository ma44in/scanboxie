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

## Cover Arts using beets

```sh
apt install python3-pip
pip3 install https://github.com/beetbox/beets/tarball/master
```

 ~/.local/bin/beet import /mnt/z/Musik

 ~/.local/bin/beet import --nocopy --nowrite --quiet /mnt/z/Musik

 ~/.config/beets/config.yaml

 plugins: fetchart

fetchart:
    cautious: true
    cover_names: front back
    sources: amazon *

~/.local/bin/beet fetchart -f Adele