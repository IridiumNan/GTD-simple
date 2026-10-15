async function create() {
    const url = "http://localhost:8080/capture/create";
    try {
        const response = await fetch(url,
            {
                method: "POST",
                body: JSON.stringify({
                    data: {
                        stuff: {
                            id: null,
                            raw: "This is test stuff",
                            is_del: false
                        }
                    }
                })
            })
        // get response status code
        console.log(response.status)

        const result = await response.json()
        console.log(result)
    } catch (error) {
        console.error(error.message)

    }
}

create()
