import * as global from "./global.js"
import { sizeToStr } from "./size_utils.js";

const $input = document.getElementById('input');
const $list = document.getElementById('list');
const $stat = document.getElementById("stat")

function renderStat() {
    console.log(`stat on render function: ${JSON.stringify(global.stat)}`)
    let itext = `stuff count: ${global.stat.count}, storage: ${sizeToStr(global.stat.size)}`
    $stat.innerText = itext
}

function renderStuff() {
    $list.innerHTML = '';
    global.stuffs.forEach((id, stf) => {
        if (stf.is_del) return
        const li = document.createElement('li');
        const span = document.createElement('span');
        span.textContent = stf.raw;

        const btn = document.createElement('button');
        btn.type = 'button';
        btn.textContent = '×';
        btn.title = 'delete';
        btn.onclick = () => { stf.is_del = true; delArr.push(id); render(); };

        li.append(span, btn);
        global.$list.append(li);
    })

}

function render() {
    renderStat()
    renderStuff()
}

export { render }
