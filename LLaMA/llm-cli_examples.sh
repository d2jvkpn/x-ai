#! /usr/bin/env bash
set -eu -o pipefail
# set -x
_wd=$(pwd)
_path=$(dirname $0 | xargs -i readlink -f {})


#### v0.2.0-dev
cargo install --git https://github.com/rustformers/llm llm-cli
llm --version

exit

[ -s .env ] && source .env
# export ai_models=~/Downloads/ai_models
ai_models=${ai_models:-~/Downloads/ai_models}
prompt=$ai_models/prompt.txt
[ -s $prompt ] || echo -e "{{PROMPT}}" > $prompt

####
llm infer --model-architecture gpt-j \
  --model-path $ai_models/gpt4all-j-q4_0.bin \
  --prompt "Tell me how cool the Rust programming language is:"

llm chat --help

llm chat --model-architecture gpt-j \
  --model-path $ai_models/gpt4all-j-q4_0.bin \
  --prelude-prompt-file $prompt \
  --message-prompt-prefix "Tell me how cool the Rust programming language is:"

####
llm chat --model-architecture gpt-j \
  --model-path $ai_models/gpt4all-j-q4_0.bin \
  --prelude-prompt-file $prompt \
  --message-prompt-prefix "English learning:"

#=> how to learn english effective

llm chat --model-architecture gpt-j \
  --model-path $ai_models/gpt4all-j-q4_0.bin \
  --prelude-prompt-file $prompt \
  --message-prompt-prefix "AI programming:"

#=> how to learn ai programming
#=> c++

####
llm chat --model-architecture gpt-j \
  --model-path $ai_models/gpt4all-j-q5_1-ggjt.bin \
  --prelude-prompt-file $prompt \
  --message-prompt-prefix "Tell me how cool the Rust programming language is:"

####
llm chat --model-architecture gpt-j \
  --use-gpu \
  --model-path $ai_models/gpt4all-j-f16.bin \
  --prelude-prompt-file $prompt \
  --message-prompt-prefix "Tell me something that are cool:"

####
llm chat --model-architecture llama \
  --model-path $ai_models/dolly-v2-7b-q5_1.bin \
  --prelude-prompt-file $prompt \
  --message-prompt-prefix "Tell me how cool the Rust programming language is:"
