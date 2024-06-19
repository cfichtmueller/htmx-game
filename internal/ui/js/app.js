(function () {
  let vw = document.getElementById("vw");
  let vh = document.getElementById("vh");

  vw.setAttribute("value", window.innerWidth);
  vh.setAttribute("value", window.innerHeight);

  let pid = document.querySelector("[name=playerId]").value;

  function notifyPlayer(...payload) {
    return fetch("/player/" + pid, {
      method: "POST",
      body: JSON.stringify({
        commands: payload,
      }),
    });
  }

  const setVelocity = (v) => ({ m: "setVelocity", v });
  const setRotation = (v) => ({ m: "setRotation", v });

  document.addEventListener("keydown", function (e) {
    console.log("key down", e);
    switch (e.code) {
      case "KeyA":
        return notifyPlayer(setVelocity(1), setRotation(Math.PI));
        s;
      case "KeyD":
        return notifyPlayer(setVelocity(1), setRotation(2 * Math.PI));
      case "KeyS":
        return notifyPlayer(setVelocity(1), setRotation(Math.PI / 2));
      case "KeyW":
        return notifyPlayer(setVelocity(1), setRotation((3 * Math.PI) / 2));
    }
  });

  document.addEventListener("keyup", function (e) {
    console.log("key up", e);
    switch (e.code) {
      case "KeyA":
      case "KeyD":
      case "KeyS":
      case "KeyW":
        return notifyPlayer(setVelocity(0));
    }
  });

  document.addEventListener("keypress", function (e) {
    let dx = 0;
    let dy = 0;
    switch (e.key) {
      case "s":
        dy = 1;
        break;
      case "a":
        dx = -1;
        break;
      case "d":
        dx = 1;
        break;
      case "p":
        notifyPlayer({ m: "respawn" }).then(() => window.location.assign("/"));
        return;
      default:
        return;
    }
  });
})();
