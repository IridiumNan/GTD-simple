import { captureCreate, captureAll, commitDel } from "./capture_request.js";
import { render } from "./render.js"
import * as stf from "./stuff.js"
import * as global from "./global.js"


const $input = document.getElementById('input');
const $list = document.getElementById('list');
const $stat = document.getElementById("stat")

async function init() {
    var stuffArr
    stuffArr, global.stat = await captureAll()
    console.log(`stat on init function ${JSON.stringify(stat)}`)

    console.log(`stuffArr: ${JSON.stringify(stuffArr)}`)

    if (stuffArr == null) {
        console.log("no stuff found on server")
        render()
        return
    }

    stuffArr.forEach((s, idx) => {
        global.stuffs.set(s, s.id)
    })

    render()
}


// enter to capture
document.getElementById('form').addEventListener('submit', e => {
    e.preventDefault();

    // try to send request to server, then get id from server
    const text = global.$input.value.trim();
    if (!text) return;

    var id
    id, global.stat = captureCreate(text)


    if (id != null) global.stuffs.set(stf.setId(id, stf.newStuff(text)), id)

    $input.value = '';
    render();
    // $input.focus();
});

document.getElementById("commit").addEventListener("click", e => {
    e.preventDefault()


    if (global.delArr.length == 0) {
        console.log("there is no stuff marked as deleted, skip request")
        return
    }

    global.stat = commitDel(delArr)
})

init()
global.$input.focus();
