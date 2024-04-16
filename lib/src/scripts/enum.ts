export enum GamePhase {
  GamePhaseLobby = "GamePhaseLobby",
  GamePhasePregame = "GamePhasePregame",
  GamePhaseStarted = "GamePhaseStarted",
  GamePhaseFinished = "GamePhaseFinished",
}

export enum ClientAction {
  ClientActionDrawRace = "ClientActionDrawRace",
  ClientActionDrawLobby = "ClientActionDrawLobby",
  ClientActionDrawNewPosition = "ClientActionDrawNewPosition",
}

export enum ServerAction {
  ServerActionNameChange = "ServerActionNameChange",
  ServerActionToggleReady = "ServerActionToggleReady",
  ServerActionChooseStartingPosition = "ServerActionChooseStartingPosition",
  ServerActionMakeMove = "ServerActionMakeMove",
  ServerActionMoveAnimationDone = "ServerActionMoveAnimationDone",
}