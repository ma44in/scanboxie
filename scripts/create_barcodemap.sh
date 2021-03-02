#!/bin/sh

# TODO Move this into go app

echo "" > ./mybarcode.csv

find /mnt/z/Musik -name "*.mp3" -exec dirname {} \; | sort | uniq | sed 's#/mnt/z/##' | while read LINE; do
	HASH="$(echo "${LINE}" | sha1sum  | cut -c1-8)"
	echo "${HASH},mpc_add_and_play,${LINE}" >> ./mybarcode.csv
done
