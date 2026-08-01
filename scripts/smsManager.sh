#!/usr/bin/env bash
function ok {
    echo "OK"
    exit 0
}

function alertMessage {
    MESSAGE="$1"
    echo -e "SMS manager has runtime errors\n$MESSAGE"
    exit 0
}

function checkLogs {
    /usr/bin/journalctl -u smsmanager.service --since=-15m | grep ERROR
}

function main {
    MESSAGE=$(checkLogs)
    if [ "$MESSAGE" != "" ]; then
        alertMessage "$MESSAGE"
    fi
    ok
}
main "$@"
