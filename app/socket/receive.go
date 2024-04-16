package socket

import (
	"github.com/gorilla/websocket"
	"github.com/kilianmandscharo/papierraser/enum"
	"github.com/kilianmandscharo/papierraser/game"
	"github.com/kilianmandscharo/papierraser/state"
)

func handleReceiveActionDisconnectPlayer(ch chan<- state.ActionRequest, gameId string, addr string, conn *websocket.Conn) {
	ch <- state.ActionRequest{
		GameId: gameId,
		UpdateFunc: func(race *game.Race) state.RenderFunc {
			race.DisconnectPlayer(addr)
			return func(target *game.Player) (enum.ClientAction, []byte) {
				if race.GamePhase == enum.GamePhaseStarted {
					return drawRace(race, target)
				}
				return drawLobby(race, target)
			}
		},
	}
	conn.Close()
}

func handleReceiveActionConnectPlayer(ch chan<- state.ActionRequest, gameId string, addr string, conn *websocket.Conn) {
	ch <- state.ActionRequest{
		GameId: gameId,
		UpdateFunc: func(race *game.Race) state.RenderFunc {
			race.ConnectPlayer(addr, conn)
			return func(target *game.Player) (enum.ClientAction, []byte) {
				// if race.Started {
				// 	return renderTrack(race, target)
				// }
				// return renderLobby(race, target)
				return drawRace(race, target)
			}
		},
	}
}

func handleReceiveActionNameChange(ch chan<- state.ActionRequest, gameId string, addr string, message messageReceive) {
	ch <- state.ActionRequest{
		GameId: gameId,
		UpdateFunc: func(race *game.Race) state.RenderFunc {
			race.UpdatePlayerName(addr, message.Data.(string))
			return func(target *game.Player) (enum.ClientAction, []byte) {
				return drawLobby(race, target)
			}
		},
	}
}

func handleReceiveActionToggleReady(ch chan<- state.ActionRequest, gameId string, addr string) {
	ch <- state.ActionRequest{
		GameId: gameId,
		UpdateFunc: func(race *game.Race) state.RenderFunc {
			race.TogglePlayerReady(addr)
			race.StartPregameIfReady()
			return func(target *game.Player) (enum.ClientAction, []byte) {
				if race.AllPlayersReady() {
					return drawRace(race, target)
				}
				return drawLobby(race, target)
			}
		},
	}
}

func handleReceiveActionChooseStartingPosition(ch chan<- state.ActionRequest, gameId string, addr string, message messageReceive) {
	ch <- state.ActionRequest{
		GameId: gameId,
		UpdateFunc: func(race *game.Race) state.RenderFunc {
			race.UpdateStartingPosition(addr, game.CastPoint(message.Data))

			return func(target *game.Player) (enum.ClientAction, []byte) {
				return drawRace(race, target)
			}
		},
	}
}

func handleReceiveActionMakeMove(ch chan<- state.ActionRequest, gameId string, message messageReceive) {
	ch <- state.ActionRequest{
		GameId: gameId,
		UpdateFunc: func(race *game.Race) state.RenderFunc {
			playerToMove := race.CurrentPlayer()
			movedTo, hasMoved := race.MakeMove(game.CastPoint(message.Data))

			return func(target *game.Player) (enum.ClientAction, []byte) {
				if hasMoved {
					return drawNewPosition(playerToMove, movedTo)
				}
				return drawRace(race, target)
			}
		},
	}
}

func handleReceiveActionMoveAnimationDone(ch chan<- state.ActionRequest, gameId string) {
	ch <- state.ActionRequest{
		GameId: gameId,
		UpdateFunc: func(race *game.Race) state.RenderFunc {
			return func(target *game.Player) (enum.ClientAction, []byte) {
				return drawRace(race, target)
			}
		},
	}
}
