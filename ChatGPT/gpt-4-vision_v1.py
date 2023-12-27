#! /usr/bin/python3
import sys, base64, requests, imghdr

# $ OPENAI_API_KEY=xxxxxxxx python3 gpt-4-vision.py my_image.jpeg

# OPENAI_API_KEY
api_key = os.getenv("OPENAI_API_KEY")
if api_key is None: sys.exit("OPENAI_API_KEY is unset")

# Path to your image
image_path = sys.argv[1]

# Function to encode the image
def encode_image(image_path):
  with open(image_path, "rb") as image_file:
    return base64.b64encode(image_file.read()).decode('utf-8')

# detect image fomat: jpeg, png
image_type = imghdr(image_path)

# Getting the base64 string
base64_image = encode_image(image_path)

headers = {
  "Content-Type": "application/json",
  "Authorization": f"Bearer {api_key}"
}

# url: https://upload.wikimedia.org/wikipedia/commons/thumb/d/dd/Gfp-wisconsin-madison-the-nature-boardwalk.jpg/2560px-Gfp-wisconsin-madison-the-nature-boardwalk.jpg
payload = {
  "model": "gpt-4-vision-preview",
  "messages": [
    {
      "role": "user",
      "content": [
        { "type": "text", "text": "What’s in this image?" },
        {
          "type": "image_url",
          "image_url": { "url": f"data:image/{image_type};base64,{base64_image}" }
        }
      ]
    }
  ],
  "max_tokens": 300
}

response = requests.post(
  "https://api.openai.com/v1/chat/completions",
  headers=headers, json=payload,
)

print(response.json())
