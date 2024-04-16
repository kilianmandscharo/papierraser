package enum

type GamePhase string

const (
	GamePhaseLobby    GamePhase = "GamePhaseLobby"
	GamePhasePregame  GamePhase = "GamePhasePregame"
	GamePhaseStarted  GamePhase = "GamePhaseStarted"
	GamePhaseFinished GamePhase = "GamePhaseFinished"
)

type ClientAction string

const (
	ClientActionDrawRace        ClientAction = "ClientActionRenderRace"
	ClientActionDrawLobby       ClientAction = "ClientActionRenderLobby"
	ClientActionDrawNewPosition ClientAction = "ClientActionDrawNewPosition"
)

type ServerAction string

const (
	ServerActionNameChange             ServerAction = "ServerActionNameChange"
	ServerActionToggleReady            ServerAction = "ServerActionToggleReady"
	ServerActionChooseStartingPosition ServerAction = "ServerActionChooseStartingPosition"
	ServerActionMakeMove               ServerAction = "ServerActionMakeMove"
	ServerActionMoveAnimationDone      ServerAction = "ServerActionMoveAnimationDone"
)
