import os, sys, textract
from random import sample

from langchain.chains.summarize import load_summarize_chain
from langchain.llms import OpenAI
# from langchain.text_splitter import CharacterTextSplitter
from langchain.text_splitter import TokenTextSplitter
from langchain.docstore.document import Document

fp = sys.argv[1]

api_key = os.environ.get('OPENAI_API_KEY')
if api_key is None:
    sys.exit("OPENAI_API_KEY is unset")

llm = OpenAI(temperature=0, openai_api_key=api_key)
summarize_chain = load_summarize_chain(llm, chain_type="map_reduce", verbose=False)

text_splitter = TokenTextSplitter(
  chunk_size = 500,
  chunk_overlap = 0,
  length_function = len,
)

content = textract.process(fp)

docs = []
for text in text_splitter.split_text(content.decode()):
    docs.append(Document(page_content=text.strip()))

if len(docs) <= 5:
    sample_docs = docs
else:
    xx = sample(docs[1:-1], 3)
    sample_docs = [docs[0], xx[0], xx[1], xx[2], docs[-1]]

# This model's maximum context length is 4097 tokens, however you requested 8700 tokens (8444 in your prompt; 256 for the completion). Please reduce your prompt; or completion length.
ans = summarize_chain.run(sample_docs)
print(ans.strip())
