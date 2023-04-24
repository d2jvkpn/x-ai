import os, sys, json, yaml, textract

from langchain.chains.summarize import load_summarize_chain
from langchain.llms import OpenAI
from langchain.text_splitter import CharacterTextSplitter
from langchain.docstore.document import Document

cf, findex = sys.argv[1:3]
cfd = os.path.dirname(cf)
findex = int(findex)

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

if findex >= len(config["sources"]):
    sys.exit("findx out of index range")

llm = OpenAI(temperature=0, openai_api_key=api_key)
summarize_chain = load_summarize_chain(llm, chain_type="map_reduce", verbose=False)
text_splitter = CharacterTextSplitter()

source = os.path.join(cfd, config["sources"][findex]["source"])
content = textract.process(source)

docs = []
for text in text_splitter.split_text(content.decode()):
    docs.append(Document(page_content=text.strip()))

ans = summarize_chain.run(docs)
println(ans.strip())
