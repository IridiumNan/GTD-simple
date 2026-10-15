#!/usr/bin/env bash

baseURL='http://localhost:8080'
id='28ee7b93e1b8876622c0'

echo -e "creating a new stuff with id: $id\n\n"

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
" | jq

curl "$baseURL/capture/all/default" | jq

curl "$baseURL/capture/fetch/$id" | jq

curl "$baseURL/capture/update" \
    -XPOST \
    -H "Content-Type: application/json" \
    -d "{
    \"data\": {
        \"stuff\": {
            \"id\": \"$id\",
            \"raw\": \"Updated stuff content\",
            \"is_del\": false
        }
    }
}
"

curl "$baseURL/capture/delete/$id" -XDELETE | jq

curl "$baseURL/capture/all/default" | jq

echo -e "\n\ncurl for none filter\n\n"

curl "$baseURL/capture/all/none" | jq

echo -e "\n\nclear $id ...\n\n"

curl "$baseURL/capture/clear" -XDELETE | jq

curl "$baseURL/capture/all/none" | jq

curl "$baseURL/capture/flush" -XPOST | jq
