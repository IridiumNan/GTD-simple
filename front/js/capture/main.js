import { initStore, getState, deleteStuff } from "./store.js";
import { render } from "./render.js"


async function bootstrap() {
    const hasData = await initStore()
    const state = getState()

    const actions = {
        delete: (stf, id) => {
            deleteStuff(stf, id)
            render(state, actions)
        },
    }

    render(state, actions)
}

bootstrap()
