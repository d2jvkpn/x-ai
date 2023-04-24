import os, sys

from langchain.llms import OpenAI
from langchain.chains.question_answering import load_qa_chain
from langchain.embeddings.openai import OpenAIEmbeddings
from langchain.vectorstores import FAISS

prefix, query = sys.argv[1:3]

api_key = os.environ.get('OPENAI_API_KEY')
if api_key is None:
    sys.exit("OPENAI_API_KEY is unset")

embeddings = OpenAIEmbeddings(openai_api_key=api_key)
llm = OpenAI(temperature=1, openai_api_key=api_key)
qa_chain = load_qa_chain(llm, chain_type="stuff")

fdir = os.path.dirname(prefix)
fname = os.path.basename(prefix)
faiss_index = FAISS.load_local(fdir, embeddings, fname)

ss = faiss_index.similarity_search(query.strip(), k=5)
ans = qa_chain.run(input_documents=ss, question=query)

print(ans.strip())
