import { sizeToStr } from "./utils/size.js"

const $input = document.getElementById("input")
const $list = document.getElementById("list")
const $stat = document.getElementById("stat")

export function renderStat(stat) {
    let itext = `stuff count: ${stat.count}, storage: ${sizeToStr(stat.size)}`
    $stat.innerText = itext
}

export function renderStuff(stuffs, actions) {
    $list.innerText = ''
    stuffs.forEach((id, stf) => {
        if (stf.is_del) {
            return
        }

        const li = document.createElement("li")
        const span = document.createElement("span")
        span.textContent = stf.raw

        let delBtn = document.createElement("button")
        delBtn.type = "button"
        delBtn.textContent = "x"
        delBtn.title = "delete"

        delBtn.onclick = () => {
            console.log("delete button clicked")
            actions.delete(stf, id)
        }

        li.append(span, delBtn)
        $list.append(li)
    })
}


export function render(state, actions) {
    renderStat(state.stat)
    renderStuff(state.stuffs, actions)
}
