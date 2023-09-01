#! /bin/bash


wget -qO- https://github.com/reddydinesh427/upgraded-config-mang/releases/download/v1.0.0/provisioner.tar | tar xvz -C /tmp/provisioner

mkdir /etc/upgraded-config-mang
mv /tmp/provisioner/ /etc/upgraded-config-mang/

chmod +x /etc/upgraded-config-mang/upgraded-config-mang-linux

cd /etc/upgraded-config-mang/

./upgraded-config-mang-linux -tasks ./tasks.json