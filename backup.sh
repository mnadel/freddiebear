#!/usr/bin/env bash

DIR=$(dirname $(readlink -f $0))
LOG=${DIR}/backup.log
DATADIR="${HOME}/Library/Group Containers/9K33E3U3T4.net.shinyfrog.bear/Application Data"
MAXLOG=$((1024 * 128))

function log {
    echo $(date "+%Y-%m-%d %T") :: $*
}

cd "${DIR}"

# if invoked without args, redirect stdout to logfile
# else assume "interactive" mode
if [ $# -eq 0 ]; then
    exec >> "${LOG}" 2>&1
fi

# truncate logfile if it's too big
if [ -f "${LOG}" ]; then
    sz=$(stat -f%z "${LOG}")
    if [ $sz -gt ${MAXLOG} ]; then
        echo -n > "${LOG}"
        log "rolled ${LOG}"
    fi
fi

log "exporting notes"
freddiebear export .

rsync --dry-run -avz "${DATADIR}/Local Files" . | grep -qs "Local Files" >/dev/null 2>&1
if [ $? -eq 0 ]; then
    log "found assets to sync"
    rsync -avz "${DATADIR}/Local Files" .
fi

changes=$(git status --porcelain 2>/dev/null | wc -l)
if [ ${changes} -eq 0 ]; then
    exit 0
fi

log "adding changes"
git add -A .

log "committing changes"
git commit -am autosync

log "pushing changes"
git push
