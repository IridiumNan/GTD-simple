# API

- `/capture/create` POST

Request

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
        "id": "isfjaoigqrfaoij"
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
    "data": null
}
```

In not found case

404

```json
{
    "msg": "stuff not found, this request is omitted, id: fjaiosaoisf ",
    "data": null
}
```

---

- `/capture/fetch/{id}` GET

Request not contains JSON

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
