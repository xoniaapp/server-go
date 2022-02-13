### Xonia REST API Endpoints
Server URL `https://gateway.xoniaapp.com`

We use sessions for authentication.

### Account Related

---

| Type | Route                        |    Description                                    |Parameter content type|
| ---- | ---------------------------- |-------------------------------------------------- |-----------------|
| GET  | `/account`                   |  Get current user                                 | null |
| PUT  | `/account`                   |  Update current user                              |  |
| PUT  | `/account/change-password`   |  Change current users password                    |  |
| POST | `/account/forgot-password`   |  Forgot password request                          |  |
| POST | `/account/login`             |  User Login                                       |  |
| POST | `/account/logout`            |  User Logout                                      |  |
| POST | `/account/register`          |  Create an account                                |  |
| POST | `/account/reset-password`    |  Reset the password                               |  |

### Friends Related

---

| Type   | Route                                   |    Description                                    |
| -------| ----------------------------------------|-------------------------------------------------- |
| GET    | `/account/me/friends`                   |  Get current users friends                        |
| GET    | `/account/me/pending`                   |  Get current users pending friend list            |
| POST   | `/account/{member-ID}/friend`           |  Send friend request                              |
| DELETE | `/account/{member-ID}/friend`           |  Remove from friend                               |
| POST   | `/account/{member-ID}/friend/accept`    |  Accept friend request                            |
| POST   | `/account/{member-ID}/friend/cancel`    |  Cancel friend request                            |
