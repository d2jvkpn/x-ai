#!/bin/bash
set -eu -o pipefail; _wd=$(pwd); _path=$(dirname $0)


[ $# -eq 0 ] && {
    >&2 echo "Hello! How can I assist you today?"
    exit 0
}

#### 1.
for m in yq jq curl; do
    command -v $m > /dev/null || { >&2 echo "command $m not found"; exit 1; }
done

#### 2.
app_dir=~/apps/chatgpt
config=$app_dir/configs/chatgpt.yaml
account=${account:-chatgpt}

ls $config > /dev/null
mkdir -p $app_dir/data

# OPENAI_API_Key
api_key=$(yq ".$account.api_key" $config)

# gpt-3.5-turbo gpt-4-turbo gpt-4 gpt-4o
model=$(yq ".$account.model" $config)

# -x socks5h://localhost:1080
proxy=$(yq ".$account.proxy" $config)

if [[ "$api_key$model$proxy" == *"null"* ]]; then
    >&2 echo "invalid config: $config"
    exit 1
fi

#### 3.
question="$*"
if [[ "$question" == @* ]]; then
    question=$(cat ${question:1})
fi

tag=$(date +%FT%T-%s | sed 's/:/-/g')
echo "==> $model@$tag: $question"

ques_file=$app_dir/data/${tag}_quesiton.json
ans_file=$app_dir/data/${tag}_answer.json

jq -n \
  --arg model "$model" \
  --arg content "$question" \
  --argjson max_tokens "1024" \
  --argjson temperature "1.0" \
  '{model: $model, messages: [{"role": "user", "content": $content}],
   max_tokens: $max_tokens, temperature: $temperature}' > $ques_file

####
curl https://api.openai.com/v1/chat/completions $proxy \
  -H 'Content-Type: application/json' -H "Authorization: Bearer $api_key" \
  -d @$ques_file > $ans_file || { rm $ans_file; exit 1; }

jq -r .choices[].message.content $ans_file || cat $ans_file

{
    echo -e "\n#### QA"
    yq -P -oy eval .  $ques_file
    echo -e "---"
    yq -P -oy eval .  $ans_file
} >> $app_dir/data/chatgpt_QA_$(date +%F).yaml

rm $ques_file $ans_file
