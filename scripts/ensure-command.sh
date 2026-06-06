#!/usr/bin/env bash
set -e

command -v $1 >/dev/null 2>&1 || $2
