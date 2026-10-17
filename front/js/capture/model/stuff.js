export function newStuff(raw) {
    return {
        id: "",
        raw: raw,
        is_del: false,
    }
}
export function setId(id, stuf) {
    stuf.id = id
    return stuf
}
