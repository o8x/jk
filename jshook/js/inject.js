(function () {
    function connect() {
        let wsUrl = "{{ws_url}}"
        console.info("try connecting", wsUrl)

        let s = new WebSocket(wsUrl)
        s.onmessage = e => {
            try {
                let data = JSON.parse(e.data)
                if (data.method === "exec_js") {
                    s.send(JSON.stringify({
                        "name": data.name,
                        "method": data.method,
                        "result": JSON.stringify(eval(data.code)),
                    }))
                }
            } catch (e) {
                console.error("exec js failed", e)
            }
        }

        s.onopen = () => console.log("connection opened")
        s.onerror = e => console.error(e.error)
        s.onclose = () => {
            s.close()
            setTimeout(connect, 5000)
        }
    }

    connect()
}())


