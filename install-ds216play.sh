#!/usr/bin/env bash
env GOOS=linux GOARCH=386 go build -o $GOBIN/directory-monitor-linux-386
#scp hello-world admin@192.168.1.10:/volume1/share/scripts/
