#!/usr/bin/env sh
#
# Generate Table of Contents for README.md
#
set -e

exec 5>&1
output="$(tocenize README.md | tee /dev/fd/5)"
[ -z "$output" ]
