The app service accepts a request for creating a user and asynchronously launches requests to microservices (users, cash) in goroutines, gradually going through the stages of preparing and then committing. When all services are committed or an error occurs, the request is completed.

The requests operate synchronously and support the transmission of HTTP status codes and error messages from the microservice to the app.

# Positive Path

POST http://127.0.0.1:8001/create
```
{
    "email": "test@test.com"
}
```

Response:

```
{
  "success": true,
  "uuid": "9c45856b-13c5-46f4-8b69-0a333d1b885e"
}
```

stdout:

```
Cash creation (preparing) successfully: 9c45856b-13c5-46f4-8b69-0a333d1b885e
User creation (preparing) successfully: 9c45856b-13c5-46f4-8b69-0a333d1b885e
Cash creation (commit) successfully: 9c45856b-13c5-46f4-8b69-0a333d1b885e
User creation (commit) successfully: 9c45856b-13c5-46f4-8b69-0a333d1b885e
```

# Negative path

POST http://127.0.0.1:8001/create
```
{
    "email": "123"
}
```

Response:

```
{
  "message": "Invalid email format",
  "success": false
}
```

stdout:

```
User creation (preparing) failed: 57216741-3db1-4b3c-9226-35530a961a5d
Cash creation (preparing) successfully: 57216741-3db1-4b3c-9226-35530a961a5d
User creation rollback: 57216741-3db1-4b3c-9226-35530a961a5d
Cash creation rollback: 57216741-3db1-4b3c-9226-35530a961a5d
```