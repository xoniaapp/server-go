### Examples with Python


We use sessions for authentication.

We recommend getting the session TOKEN using Postman (Each is valid for 7days)


Example Cookie `xa=MTY0MTU0NTk4NXxOd3dBTkRVMFRVZFFORVl6VWtsTFMxQTJRazlMTkZSSU4wUlhOREkxV0ZGWVR6ZE5NMDR5UVU4eVJWWkZTbFJQVERkVVdFTkdVRkU9fKt_QL_t2nDecucCWu4sdhjklN9ad6puZI7gh3UL4_ss`


---
#### Get account Info

```py
import requests

url = "https://gateway.xoniaapp.com/api/account"

payload={}
headers = {
  'Cookie': 'xa=MTY0MTU0NjEwM3xOd3dBTkRVMFRVZFFORVl6VWtsTFMxQTJRazlMTkZSSU4wUlhOREkxV0ZGWVR6ZE5NMDR5UVU4eVJWWkZTbFJQVERkVVdFTkdVRkU9fGHOs7Be3JdwcESK7d7Sa2RkwlRkz3IZrk7ya2MNcUvE'
}

response = requests.request("GET", url, headers=headers, data=payload)

print(response.text)

```
---

#### Update the account
---

```py
import requests

url = "https://gateway.xoniaapp.com/api/account"

payload={}
files=[
  ('image',('file',open('/path/to/file','rb'),'application/octet-stream'))
]
headers = {
  'Cookie': 'xa=MTY0MTU0NjU4M3xOd3dBTkZCWE5WUlZWRFJSVkZCUFJEVlBXbFpNTjFwR1QwVktSRk15VTFwV05VWlRXVXRDVVVoWU1rdFROVXhQVFU1UFdUTkpWa0U9fNFDuwfftYXD9rLM85W5PPm9nYwlN-FNQkkR8vxVgMjY'
}

response = requests.request("PUT", url, headers=headers, data=payload, files=files)

print(response.text)

```
---
