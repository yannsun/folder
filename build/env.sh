#!/bin/bash
# Check env

cd `dirname $0`/.. || exit 1;
projectpath=`pwd`
export GOPATH="${projectpath}"
export GOBIN="${projectpath}/bin"
export PATH=$PATH:$projectpath:${projectpath}/bin

gobin=`which go`
if [ ! -x "$gobin" ]; then
    echo "No gobin found."
    exit 2
fi

git_rev ()
{
    d=`date +%Y%m%d`
    c=`git rev-list --full-history --all --abbrev-commit | wc -l | sed -e 's/^ *//'`
    h=`git rev-list --full-history --all --abbrev-commit | head -1`
    echo ${c}:${h}:${d}
}
