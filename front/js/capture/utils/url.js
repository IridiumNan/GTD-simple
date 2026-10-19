
export function Join(baseUrl, ...args) {
    let resURL = baseUrl

    for (let i = 0; i < args.length; i++) {

        if (resURL.endsWith("/") && args[i].startsWith("/")) {
            resURL += args[i].substring(1, args[i].length - 1);
            continue
        }

        if (!resURL.endsWith("/") && !args[i].startsWith("/")) {

            resURL += "/" + args[i]
            continue
        }
        resURL += args[i]

    }

    return resURL
}
