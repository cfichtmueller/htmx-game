(function() {
    let vw = document.getElementById("vw")
    let vh = document.getElementById("vh")

    vw.setAttribute("value", window.innerWidth)
    vh.setAttribute("value", window.innerHeight)

    let pid = document.querySelector("[name=playerId]").value
    
    document.addEventListener("keypress", function(e) {
        let dx = 0
        let dy = 0
        switch (e.key) {
            case "w":
                dy = -1
                break
            case "s":
                dy = 1
                break
            case "a":
                dx = -1
                break;
            case "d":
                dx = 1
                break
            case "p":
                fetch("/player/" + pid, {
                    method: "POST",
                    body: JSON.stringify({
                        action: "respawn"
                    })
                }).then(function() {
                    window.location.assign("/")
                })
                return
            default:
                return
        }
        fetch("/player/" + pid, {
            method: "POST",
            body: JSON.stringify({
                action: "move",
                dx,
                dy
            })
        })
    })
})()