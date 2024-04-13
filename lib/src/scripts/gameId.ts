function getGameId() {
  return localStorage.getItem("gameId");
}

function setGameId(id: string) {
  localStorage.setItem("gameId", id);
}

export function initGameId() {
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

