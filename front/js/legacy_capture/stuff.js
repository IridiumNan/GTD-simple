const newStuff = (raw) => {
    return {
        id: "",
        raw: raw,
        is_del: false,
    }
}
const setId = (id, stuf) => {
    stuf.id = id
    return stuf
}
