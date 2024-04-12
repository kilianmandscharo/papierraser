function main() {
  const gameId = initGameId();

  const socket = new WebSocket(`ws://localhost:8080/ws?id=${gameId}`);

  const contentContainer = document.getElementById("content");

  socket.onopen = () => {
    console.log("connected");
  };

  socket.onclose = () => {
    console.log("disconnected");
  };

  socket.addEventListener("message", (evt) => {
    const payload = JSON.parse(evt.data);

    if (!payload.type || !payload.data) {
      console.error("ERROR: wrong message format");
      return;
    }

    switch (payload.type) {
      case "Lobby":
        handleLobbyUpdate(contentContainer, payload.data, gameId, socket);
        break;
      case "Track":
        handleTrackUpdate(contentContainer, payload.data, socket);
        break;
      default:
        break;
    }
  });
}

function handleTrackUpdate(contentContainer, html, socket) {
  contentContainer.innerHTML = html;
  connectOptions(
    ".starting-position-option",
    "ActionChooseStartingPosition",
    socket,
  );
  connectOptions(".player-option", "ActionMakeMove", socket);
}

function connectOptions(className, actionType, socket) {
  document.querySelectorAll(className)?.forEach((option) => {
    const [_, rest] = option.id.split("-");
    const [x, y] = rest.split(",");
    option.addEventListener("click", () => {
      if (socket && option.className.baseVal.length > 0) {
        const payload = newPayload(actionType, {
          x: parseInt(x),
          y: parseInt(y),
        });
        socket.send(JSON.stringify(payload));
      }
    });
  });
}

function handleLobbyUpdate(contentContainer, html, gameId, socket) {
  contentContainer.innerHTML = html;

  const inviteButton = document.getElementById("invite-button");
  if (inviteButton) {
    inviteButton.addEventListener("click", () => {
      navigator.clipboard.writeText(`${document.location.host}?id=${gameId}`);
    });
  }

  const changeNameButton = document.getElementById("change-name-button");
  const changeNameInput = document.getElementById("change-name-input");
  if (changeNameButton && changeNameInput) {
    changeNameButton.addEventListener("click", () => {
      const value = changeNameInput.value;
      if (value.length === 0 || value.length > 20) return;
      const payload = newPayload("ActionNameChange", value);
      socket.send(JSON.stringify(payload));
      changeNameInput.value = "";
    });
  }

  const startButton = document.getElementById("start-button");
  if (startButton) {
    startButton.addEventListener("click", () => {
      const payload = newPayload("ActionToggleReady", "");
      socket.send(JSON.stringify(payload));
    });
  }
}

function newPayload(type, data) {
  return { type, data };
}

function getGameId() {
  return localStorage.getItem("gameId");
}

function setGameId(id) {
  localStorage.setItem("gameId", id);
}

function initGameId() {
  let gameId = new URLSearchParams(window.location.search)["id"];
  if (gameId) {
    return gameId;
  }
  gameId = getGameId();
  if (!gameId) {
    gameId = crypto.randomUUID();
    setGameId(gameId);
  }
  return gameId;
}

main();
