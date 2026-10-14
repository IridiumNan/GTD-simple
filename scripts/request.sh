#!/usr/bin/env bash

baseURL='http://localhost:8080'

curl "$baseURL/capture/create" \
    -XPOST \
    -H "Content-Type: application/json" \
    -d '{
    "data": {
        "stuff": {
            "id": "",
            "raw": "This is test for creating",
            "is_del": false
        }
    }
}
'

id='d9660cdafa2d8a280ba'
curl "$baseURL/capture/delete/$id" \
    -XDELETE
