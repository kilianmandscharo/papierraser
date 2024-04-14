export enum ClientAction {
  ClientActionLobby = "ClientActionLobby",
  ClientActionTrack = "ClientActionTrack",
  ClientActionMove = "ClientActionMove",
}

export enum ServerAction {
  ServerActionNameChange = "ServerActionNameChange",
  ServerActionToggleReady = "ServerActionToggleReady",
  ServerActionChooseStartingPosition = "ServerActionChooseStartingPosition",
  ServerActionMakeMove = "ServerActionMakeMove",
  ServerActionMoveAnimationDone = "ServerActionMoveAnimationDone",
}
