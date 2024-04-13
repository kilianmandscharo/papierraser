import { initGameId } from "./gameId";
import { initSocket } from "./socket";

function main() {
  const gameId = initGameId();
  const contentContainer = document.getElementById("content");
  initSocket(gameId, contentContainer);
}

main();
