
export function createStuffNode(stf, id, actions) {
    const li = document.createElement('li')

    const span = document.createElement("span")
    span.textContent = stf.raw
    li.append(span)

    let btn = document.createElement("button")

    if (stf.is_del) {

        // recover button
        btn.type = 'button'
        btn.className = "iconfont icon-huifu"
        btn.title = "recover"
        btn.onclick = () => {
            console.log("recover button is pressed")
            actions.recover(stf, id)
        }
        li.append(btn)
    }



    // update button
    // this button fill the input box with stuff raw then delete this stuff on server

    btn = document.createElement("button")
    btn.type = "button"
    btn.className = "iconfont icon-bianji"
    btn.style.fontSize = "20px"
    btn.title = "edit"
    btn.onclick = () => {
        console.log("updating button is pressed")
        actions.update(stf, id)
    }
    li.append(btn)

    // soft-delete button
    btn = document.createElement('button')
    btn.type = 'button'
    btn.className = "iconfont icon-shanchu"
    btn.title = "delete"
    btn.onclick = () => {
        console.log("delete button is pressed")
        actions.delete(stf, id)
    }


    li.append(btn)


    return li
}
