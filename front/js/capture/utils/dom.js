export function mountText(target, text) {
    const elm = typeof target === 'string' ?
        document.getElementById(target) : target

    if (elm) {
        elm.innerText = text
    }
}

export function mountList(container, children) {
    const parent = typeof container == 'string' ?
        document.getElementById(container) : container

    if (!parent) return

    parent.innerText = ''

    if (Array.isArray(children) && children.length > 0) {
        parent.append(...children)
    }

}
