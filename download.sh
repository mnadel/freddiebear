#!/usr/bin/env bash

function log {
    local msg=$1
    >&2 echo "$(date) ** ${msg}"
}

if [ -f "${gopath}/bin/freddiebear" ]; then
    log "copying ${gopath}/bin/freddiebear"
    cp "${gopath}/bin/freddiebear" .
    echo ok
    exit 0
elif [ -f "${gopath}/src/github.com/mnadel/freddiebear/freddiebear" ]; then
    log "copying ${gopath}/src/github.com/mnadel/freddiebear/freddiebear"
    cp "${gopath}/src/github.com/mnadel/freddiebear/freddiebear" .
    echo ok
    exit 0
fi

arch=$(uname -m)

log "downloading ver=${alfred_workflow_version} arch=${arch}"
curl -LO https://github.com/mnadel/freddiebear/releases/download/${alfred_workflow_version}/freddiebear.${arch}.gz >/dev/stderr
[ $? -ne 0 ] && echo error && exit 1

log "uncompressing"
gzip -d freddiebear.${arch}.gz >/dev/stderr
[ $? -ne 0 ] && echo error && exit 1

log "renaming"
mv freddiebear.${arch} freddiebear >/dev/stderr
[ $? -ne 0 ] && echo error && exit 1

log "chmodding"
chmod +x freddiebear >/dev/stderr
[ $? -ne 0 ] && echo error && exit 1

echo ok
