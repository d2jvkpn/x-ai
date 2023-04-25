import os, sys, textract

from langchain.chains.summarize import load_summarize_chain
from langchain.llms import OpenAI
from langchain.text_splitter import CharacterTextSplitter
from langchain.docstore.document import Document

fp = sys.argv[1]

api_key = os.environ.get('OPENAI_API_KEY')
if api_key is None:
    sys.exit("OPENAI_API_KEY is unset")

llm = OpenAI(temperature=0, openai_api_key=api_key)
summarize_chain = load_summarize_chain(llm, chain_type="map_reduce", verbose=False)
text_splitter = CharacterTextSplitter()

content = textract.process(fp)

docs = []
for text in text_splitter.split_text(content.decode()):
    docs.append(Document(page_content=text.strip()))

# This model's maximum context length is 4097 tokens, however you requested 8700 tokens (8444 in your prompt; 256 for the completion). Please reduce your prompt; or completion length.
ans = summarize_chain.run(docs)
print(ans.strip())
