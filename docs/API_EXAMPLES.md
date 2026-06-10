# API Examples

These examples assume the backend is running on `http://localhost:8080`.

## Register

```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "email": "alice@example.com",
    "password": "password123"
  }'
```

Expected response:

```json
{
  "message": "註冊成功！歡迎加入 Mercari Clone！"
}
```

## Login

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "password": "password123"
  }'
```

Expected response:

```json
{
  "token": "<jwt-token>"
}
```

For the protected examples below, store the token:

```bash
TOKEN="<jwt-token>"
```

## List Items

```bash
curl "http://localhost:8080/items?limit=10&offset=0"
```

## Search Items

```bash
curl "http://localhost:8080/items?limit=10&offset=0&search=ps5"
```

## Create Item

```bash
curl -X POST http://localhost:8080/items \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "PlayStation 5",
    "price": 15000
  }'
```

## Buy Item

```bash
curl -X POST http://localhost:8080/items/1/buy \
  -H "Authorization: Bearer $TOKEN"
```

Possible responses:

- `200 OK` when the purchase succeeds.
- `409 Conflict` when the item has already been sold.
- `401 Unauthorized` when the token is missing or invalid.

## Upload Item Image

```bash
curl -X POST http://localhost:8080/items/1/image \
  -H "Authorization: Bearer $TOKEN" \
  -F "image=@/path/to/image.jpg"
```

Expected response:

```json
{
  "image_url": "/uploads/1_image.jpg",
  "message": "圖片上傳成功!"
}
```

The uploaded image is served from:

```text
http://localhost:8080/uploads/1_image.jpg
```
