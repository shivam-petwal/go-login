# Go Login API

## Endpoint

**POST** `/login`

## Request

```json
{
  "email": "test@test.com",
  "password": "pass"
}
```

## Response

**Success (200):**
```json
{
  "token": "eyJhbGciOiiIsInR5cCI6IkpXVCJ9..."
}
```
