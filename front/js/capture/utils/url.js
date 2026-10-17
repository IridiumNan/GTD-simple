
export function Join(baseUrl, ...args) {
    let resURL = baseUrl
    console.log(`the length of args: ${args.length}`)

    for (let i = 0; i < args.length; i++) {
        console.log(args[i])

        if (resURL.endsWith("/") && args[i].startsWith("/")) {
            resURL += args[i].substring(1, args[i].length - 1);
            continue
        }

        if (!resURL.endsWith("/") && !args[i].startsWith("/")) {

            resURL += "/" + args[i]
            continue
        }
        resURL += args[i]

        console.log(`the resURL: ${resURL}`)
    }

    return resURL
}
