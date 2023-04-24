import os, sys, json

import yaml, textract
from langchain.embeddings.openai import OpenAIEmbeddings
from langchain.text_splitter import CharacterTextSplitter
from langchain.vectorstores import FAISS

"""
sources:
- {name: "2023_GPT4All_Technical_Report", type: "pdf", source: "2023_GPT4All_Technical_Report.pdf"}
"""

####
cf, prefix = sys.argv[1:3]
cfd = os.path.dirname(cf)

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

api_key = os.environ['OPENAI_API_KEY']
embeddings = OpenAIEmbeddings(openai_api_key=api_key)

text_docs = []
for d in config["sources"]:
    # d["type"]
    source = os.path.join(cfd, d["source"])
    text = textract.process(source)
    text_docs.append(text.decode().strip())

text_splitter = CharacterTextSplitter(
  separator = "\n",
  chunk_size = 1000,
  chunk_overlap = 200,
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
