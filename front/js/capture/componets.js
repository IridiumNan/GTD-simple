
export function createStuffNode(stf, id, actions) {
    const li = document.createElement('li')

    const span = document.createElement("span")
    span.textContent = stf.raw

    let btn = document.createElement("button")
    btn.type = 'button'
    btn.className = "recoverBtn"
    btn.textContent = 'r'
    btn.onclick = () => {
        console.log("recover button is pressed")
        actions.recover(stf, id)
    }
    li.append(span, btn)

    btn = document.createElement('button')
    btn.type = 'button'
    btn.textContent = 'x'
    btn.className = "deleteBtn"
    btn.onclick = () => {
        console.log("delete button is pressed")
        actions.delete(stf, id)
    }

    li.append(span, btn)



    return li
}
