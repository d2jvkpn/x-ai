import os, sys, json, yaml, textract

from langchain.embeddings.openai import OpenAIEmbeddings
# from langchain.text_splitter import CharacterTextSplitter
from langchain.text_splitter import TokenTextSplitter
from langchain.vectorstores import FAISS

"""
config file UUID.yaml
```yaml
id: UUID
title: 2023 GPT4All Technical Report
created: 1682391010
sources:
- { title: "2023_GPT4All_Technical_Report", type: "pdf", source: "UUID_doc001.pdf", size: 3651423 }
```
"""

####
cf, prefix = sys.argv[1:3]
cfd = os.path.dirname(cf)

api_key = os.environ.get('OPENAI_API_KEY')
if api_key is None:
    sys.exit("OPENAI_API_KEY is unset")

f = open(cf, "r")

if cf.endswith(".json"):
    config = json.loads(f.read())
    # prefix = cf[:-5]
elif cf.endswith(".yaml"):
    config = yaml.safe_load(f)
    # prefix = cf[:-5]
elif cf.endswith(".yml"):
    config = yaml.safe_load(f)
    # prefix = cf[:-4]
else:
    sys.exit("unknown input file type")

f.close()

embeddings = OpenAIEmbeddings(openai_api_key=api_key)

text_docs = []
for d in config["sources"]:
    # d["type"]
    source = os.path.join(cfd, d["source"])
    text = textract.process(source)
    text_docs.append(text.decode().strip())

#text_splitter = CharacterTextSplitter(
#  separator = "\n",
#  chunk_size = 1000,
#  chunk_overlap = 100,
#  length_function = len,
#)

text_splitter = TokenTextSplitter(
  chunk_size = 500,
  chunk_overlap = 50,
  length_function = len,
)

chunks = []
for d in text_docs:
    chunks.extend(text_splitter.split_text(d))
# sum([[1, 2], [3, 4]], [])

faiss_index = FAISS.from_texts(chunks, embeddings)
fdir = os.path.dirname(prefix)
fname = os.path.basename(prefix)
faiss_index.save_local(fdir, fname)
faiss_index = FAISS.load_local(fdir, embeddings, fname)

print("ok")
