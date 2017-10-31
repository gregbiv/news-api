#!/usr/bin/env sh
#
# Capture and print stdout/stderr, since golint doesn't use proper exit codes
#
set -e

exec 5>&1
rtn=0
for file in "$@"; do
    output="$(golint "$file" 2>&1 | tee /dev/fd/5)"
    if [ ! -z "$output" ]; then
        rtn=1
    fi
done
exit ${rtn}
