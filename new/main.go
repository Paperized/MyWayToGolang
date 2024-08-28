package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"strings"
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

	var test string
	readStdinTrim(&test)
	return 1
}

func onPlayerScore(playerCode string, scoreIndex uint8, game *Game) {

}

func onPlayerFail(playerCode string, scoreIndex uint8, game *Game) {

}

func onPlayerWon(playerCode string, game *Game) {

}

func readInputsAndSetupMatch(game *Game) {
	var input string
	var isBot bool
	var player1Name string
	var player2Name string

	fmt.Println("Choose 1) Play with a friend 2) Play vs Bot")

	for readStdinTrim(&input); input != "1" && input != "2"; readStdinTrim(&input) {
	}
	isBot = input == "2"

	fmt.Println("Insert player 1 name")
	for readStdinTrim(&player1Name); player1Name == ""; readStdinTrim(&player1Name) {
	}

	fmt.Println("Insert player 2 name")
	for readStdinTrim(&player2Name); player2Name == ""; readStdinTrim(&player2Name) {
	}

	game.SetupMatch(player1Name, player2Name, isBot)
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
	printTitle()
	fmt.Println(" 1  2  3  4  5")

	choices := " XXXXX"
	// Convert the string to a mutable list of characters and set the correct index
	copyRunes := []rune(choices)
	copyRunes[index] = 'O'

	for i := 2; i < len(copyRunes); i += 2 {
		copyRunes = slices.Insert(copyRunes, i, ' ', ' ')
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
