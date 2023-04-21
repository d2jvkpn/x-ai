# pip install langchain openai PyPDF2 faiss-cpu
import os, sys, yaml

from langchain.document_loaders import PyPDFLoader
from langchain.embeddings.openai import OpenAIEmbeddings
from langchain.text_splitter import CharacterTextSplitter
from langchain.vectorstores import FAISS, ElasticVectorSearch, Pinecone, Weaviate
from langchain.chains.question_answering import load_qa_chain
from langchain.llms import OpenAI
from langchain.chains.summarize import load_summarize_chain

# yf, output = "config.yaml", "data/result.yaml"
yf, output = sys.argv[1:3]

f = open(yf , "r");
config = yaml.safe_load(f)
f.close()

# os.environ['OPENAI_API_KEY'] = config["chatgpt"]["api_key"]
# os.environ["OPENAI_API_BASE"] = config["chatgpt"]["url"]
api_key = config["chatgpt"]["api_key"]
llm = OpenAI(temperature=1, openai_api_key=api_key)
embeddings = OpenAIEmbeddings(openai_api_key=api_key)

# ! wget https://s3.amazonaws.com/static.nomic.ai/gpt4all/2023_GPT4All_Technical_Report.pdf
pdf = config["parameters"]["files"][0]
loader = PyPDFLoader(pdf)
docs = loader.load()

def pdf2texts(fp):
    loader = PyPDFLoader(fp)
    pages = loader.load_and_split()

    raw_text = ''
    for i, p in enumerate(pages):
        text = p.page_content
        if text: raw_text += "\n" + text.strip()

    text_splitter = CharacterTextSplitter(
      separator = "\n",
      chunk_size = 1000,
      chunk_overlap = 200,
      length_function = len,
    )

    texts = text_splitter.split_text(text=raw_text)
    return texts

def docs2texts(docs):
    raw_text = ''
    for i, p in enumerate(docs):
        text = p.page_content
        if text: raw_text += "\n" + text.strip()

    text_splitter = CharacterTextSplitter(
      separator = "\n",
      chunk_size = 1000,
      chunk_overlap = 200,
      length_function = len,
    )

    texts = text_splitter.split_text(text=raw_text)
    return texts

chain = load_qa_chain(llm, chain_type="stuff")

# texts = pdf2texts(pdf)
texts = docs2texts(docs)
docsearch = FAISS.from_texts(texts, embeddings)
result = {"queries": []}

for q in config["parameters"]["queries"]:
    q = q.strip()
    ss = docsearch.similarity_search(q, k=2)
    ans = chain.run(input_documents=ss, question=q)
    result["queries"].append({"query": q, "answer": ans.strip()})

if config["parameters"]["summarize"]:
    chain = load_summarize_chain(llm, chain_type="map_reduce", verbose=True)
    ans = chain.run(docs).strip()
    result["summary"] = ans.strip()

os.makedirs(os.path.dirname(output), exist_ok=True)

with open(output, 'w') as f:
    yaml.dump(result, f, default_flow_style=False)
