#!/bin/bash -e

log() {
  echo -e "${NAMI_DEBUG:+${CYAN}${MODULE} ${MAGENTA}$(date "+%T.%2N ")}${RESET}${@}" >&2
}

log "Start ${COMMAND_NAME}"
gitlab-registry-cleaner ${COMMAND_NAME}
