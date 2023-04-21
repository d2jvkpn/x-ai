import os, sys, json

import yaml, textract
from langchain.embeddings.openai import OpenAIEmbeddings
from langchain.text_splitter import CharacterTextSplitter
from langchain.vectorstores import FAISS

"""
langchain:
  sources:
  - {type: "pdf", source: "2023_GPT4All_Technical_Report.pdf"}
"""

####
cf = sys.argv[1]
prefix = ""

f = open(cf, "r")

if cf.endswith(".json"):
    config = json.loads(f.read())
    prefix = cf[:-5]
elif cf.endswith(".yaml"):
    config = yaml.safe_load(f)
    prefix = cf[:-5]
elif cf.endswith(".yml"):
    config = yaml.safe_load(f)
    prefix = cf[:-4]
else:
    sys.exit("unknown input file type")

f.close()

api_key = os.environ['OPENAI_API_KEY']
embeddings = OpenAIEmbeddings(openai_api_key=api_key)

text_docs = []
for d in config["langchain"]["sources"]:
    # d["type"]
    text = textract.process(d["source"])
    text_docs.append(text.decode().strip())

text_splitter = CharacterTextSplitter(
  separator = "\n",
  chunk_size = 1000,
  chunk_overlap = 200,
  length_function = len,
)

texts = text_splitter.split_text(text="\n".join(text_docs))

faiss_index = FAISS.from_texts(texts, embeddings)
fdir = os.path.dirname(prefix)
fname = os.path.basename(prefix)
faiss_index.save_local(fdir, fname)
faiss_index = FAISS.load_local(fdir, embeddings, fname)

print("ok")
