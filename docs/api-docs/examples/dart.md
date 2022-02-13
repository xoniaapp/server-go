### Examples with Dart

We use sessions for authentication.

We recommend getting the session TOKEN using Postman (Each is valid for 7days)

Example Cookie `xa=MTY0MTU0NTk4NXxOd3dBTkRVMFRVZFFORVl6VWtsTFMxQTJRazlMTkZSSU4wUlhOREkxV0ZGWVR6ZE5NMDR5UVU4eVJWWkZTbFJQVERkVVdFTkdVRkU9fKt_QL_t2nDecucCWu4sdhjklN9ad6puZI7gh3UL4_ss`

---

#### Get account Info


```dart
var headers = {
  'Cookie': 'xa=MTY0MTU0NjEwM3xOd3dBTkRVMFRVZFFORVl6VWtsTFMxQTJRazlMTkZSSU4wUlhOREkxV0ZGWVR6ZE5NMDR5UVU4eVJWWkZTbFJQVERkVVdFTkdVRkU9fGHOs7Be3JdwcESK7d7Sa2RkwlRkz3IZrk7ya2MNcUvE'
};
var request = http.Request('GET', Uri.parse('https://gateway.xoniaapp.com/api/account'));

request.headers.addAll(headers);

http.StreamedResponse response = await request.send();

if (response.statusCode == 200) {
  print(await response.stream.bytesToString());
}
else {
  print(response.reasonPhrase);
}
```

---

#### Update the account
---

```dart
var headers = {
  'Cookie': 'xa=MTY0MTU0NjU4M3xOd3dBTkZCWE5WUlZWRFJSVkZCUFJEVlBXbFpNTjFwR1QwVktSRk15VTFwV05VWlRXVXRDVVVoWU1rdFROVXhQVFU1UFdUTkpWa0U9fNFDuwfftYXD9rLM85W5PPm9nYwlN-FNQkkR8vxVgMjY'
};
var request = http.MultipartRequest('PUT', Uri.parse('https://gateway.xoniaapp.com/api/account'));
request.fields.addAll({
  'email': 'string',
  'username': 'string'
});
request.files.add(await http.MultipartFile.fromPath('image', '/path/to/file'));
request.headers.addAll(headers);

http.StreamedResponse response = await request.send();

if (response.statusCode == 200) {
  print(await response.stream.bytesToString());
}
else {
  print(response.reasonPhrase);
}

```
---
