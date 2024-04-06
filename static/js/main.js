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

const gameId = initGameId();

const socket = new WebSocket(`ws://localhost:8080/ws?id=${getGameId()}`);

const content = document.getElementById("content");

socket.onopen = () => {
  console.log("connected");
};

socket.onclose = () => {
  console.log("disconnected");
};

socket.addEventListener("message", (evt) => {
  content.innerHTML = evt.data;
});

// function connectOptions() {
//   document.querySelectorAll(".player-option")?.forEach((option) => {
//     const [_, rest] = option.id.split("-");
//     const [x, y] = rest.split(",");
//     option.addEventListener("click", () => {
//       if (socket) {
//         socket.send(JSON.stringify({ x: parseInt(x), y: parseInt(y) }));
//       }
//     });
//   });
// }
