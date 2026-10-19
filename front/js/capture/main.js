import { initStore, getState, deleteStuff, createStuff, deleteCommit, clearDelStuff, recoverStuff, updateStuff } from "./store.js";
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
        // this action just soft-delete locally
        delete: (targetStf, targetId) => {
            deleteStuff(targetStf, targetId)
            updateView()
        },
        recover: (targetStf) => {
            recoverStuff(targetStf)
            updateView()
        },
        update: (stf, id) => {
            if ($input.value) {
                alert("please commit current content on input box")
                return;
            }

            $input.value = stf.raw
            updateStuff(stf, id)

            updateView()
        }
    }

    const stuffNodes = []
    state.stuffs.forEach((id, stf) => {
        const node = createStuffNode(stf, id, actions)
        if (stf.is_del) {

            node.style.opacity = '0.3'
        }
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
    updateView()
})

$clear.addEventListener("click", async (e) => {
    e.preventDefault()

    await clearDelStuff()
    updateView()
})


async function bootstrap() {
    await initStore()

    updateView()
}

bootstrap()

