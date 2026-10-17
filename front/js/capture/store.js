import { captureAll } from "./utils/request.js";

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

export function clearDelArr() {
    state.delArr = []
}

export function createStuff(text) {

}
