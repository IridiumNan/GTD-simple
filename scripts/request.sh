#!/usr/bin/env bash

baseURL='http://localhost:8080'
id='28ee7b93e1b8876622c0'

curl "$baseURL/capture/create" \
    -XPOST \
    -H "Content-Type: application/json" \
    -d "{
    \"data\": {
        \"stuff\": {
            \"id\": \"$id\",
            \"raw\": \"This is test for creating\",
            \"is_del\": false
        }
    }
}
"

curl "$baseURL/capture/fetch/$id"

curl "$baseURL/capture/delete/$id" \
    -XDELETE
