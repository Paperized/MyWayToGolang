package main

import "math/rand"

const gameNeverInitialized string = "NEVER"
const gameInitializedCode string = "INITIALIZED"
const gameInProgressCode string = "PROGRESS"
const gameEndedCode string = "ENDED"

const Player1Code string = "P1"
const Player2Code string = "P2"

type PlayerState struct {
	name  string
	score uint8
	isBot bool
}

type GameEvents struct {
	// string -> P1 | P2
	onPlayerInput func(string, *Game) uint8
	onPlayerScore func(string, uint8, *Game)
	onPlayerFail  func(string, uint8, *Game)
	onPlayerWon   func(string, *Game)
}

type Game struct {
	player1 PlayerState
	player2 PlayerState
	state   string
	turn    string
	events  GameEvents
}

func MakeGame(events GameEvents) Game {
	return Game{
		state:  gameNeverInitialized,
		events: events,
	}
}

func (g *Game) SetupMatch(player1Name string, player2Name string, isBot bool) {
	g.player1.name = player1Name
	g.player2.name = player2Name
	g.player1.isBot = false
	g.player2.isBot = isBot

	g.player1.score = 0
	g.player2.score = 0

	g.state = gameInitializedCode
	g.turn = Player1Code
}

func (g *Game) PlayMatch() bool {
	if g.state != gameInitializedCode {
		return false
	}

	g.state = gameInProgressCode
	events := g.events

	// match loop, until someone reaches 3
	for g.player1.score < 3 && g.player2.score < 3 {
		turnPlayerState := GetPlayerFromCode(g.turn, g)
		if turnPlayerState == nil {
			return false
		}
		// handle both player and bot input
		choice := _handlePlayerInput(turnPlayerState, g)

		scoreIndex := uint8(rand.Intn(5))
		if choice == scoreIndex {
			turnPlayerState.score++
			events.onPlayerScore(g.turn, scoreIndex, g)
		} else {
			events.onPlayerFail(g.turn, scoreIndex, g)
		}

		if g.turn == Player1Code {
			g.turn = Player2Code
		} else if g.turn == Player2Code {
			g.turn = Player1Code
		} else {
			return false
		}
	}

	winner := Player1Code
	if g.player2.score == 3 {
		winner = Player2Code
	}

	events.onPlayerWon(winner, g)
	g.state = gameEndedCode
	return true
}

func _handlePlayerInput(playerState *PlayerState, g *Game) uint8 {
	if !playerState.isBot {
		return g.events.onPlayerInput(g.turn, g)
	}

	return _botPlayerInput()
}

func _botPlayerInput() uint8 {
	// just a random function
	return uint8(rand.Intn(5))
}

func GetPlayerFromCode(label string, g *Game) *PlayerState {
	if label == Player1Code {
		return &g.player1
	} else if label == Player2Code {
		return &g.player2
	}

	return nil
}
