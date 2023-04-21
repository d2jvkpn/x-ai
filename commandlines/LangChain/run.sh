#! /usr/bin/env bash
set -eu -o pipefail
_wd=$(pwd)
_path=$(dirname $0 | xargs -i readlink -f {})

apt update
apt -y install vim

pip3 install ipython
pip3 install -r requirements.txt

mkdir data

python3 langchain_02.py configs/config.yaml data/result.yaml
