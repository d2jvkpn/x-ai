#! /usr/bin/env bash
set -eu -o pipefail
_wd=$(pwd)
_path=$(dirname $0 | xargs -i readlink -f {})

curl https://api.openai.com/v1/models \
  -H "Authorization: Bearer $OPENAI_API_KEY"
  # -x socks5://localhost:1081          \
  # -H "OpenAI-Organization: $OPENAI_ORG_ID" > OpenAI_Models.json

exit
jq '.data | sort_by(.created) | reverse' OpenAI_Models.json |
  less
