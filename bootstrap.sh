#! /bin/bash
export PATH=$PATH:/usr/local/go/bin
which go
if [ $? -eq 0 ]; then
        echo "Go is installed"
        export PATH=$PATH:/usr/local/go/bin
else
curl -OL https://golang.org/dl/go1.21.0.linux-amd64.tar.gz; tar -C /usr/local -xvf go1.21.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
fi
apt update; apt install git
git clone https://github.com/reddydinesh427/upgraded-config-mang.git
cd upgraded-config-mang

go run main.go -tasks tasks.json
