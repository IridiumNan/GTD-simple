import { captureAll, captureClear, captureCreate, captureDel, captureMultiDel } from "./utils/request.js";
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
    state.stat = await captureMultiDel(state.delArr)

    clearDelArr()

    state.stuffs.clear()
    await initStore()

    console.log(state.stuffs)

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

export async function recoverStuff(stf, id) {
    if (stf.is_del == false) return;
    stf.is_del = false

    let newDelArr = []
    for (let i = 0; i < state.delArr.length; i++) {
        if (state.delArr[i] == id) {
            console.log(`removing the id ${id} from delArr`)
            continue
        }
        newDelArr.push(state.delArr[i])
    }

    state.delArr = newDelArr
}

export async function updateStuff(stf, id) {
    state.stuffs.delete(stf)
    state.stat = await captureDel(id)

    console.log(`deleting stuff on server, id: ${id}`)
}
