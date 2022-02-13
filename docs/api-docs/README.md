### Xonia REST API Endpoints
Server URL `https://gateway.xoniaapp.com`

We use sessions for authentication.

### Account Related

---

| Type | Route                        |    Description                   
| ---- | ---------------------------- |---------------------------------|
| GET  | `/account`                   |  Get current user               |
| PUT  | `/account`                   |  Update current user            |
| PUT  | `/account/change-password`   |  Change current users password  |
| POST | `/account/forgot-password`   |  Forgot password request        |
| POST | `/account/login`             |  User Login                     |
| POST | `/account/logout`            |  User Logout                    |
| POST | `/account/register`          |  Create an account              |
| POST | `/account/reset-password`    |  Reset the password             |

### Friends Related

---

| Type   | Route                                   |    Description                               
| -------| ----------------------------------------|----------------------------------------------|
| GET    | `/account/me/friends`                   |  Get current users friends                   |
| GET    | `/account/me/pending`                   |  Get current users pending friend list       |
| POST   | `/account/{member-id}/friend`           |  Send friend request                         |
| DELETE | `/account/{member-id}/friend`           |  Remove from friend                          |
| POST   | `/account/{member-id}/friend/accept`    |  Accept friend request                       |
| POST   | `/account/{member-id}/friend/cancel`    |  Cancel friend request                       |

### Channels Related

---

| Type   | Route                               |    Description                               
| -------| ------------------------------------|--------------------------------
| GET    | `/channels/me/dm`                   | Get users DMs			        |
| PUT    | `/channels/{channel-id}`            | Edit channel					|
| POST   | `/channels/{channel-id}/dm`         | Get or create DM				|
| GET    | `/channels/{channel-id}/members`    | Get members of given channel   |
| GET    | `/channels/{guild-id}`    		   | Get guild channels             |
| POST   | `/channels/{guild-id}`    		   | Create channel                 |
| DELETE | `/channels/{id}`    				   | Delete channel                 |
| DELETE | `/channels/{id}/dm`    			   | Close DM                       |
