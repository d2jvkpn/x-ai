#! /usr/bin/env bash
set -eu -o pipefail

_wd=$(pwd)
_path=$(dirname $0 | xargs -i readlink -f {})
# set -x

####
## v0.1.1
# cargo install llm-cli

## v0.2.0-dev
cargo install --git https://github.com/rustformers/llm llm-cli

exit

echo -e "{{PROMPT}}" > prompt.txt

####
llm infer --model-architecture gpt-j \
  --model-path gpt4all-j-q4_0.bin \
  --prompt "Tell me how cool the Rust programming language is:"

llm chat --help

llm chat --model-architecture gpt-j \
  --model-path gpt4all-j-q4_0.bin \
  --prelude-prompt-file prompt.txt \
  --message-prompt-prefix "Tell me how cool the Rust programming language is:"

####
llm chat --model-architecture gpt-j \
  --model-path gpt4all-j-q4_0.bin \
  --prelude-prompt-file prompt.txt \
  --message-prompt-prefix "English learning:"

#=> how to learn english effective

llm chat --model-architecture gpt-j \
  --model-path gpt4all-j-q4_0.bin \
  --prelude-prompt-file prompt.txt \
  --message-prompt-prefix "AI programming:"

#=> how to learn ai programming
#=> c++

####
llm chat --model-architecture gpt-j \
  --model-path gpt4all-j-q5_1-ggjt.bin \
  --prelude-prompt-file prompt.txt \
  --message-prompt-prefix "Tell me how cool the Rust programming language is:"

####
llm chat --model-architecture gpt-j \
  --use-gpu \
  --model-path gpt4all-j-f16.bin \
  --prelude-prompt-file prompt.txt \
  --message-prompt-prefix "Tell me something that are cool:"

####
llm chat --model-architecture llama \
  --model-path dolly-v2-7b-q5_1.bin \
  --prelude-prompt-file prompt.txt \
  --message-prompt-prefix "Tell me how cool the Rust programming language is:"
