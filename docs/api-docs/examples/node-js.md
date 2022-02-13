### Examples with Nodejs

We use sessions for authentication.

We recommend getting the session TOKEN using Postman (Each is valid for 7days) 

Example Cookie `xa=MTY0MTU0NTk4NXxOd3dBTkRVMFRVZFFORVl6VWtsTFMxQTJRazlMTkZSSU4wUlhOREkxV0ZGWVR6ZE5NMDR5UVU4eVJWWkZTbFJQVERkVVdFTkdVRkU9fKt_QL_t2nDecucCWu4sdhjklN9ad6puZI7gh3UL4_ss`


---
#### Get account Info

```js
var request = require('request');
var options = {
  'method': 'GET',
  'url': 'https://gateway.xoniaapp.com/api/account',
  'headers': {
    'Cookie': 'xa=MTY0MTU0NjEwM3xOd3dBTkRVMFRVZFFORVl6VWtsTFMxQTJRazlMTkZSSU4wUlhOREkxV0ZGWVR6ZE5NMDR5UVU4eVJWWkZTbFJQVERkVVdFTkdVRkU9fGHOs7Be3JdwcESK7d7Sa2RkwlRkz3IZrk7ya2MNcUvE'
  }
};
request(options, function (error, response) {
  if (error) throw new Error(error);
  console.log(response.body);
});

```
---

#### Update the account
---

```js
var request = require('request');
var fs = require('fs');
var options = {
  'method': 'PUT',
  'url': 'https://gateway.xoniaapp.com/api/account',
  'headers': {
    'Cookie': 'xa=MTY0MTU0NjU4M3xOd3dBTkZCWE5WUlZWRFJSVkZCUFJEVlBXbFpNTjFwR1QwVktSRk15VTFwV05VWlRXVXRDVVVoWU1rdFROVXhQVFU1UFdUTkpWa0U9fNFDuwfftYXD9rLM85W5PPm9nYwlN-FNQkkR8vxVgMjY'
  },
  formData: {
    'email': 'string',
    'image': {
      'value': fs.createReadStream('/path/to/file'),
      'options': {
        'filename': 'filename'
        'contentType': null
      }
    },
    'username': 'string'
  }
};
request(options, function (error, response) {
  if (error) throw new Error(error);
  console.log(response.body);
});

```
---
