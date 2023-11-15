#! /usr/bin/env bash
set -eu -o pipefail
_wd=$(pwd)
_path=$(dirname $0 | xargs -i readlink -f {})

app_dir=~/.local/apps/chatgpt
mkdir -p $app_dir/data
# token=${ChatGPT_Token:-Your_Default_ChatGPT_API_Key}
[ -f $app_dir/env ] && source $app_dir/env

set_proxy=""
# CURL_Proxy=socks5h://localhost:1081
CURL_Proxy=$(printenv CURL_Proxy || true)
[ ! -z "$CURL_Proxy" ] && set_proxy="-x $CURL_Proxy"

curl https://api.openai.com/v1/models \
  -H "Authorization: Bearer $OPENAI_API_KEY" \
  $set_proxy
  # -H "OpenAI-Organization: $OPENAI_ORG_ID" > OpenAI_Models.json

exit
jq '.data | sort_by(.created) | reverse' OpenAI_Models.json |
  less
