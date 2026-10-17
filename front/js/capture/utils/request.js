const captureURL = "http://localhost:8080/capture/"

// captureCreate return the created stuff id
async function captureCreate(raw) {
    let url = captureURL + "create"
    const bodyData = {
        data: {
            stuff: {
                id: "",
                raw: raw,
                is_del: false,
            }
        }
    }

    let bodyStr = JSON.stringify(bodyData)

    console.log(`Request: target url: ${url}, body: ${bodyStr}`)

    const resp = await fetch(url, {
        method: "POST",
        body: bodyStr
    })

    const json = await resp.json()
    console.log(`status: ${resp.status}`)
    console.log(json.msg)

    if (resp.status != 201) {
        return null
    }


    return json.data.id, json.data.stat

}

async function captureAll() {
    let url = captureURL + "all" + "/" + "default"

    console.log(`Request: target url: ${url}`)
    const resp = await fetch(url)

    const json = await resp.json()
    console.log(`status: ${resp.status}`)
    console.log(`stuffs: ${JSON.stringify(json.data.stuffs)}`)

    return [json.data.stuffs, json.data.stat]

}

async function commitDel(idArr) {
    let url = captureURL + "multiDel"

    const bodyData = {
        data: {
            id_list: idArr
        }
    }
    let bodyStr = JSON.stringify(bodyData)

    console.log(`Request: target url: ${url}, body: ${bodyStr}`)

    const resp = await fetch(url, {
        method: "DELETE",
        body: bodyStr
    })

    const json = await resp.json()

    console.log(`status: ${resp.status}`)
    console.log(json.msg)
    console.log(json.data.stat)

    return json.data.stat
}


export { captureCreate, captureAll, commitDel }


