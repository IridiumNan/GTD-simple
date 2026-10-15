# API

- `/capture/create` POST

Request

> [!NOTE]
> Let the id as empty string so backend will auto generate

```json
{
    "data": {
        "stuff": {
            "id": "",
            "raw": "This is a test stuff",
            "is_del": false
        }
    }
}
```

Response

```json
{
    "msg": "submit success",
    "data": {
        "id": "isfjaoigqrfaoij",

        "stat": {
            "count": 3,
            "size": 12345
        }
    }
}
```

---

- `/capture/delete/{id}` DELETE

Request not contains JSON

Response

In success case

200

```json
{
    "msg": "stuff has been marked is deleted id: fisdfjaifnasdfi ",
    "data": {
        "stat": {
            "count": 3,
            "size": 12345
        }
    }
}
```

In not found case

404

```json
{
    "msg": "stuff not found, this request is omitted, id: fjaiosaoisf ",
    "data": {

        "stat": {
            "count": 3,
            "size": 12345
        }
    }
}
```

---

- `/capture/fetch/{id}` GET

Request not contains JSON

Response

In success case

200

```json
{
    "msg": "stuff found",
    "data": {
        "stuff": {
            "id": "",
            "raw": "This is a test stuff",
            "is_del": false
        }
    }
}
```

In not found case

404

```json
{
    "msg": "stuff not found, id: {id}",
    "data": null
    
}
```

---

- `/capture/clear` DELETE

Request not contains JSON

Response

200

```json
{
    "msg": "all not used stuff has been removed",
    "data": {
        "stat": {
            "count": 3,
            "size": 12345
        }
    }
}
```

---
