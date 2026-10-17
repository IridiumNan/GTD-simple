import { initStore, getState, deleteStuff, createStuff, deleteCommit, clearDelStuff } from "./store.js";
import { mountText, mountList } from "./utils/dom.js";
import { createStuffNode } from "./componets.js";
import { sizeToStr } from "./utils/size.js";

const $form = document.getElementById("form")
const $input = document.getElementById('input')
const $commit = document.getElementById('commit')
const $clear = document.getElementById("clear")

const state = getState()

const updateView = () => {
    const statText = `stuff count: ${state.stat.count}, storage: ${sizeToStr(state.stat.size)}`
    mountText('stat', statText)


    const actions = {
        delete: (targetStf, targetId) => {
            deleteStuff(targetStf, targetId)
            updateView()
        }
    }

    const stuffNodes = []
    state.stuffs.forEach((id, stf) => {
        if (stf.is_del) {
            console.log(`this stuff have been deleted, skip it, id: ${id}`)
            return
        }
        const node = createStuffNode(stf, id, actions)
        stuffNodes.push(node)
    })

    console.log(`the length of structure stuff nodes: ${stuffNodes.length}`)
    mountList('list', stuffNodes)
}

$form.addEventListener('submit', async (e) => {
    e.preventDefault()

    const text = $input.value.trim()
    await createStuff(text)
    $input.value = ''

    updateView()

})

$commit.addEventListener('click', async (e) => {
    e.preventDefault()

    let msg = await deleteCommit()
    console.log(msg)
})

$clear.addEventListener("click", async (e) => {
    e.preventDefault()

    await clearDelStuff()
    updateView()
})

// TODO: Write an export button


async function bootstrap() {
    await initStore()

    updateView()
}

bootstrap()

