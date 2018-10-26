#!/usr/bin/env bash
env GOOS=linux GOARCH=386 go build -o $GOBIN/directory-monitor-linux-386
cp $GOBIN/directory-monitor-linux-386 ~/Dropbox/tmp/
cp -r /home/mha/.go/* ~/Dropbox/tmp/
