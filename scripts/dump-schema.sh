#!/usr/bin/env bash

# Variables
HOST=
COOKIE=
OUT=../data/schema.json
verbose=

function usage() {
  echo "Download the schema file from the IPA host."
  echo "  This is used to update the schema.json to the latest version."
  echo ""
  echo "Usage: --host {ipa host url} --cookie {ipa_session cookie} [--out schema file] [-v|--verbose]"
  echo ""
  echo "  --host ipa host url   The url to your IPA host server."
  echo "                          This should include the protocol (https), but no trailing slash."
  echo "                          ex: 'https://dc1.test.local'"
  echo "  --cookie              The ipa session cookie to use."
  echo "                          This comes from your browser."
  echo "                          See developing.md for details on how to get this value."
  echo "  --out schema file     The full or relative path and file name to the schema file."
  echo "                          This defaults to '$OUT'"
  echo "  -v, --verbose         Turn on verbose logging."
  echo "                          The default is off/not set."
  echo ""
}

while [ -n "${1-}" ]; do # while loop starts

  case "$1" in

  -h | --help | /?) usage && exit ;;

  --host) shift && HOST=$1 ;;

  --cookie) shift && COOKIE=$1 ;;

  --out) shift && OUT=$1 ;;

  -v | --verbose) verbose=-v ;;

  esac

  shift

done

if [[ -z $HOST || -z $COOKIE ]]; then
  echo "You must specify both a host name and cookie value."
  exit
fi

if [[ ! -z "$verbose" ]]; then
  echo "Host = '${HOST}'"
  echo "Cookie = '$COOKIE'"
  echo "Schema file = '$OUT'"
fi


curl "$HOST/ipa/session/json" -H "Origin: $HOST" -H 'Content-Type: application/json' \
  -H "Accept: application/json" -H "Cookie: ipa_session=$COOKIE" -H "Referer: $HOST/ipa/ui/" \
  --data-binary '{"method":"schema","params":[[],{"version":"2.231"}]}' --insecure > $OUT
