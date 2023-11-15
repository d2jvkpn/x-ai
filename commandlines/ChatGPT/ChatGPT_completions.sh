#! /usr/bin/env bash
set -eu -o pipefail
_wd=$(pwd)
_path=$(dirname $0 | xargs -i readlink -f {})

for m in yq jq curl; do
    command -v $m > /dev/null || { >&2 echo "command $m not found"; exit 1; }
done

app_dir=~/.local/apps/chatgpt
mkdir -p $app_dir/data
# token=${ChatGPT_Token:-Your_Default_ChatGPT_API_Key}
[ -f $app_dir/env ] && source $app_dir/env

[ $# -eq 0 ] && { >&2 echo "Pass your question as argument(s)!"; exit 1; }
[ -z "${ChatGPT_Token}" ] && { >&2 echo "ChatGPT_Token is unset"; exit 1; }

question="$*"
if [[ "$question" == @* ]]; then
    question=$(cat ${question:1})
fi

tag=$(date +%FT%T-%s | sed 's/:/-/g')
echo ">>> $tag: $question"

ques_file=$app_dir/data/${tag}_quesiton.json
ans_file=$app_dir/data/${tag}_answer.json

# --arg model "${ChatGPT_Model:-gpt-3.5-turbo}" \
jq -n \
  --arg model "${ChatGPT_Model:-gpt-4}" \
  --arg content "$question" \
  --argjson max_tokens "${ChatGPT_MaxTokens:-1024}" \
  --argjson temperature "${ChatGPT_Temperature:-1.0}" \
  '{model: $model, messages: [{"role": "user", "content": $content}],
    max_tokens: $max_tokens, temperature: $temperature}' > $ques_file

set_proxy=""
# CURL_Proxy=socks5h://localhost:1081
[ ! -z "${CURL_Proxy:-}" ] && set_proxy="-x $CURL_Proxy"

curl https://api.openai.com/v1/chat/completions \
  $set_proxy                                  \
  -H 'Content-Type: application/json'         \
  -H "Authorization: Bearer $ChatGPT_Token"   \
  -d @$ques_file > $ans_file || { rm $ans_file; exit 1; }

jq -r .choices[].message.content $ans_file || cat $ans_file

{
    echo -e "\n#### QA"
    yq -P -oy eval .  $ques_file
    echo -e "---"
    yq -P -oy eval .  $ans_file
} >> $app_dir/data/chatgpt_QA_$(date +%F).yaml

rm $ques_file $ans_file
