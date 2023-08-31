#! /bin/bash


if [ "$(uname)" == "Darwin" ]; then
    # Mac OS X platform
    which go
    if [ $? -eq 0 ]
    then
        echo "Go is already installed."
        exit 0
    else
        brew install go
    fi    
         
elif [ "$(uname -s)" == "Linux" ]; then
    # Do something under GNU/Linux platform
    which js
    if [ $? -eq 0 ]
    then
        echo "Go is already installed."
        exit 0
    else
        wget https://dl.google.com/go/go1.21.0.linux-amd64.tar.gz
        tar -C /usr/local/ -xzf go1.21.0.linux-amd64.tar.gz
        export PATH=$PATH:/usr/local/go/bin
    fi
else
    echo "This can be built only on linux or darwin systems"
fi
