import { captureAll, captureClear, captureCreate, captureDel } from "./utils/request.js";
import { newStuff, setId } from "./model/stuff.js";

// private state
const state = {
    stuffs: new Map(),
    delArr: [],
    stat: {
        count: 0,
        size: 0
    }
}

export function getState() {
    return state;
}

export async function initStore() {
    const [stuffArr, stat] = await captureAll();
    // console.log(`stat on init function ${stat}`)
    console.log(`stat on init function: ${JSON.stringify(stat)}`)
    console.log(`stuffArr: ${JSON.stringify(stuffArr)}`)

    if (stat) {
        state.stat = stat;
    }

    if (stuffArr == null) {
        console.log("no stuff found on server")
        return false
    }

    stuffArr.forEach((s) => {
        state.stuffs.set(s, s.id)
    })

    return true
}

export function deleteStuff(stf, id) {
    stf.is_del = true
    state.delArr.push(id)
}

export async function deleteCommit() {
    if (!state.delArr || state.delArr.length == 0) {
        return `there is no stuff on del arrary, skip this commit`
    }
    state.stat = await captureDel(state.delArr)

    clearDelArr()

    return `delete stuffs commited`
}

export function clearDelArr() {
    state.delArr = []
}

export async function createStuff(raw) {
    var id
    [id, state.stat] = await captureCreate(raw)


    if (id == null) {
        console.log("get null id, drop this stuff")
        return
    }

    let stf = newStuff(raw)
    setId(id, stf)
    state.stuffs.set(stf, id)
}

export async function clearDelStuff() {
    state.stat = await captureClear()
}
