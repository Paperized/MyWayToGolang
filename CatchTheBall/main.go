package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"time"
)

func main() {
	game := MakeGame(GameEvents{
		onPlayerInput: onPlayerInput,
		onPlayerScore: onPlayerScore,
		onPlayerFail:  onPlayerFail,
		onPlayerWon:   onPlayerWon,
	})

	printTitle()
	printInstructions()

	readInputsAndSetupMatch(&game)
	res := game.PlayMatch()
	if !res {
		// would have been better a string to print the actual error, but it's just a learning project
		fmt.Println("Game terminated because of an error")
	}
}

func onPlayerInput(playerCode string, game *Game) uint8 {
	clearConsole()

	printTitle()
	printGameStats(game)
	printHiddenChoices()

	p := GetPlayerFromCode(playerCode, game)
	// in case of bot we updated only the visual and simulated a wait
	if p.isBot {
		fmt.Print("Opponent thinking")
		for i := 0; i < 3; i++ {
			time.Sleep(700 * time.Millisecond) // wait for a total of 2.1 seconds, append a dot in the console
			fmt.Print(".")
		}
		fmt.Print("\n")
		return 0
	}

	var input string
	for readStdinTrim(&input); !isValidChoice(input); readStdinTrim(&input) {
	}
	choice, _ := strconv.Atoi(input)
	return uint8(choice)
}

func onPlayerScore(playerCode string, scoreIndex uint8, game *Game) {
	clearConsole()

	printTitle()
	printGameStats(game)
	printShownChoices(scoreIndex)

	p := GetPlayerFromCode(playerCode, game)
	fmt.Printf("%v Scored!\n", p.name)
	time.Sleep(time.Second)

	var _unused string
	readStdinTrim(&_unused)
}

func onPlayerFail(playerCode string, scoreIndex uint8, game *Game) {
	clearConsole()

	printTitle()
	printGameStats(game)
	printShownChoices(scoreIndex)

	p := GetPlayerFromCode(playerCode, game)
	fmt.Printf("%v Failed!\n", p.name)
	time.Sleep(time.Second)

	var _unused string
	readStdinTrim(&_unused)
}

func onPlayerWon(playerCode string, game *Game) {
	clearConsole()

	printTitle()
	printGameStats(game)

	p := GetPlayerFromCode(playerCode, game)
	fmt.Printf("%v won!! Congratulation!\n", p.name)

	var _unused string
	readStdinTrim(&_unused)
}

func readInputsAndSetupMatch(game *Game) {
	var input string
	var isBot1 bool = false
	var isBot2 bool = false
	var player1Name string
	var player2Name string

	fmt.Println("Choose 1) Play with a friend 2) Play vs Bot 3) Bot vs Bot :')")

	for readStdinTrim(&input); input != "1" && input != "2" && input != "3"; readStdinTrim(&input) {
	}
	if input == "2" {
		isBot1 = false
		isBot2 = true
	} else if input == "3" {
		isBot1 = true
		isBot2 = true
	}

	fmt.Println("Insert player 1 name")
	for readStdinTrim(&player1Name); player1Name == ""; readStdinTrim(&player1Name) {
	}

	fmt.Println("Insert player 2 name")
	for readStdinTrim(&player2Name); player2Name == ""; readStdinTrim(&player2Name) {
	}

	game.SetupMatch(player1Name, player2Name, isBot1, isBot2)
}

func isValidChoice(input string) bool {
	switch input {
	case "1", "2", "3", "4", "5":
		return true
	default:
		return false
	}
}

func readStdinTrim(output *string) {
	fmt.Print("> ")
	fmt.Scanf("%s\n", output)

	*output = strings.Trim(*output, " ")
}

func printGameStats(game *Game) {
	fmt.Printf("%s %-4v %4v %s\n", game.player1.name, game.player1.score, game.player2.score, game.player2.name)
}

func printHiddenChoices() {
	fmt.Println(" 1  2  3  4  5")
	fmt.Println(" X  X  X  X  X")
}

func printShownChoices(index uint8) {
	fmt.Println(" 1  2  3  4  5")

	// considering the placeholders XXXXX, set the right index and pad them with 2 spaces between each others
	choices := " XXXXX"
	choicesLen := len(choices)
	// Convert the string to a mutable list of characters and set the correct index
	copyRunes := []rune(choices)
	copyRunes[index] = 'O'

	mainIndex := 1 // consider the first space in choices
	for i := 0; i < choicesLen-1; i++ {
		copyRunes = slices.Insert(copyRunes, mainIndex+1, ' ', ' ')
		mainIndex += 3
	}

	choices = string(copyRunes)
	fmt.Println(choices)
}

func printInstructions() {
	fmt.Println("This game can be player in 2 or with the PC")
	fmt.Println("One have to hide the ball, the other have to find it to score a point")
	fmt.Print("Players will take turns until one reach 3 points\n\n")
}

func printTitle() {
	fmt.Print("< WHERE IS THE BALL? >\n\n")
}

// needs this cross platform code to clear the console each round
func clearConsole() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux", "darwin": // Unix-like systems
		cmd = exec.Command("clear")
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		fmt.Println("Unsupported platform")
		return
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}
