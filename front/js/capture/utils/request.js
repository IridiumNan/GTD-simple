import * as URL from "./url.js"


let captureURL = "/capture/"

// captureCreate return the created stuff id
export async function captureCreate(raw) {
    let url = URL.Join(captureURL, "create")
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


    return [json.data.id, json.data.stat]

}

export async function captureAll() {
    let url = URL.Join(captureURL, "all", "default")

    console.log(`Request: target url: ${url}`)
    const resp = await fetch(url)

    const json = await resp.json()
    console.log(`status for capture all: ${resp.status}`)
    console.log(json.msg)

    console.log(json.data.stat)

    console.log(JSON.stringify(json.data))

    return [json.data.stuffs, json.data.stat]

}

export async function captureMultiDel(idArr) {
    let url = URL.Join(captureURL, "multiDel")

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

    console.log(`status for capture multi delete: ${resp.status}`)
    console.log(json.msg)
    console.log(json.data.stat)

    return json.data.stat
}

export async function captureDel(id) {
    let url = URL.Join(captureURL, "delete", id)

    const resp = await fetch(url, {
        method: "DELETE"
    })

    const json = await resp.json()

    return json.data.stat
}


// captureClear return the latest stat
export async function captureClear() {
    let url = URL.Join(captureURL, "clear")

    const resp = await fetch(url, {
        method: "DELETE",
    })

    const json = await resp.json()

    console.log(`status for capture clear ${resp.status}`)
    console.log(json.msg)

    return json.data.stat
}
