#!/usr/bin/env bash

mkdir -p $GOPATH/src/github.com/nitrous-io
mv ~/code/go-starter-kit $GOPATH/src/github.com/nitrous-io/
ln -s $GOPATH/src/github.com/nitrous-io/go-starter-kit ~/code/go-starter-kit

cd $GOPATH/src/github.com/nitrous-io/go-starter-kit

echo 'Installing dependencies using npm - this might take awhile...'
npm install --no-progress

echo 'Installing slrt...'
go get github.com/olebedev/srlt
echo 'Unpacking Golang dependencies...'
srlt restore

echo 'Installing github.com/jteeuwen/go-bindata...'
go get github.com/jteeuwen/go-bindata/...

echo 'Installing fswatch...'
cd ~
wget https://github.com/emcrisostomo/fswatch/releases/download/1.8.0/fswatch-1.8.0.tar.gz
tar xf fswatch-1.8.0.tar.gz
cd fstwatch-1.8.0
./configure
make
make install
sudo ldconfig
rm -rf fswatch-1.8.0.tar.gz fswatch-1.8.0
