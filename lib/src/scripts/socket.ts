import { ClientAction, ServerAction } from "./message";

export function initSocket(
  gameId: string,
  contentContainer: HTMLElement | null,
) {
  if (!contentContainer) {
    return;
  }

  const socket = new WebSocket(`ws://localhost:8080/ws?id=${gameId}`);

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
      case ClientAction.ClientActionLobby:
        handleLobbyUpdate(contentContainer, payload.data, gameId, socket);
        break;
      case ClientAction.ClientActionTrack:
        handleTrackUpdate(contentContainer, payload.data, socket);
        break;
      case ClientAction.ClientActionMove:
        handleMove(JSON.parse(payload.data), socket);
        break;
      default:
        break;
    }
  });
}

function handleMove(
  data: { point: { x: number; y: number }; playerId: number },
  socket: WebSocket,
) {
  document.querySelectorAll(".player-option")?.forEach((opt) => {
    opt.remove();
  });

  const { point, playerId } = data;

  const player = document.getElementById(
    `player-${playerId}`,
  ) as SVGCircleElement | null;

  const canvas = document.getElementById("canvas") as SVGSVGElement | null;

  if (!player || !canvas) {
    return;
  }

  var path = document.createElementNS("http://www.w3.org/2000/svg", "line");
  const playerX = player.getAttribute("cx") ?? "";
  const playerY = player.getAttribute("cy") ?? "";
  const playerColor = player.getAttribute("fill") ?? "";
  console.log(playerX, playerY, playerColor);
  path.setAttribute("x1", playerX);
  path.setAttribute("y1", playerY);
  path.setAttribute("x2", playerX);
  path.setAttribute("y2", playerY);
  path.setAttribute("stroke", playerColor);
  canvas.appendChild(path);

  animatePlayerToPosition(player, path, point.x * 5, point.y * 5, 500, () => {
    socket.send(newPayload(ServerAction.ServerActionMoveAnimationDone));
  });
}

function animatePlayerToPosition(
  player: SVGCircleElement,
  path: SVGLineElement,
  newX: number,
  newY: number,
  duration: number,
  callback: () => void,
) {
  const start = performance.now();
  const startX = parseFloat(player.getAttribute("cx") || "0");
  const startY = parseFloat(player.getAttribute("cy") || "0");

  function step(timestamp: number) {
    const progress = Math.min((timestamp - start) / duration, 1);
    const x = (startX + (newX - startX) * progress).toString();
    const y = (startY + (newY - startY) * progress).toString();
    player.setAttribute("cx", x);
    player.setAttribute("cy", y);
    path.setAttribute("x2", x);
    path.setAttribute("y2", y);

    if (progress < 1) {
      requestAnimationFrame(step);
    } else {
      callback();
    }
  }

  requestAnimationFrame(step);
}

function handleTrackUpdate(
  contentContainer: HTMLElement,
  html: string,
  socket: WebSocket,
) {
  contentContainer.innerHTML = html;
  connectOptions(
    ".starting-position-option",
    ServerAction.ServerActionChooseStartingPosition,
    socket,
  );
  connectOptions(".player-option", ServerAction.ServerActionMakeMove, socket);
}

function handleLobbyUpdate(
  contentContainer: HTMLElement,
  html: string,
  gameId: string,
  socket: WebSocket,
) {
  contentContainer.innerHTML = html;

  const inviteButton = document.getElementById("invite-button");
  if (inviteButton) {
    inviteButton.addEventListener("click", () => {
      navigator.clipboard.writeText(`${document.location.host}?id=${gameId}`);
    });
  }

  const changeNameButton = document.getElementById("change-name-button");
  const changeNameInput = document.getElementById(
    "change-name-input",
  ) as HTMLInputElement | null;
  if (changeNameButton && changeNameInput) {
    changeNameButton.addEventListener("click", () => {
      const value = changeNameInput.value;
      if (value.length === 0 || value.length > 20) return;
      const payload = newPayload(ServerAction.ServerActionNameChange, value);
      socket.send(payload);
      changeNameInput.value = "";
    });
  }

  const startButton = document.getElementById("start-button");
  if (startButton) {
    startButton.addEventListener("click", () => {
      const payload = newPayload(ServerAction.ServerActionToggleReady);
      socket.send(payload);
    });
  }
}

function connectOptions(
  className: string,
  actionType: ServerAction,
  socket: WebSocket,
) {
  document.querySelectorAll(className)?.forEach((opt) => {
    const option = opt as SVGCircleElement;
    const [_, rest] = option.id.split("-");
    const [x, y] = rest.split(",");
    option.addEventListener("click", () => {
      if (socket && option.classList.length > 0) {
        const payload = newPayload(actionType, {
          x: parseInt(x),
          y: parseInt(y),
        });
        socket.send(payload);
      }
    });
  });
}

function newPayload(type: string, data?: any) {
  return JSON.stringify({ type, data });
}
