const KEY = 'gtd-simple'
// const $empty = document.getElementById('empty');


// the key is stuff obj and value is id
let stuffs = new Map();
let delArr = [];
let stat = {
    count: 0,
    size: 0
}

function SetStat(newStat) {
    stat = newStat
}

function SetStuffs(newStuffs) {
    stuffs = newStuffs
}

function PushDelArr(id) {
    delArr.push(id)
}

function ClearDelArr() {
    delArr = []
}

export { stuffs, delArr, stat, SetStat, SetStuffs, PushDelArr, ClearDelArr }
