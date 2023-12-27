#! /bin/env python3

import base64, requests, sys
from os import path
from timeit import default_timer as timer
from datetime import timedelta
# import time
# time.sleep(5) # sleep 5s

import yaml

# Path to your image
yaml_file, output = sys.argv[1], sys.argv[2]

# Function to encode the image
def image_to_dict(image_path):
  with open(image_path, "rb") as image_file:
    base64_image = base64.b64encode(image_file.read()).decode('utf-8')

  ans = {
    "type": "image_url",
    "image_url": {
      "url": f"data:image/jpeg;base64,{base64_image}"
    },
  }

  return ans

def gpt_vision(settings, item, images):
  # Getting the base64 string
  api_key = settings["api_key"]
  prompt = settings["prompt"]
  model = settings.get("model", "gpt-4-vision-preview")
  url = settings.get("url", "https://api.openai.com")

  headers = {"Content-Type": "application/json", "Authorization": f"Bearer {api_key}"}

  payload = {
    "model": model,
    "max_tokens": 300,
    "messages": [
      {
        "role": "user",
        "content": [
          { "type": "text", "text": prompt }
        ]
      }
    ],
  }

  # { "type": "image_url", "image_url": { "url": f"data:image/jpeg;base64,{base64_image}" }}
  payload["messages"][0]["content"].extend([img[0] for img in images])

  response = requests.post(f"{url}/v1/chat/completions", headers=headers, json=payload)

  return (response.status_code, response.json())

####
with open(yaml_file, 'r', encoding='utf-8') as file:
  config = yaml.safe_load(file, )

settings = config["openai"] # keys: api_key, url, prompt, model
items = config["items"]

results = []
for item in items:
  print("==> process:", item["name"])
  images = [(image_to_dict(path.join(item["name"], p)), p) for p in item["images"]]

  start = timer()
  result = gpt_vision(settings, item, images)
  ans = result[1]
  end = timer()
  elapsed = str(timedelta(seconds=end-start))

  if "choices" in ans:
    print("~~~ elapsed:", elapsed, "ans:", ans["choices"][0]["message"]["content"])
  else:
    print("!!! elapsed:", elapsed, "ans:", ans["error"]["message"])

  results.append({
    "name": item["name"],
    "images": [img[1] for img in images],
    "elpased": elapsed,
    "status_code": result[0],
    "answer": ans,
  })

with open(output, 'w', encoding='utf-8') as file:
  yaml.dump(results, file, allow_unicode=True, sort_keys=False)

print(f"==> saved output:", output)
