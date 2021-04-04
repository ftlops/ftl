#!/usr/bin/env sh
set -x -e -u

ssh $1 mkdir -p /root/go/src/github.com/ngrash/ftl
scp *.go $1:/root/go/src/github.com/ngrash/ftl/


date=$(date +%s)
dst=/var/log/ftl/$date/plans
ssh $1 mkdir -p $dst
scp $2 $1:$dst/../$2
ssh $1 go run $dst/../$2
