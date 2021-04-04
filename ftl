#!/usr/bin/env sh
set -x -e -u

ssh $1 rm -rf /root/go/src/github.com/ftlops/ftl/
ssh $1 mkdir -p /root/go/src/github.com/ftlops/ftl/{log,ops}
scp *.go $1:/root/go/src/github.com/ftlops/ftl/
scp log/*.go $1:/root/go/src/github.com/ftlops/ftl/log
scp ops/*.go $1:/root/go/src/github.com/ftlops/ftl/ops

date=$(date +%s)
dst=/var/log/ftl/$date/plans
ssh $1 mkdir -p $dst
scp $2 $1:$dst/../$2
ssh $1 go run $dst/../$2
