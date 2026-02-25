#!/bin/sh
case "$1" in
    calc)    exec ./calc "${@:2}" ;;
    plot)    exec ./plot "${@:2}" ;;
    *)       exec ./calc "$@" ;;
esac
