#!/usr/bin/env bash
#
# Check host temperature
#
ALERT_LIMIT=80

function ok {
    echo "OK"
    exit 0
}

function alertMessage {
    MESSAGES="$1"
    echo -e "High temperature detected:\n$MESSAGES"
    exit 0
}

function getTemperatures {
    sensors -j 2>/dev/null | jq -r '
        to_entries[] |
        .key as $chip |
        .value | to_entries[] |
        select(.value | type == "object") |
        .key as $sensor |
        .value | to_entries[] |
        select(.key | endswith("_input")) |
        [$chip, $sensor, (.value | tostring)] |
        join("|")
    '
}

function checkTemperatures {
    IS_TRIGGERED=0
    ALERT_MESSAGES=""
    while IFS= read -r LINE; do
        CHIP=$(echo "$LINE" | cut -d'|' -f1)
        SENSOR=$(echo "$LINE" | cut -d'|' -f2)
        TEMP=$(echo "$LINE" | cut -d'|' -f3 | awk '{printf "%d", $1}')
        if [ "$TEMP" -ge "$ALERT_LIMIT" ] 2>/dev/null; then
            ALERT_MESSAGES="${ALERT_MESSAGES}${CHIP} / ${SENSOR} = ${TEMP}C (limit: ${ALERT_LIMIT}C)\n"
            IS_TRIGGERED=1
        fi
    done
    if [ "$IS_TRIGGERED" -eq 1 ]; then
        alertMessage "$(echo -e "$ALERT_MESSAGES" | sed '/^$/d')"
    fi
    ok
}

function main {
    getTemperatures | checkTemperatures
}
main "$@"
exit 0
