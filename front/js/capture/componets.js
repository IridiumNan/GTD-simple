
export function createStuffNode(stf, id, actions) {
    const li = document.createElement('li')

    const span = document.createElement("span")
    span.textContent = stf.raw

    const btn = document.createElement('button')
    btn.type = 'button'
    btn.textContent = 'x'
    btn.onclick = () => {
        console.log("delete button is pressed")
        actions.delete(stf, id)
    }

    li.append(span, btn)

    return li
}
