package main

// Go dev tools: https://golang.org/dl/
// download source using go get github.com/mattweidner/gohangman
// compile using: go build github.com/mattweidner/gohangman
import (
	"bufio"
	"fmt"
	"github.com/mattweidner/gohangman/newline"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

type State struct {
	Won, Lost                   bool
	Wordlist                    []string
	Wins, Losses, Word          int
	WordState, IncorrectLetters []byte
}

func drawBoard(s *State) {
	for _, x := range s.WordState {
		fmt.Print(string(x))
	}
	fmt.Print(newline.STR)
	fmt.Print("Wrong Letters: ")
	for _, x := range s.IncorrectLetters {
		fmt.Print(string(x) + " ")
	}
	fmt.Print(newline.STR + newline.STR)
	fmt.Print("Guess a letter: ")
}

func playerInput(s *State) {
	r := bufio.NewReader(os.Stdin)
	l, _ := r.ReadString('\n')
	if !((l[0] > 0x40 && l[0] < 0x5b) || (l[0] > 0x60 && l[0] < 0x7b)) {
		fmt.Println(newline.STR + "*** Please only use alphabetic characters.")
		return
	}
	c := strings.ToLower(string(l[0]))
	if contains(s.IncorrectLetters, byte(c[0])) {
		fmt.Println(newline.STR + "*** You already guessed that letter!")
		return
	}
	if !strings.Contains(s.Wordlist[s.Word], c) {
		// word DOES NOT contain the guessed letter.
		s.IncorrectLetters = append(s.IncorrectLetters, byte(c[0]))
		fmt.Println(newline.STR + "*** Incorrect!")
		return
	} else {
		// word CONTAINS the guessed letter.
		for i, a := range s.Wordlist[s.Word] {
			if a == rune(c[0]) {
				s.WordState[i] = byte(a)
			}
		}
	}
}

func initWordState(s *State) {
	// Select a random word
	s.Word = getRandomWord(s.Wordlist)
	s.WordState = []byte{}
	for i := 0; i < len(s.Wordlist[s.Word]); i++ {
		s.WordState = append(s.WordState, '_')
	}
}

func getRandomWord(wordlist []string) int {
	// Seed PRNG using current time.
	rand.Seed(time.Time.UnixNano(time.Now()))
	return rand.Intn(len(wordlist))
}

func initState() State {
	// Initialize game state
	state := State{}
	// Read wordlist and add to game state
	f, e := os.Open("wordlist.txt")
	if e != nil {
		panic(e)
	}
	r := bufio.NewReader(f)
	for w, e := r.ReadString('\n'); e != io.EOF; w, e = r.ReadString('\n') {
		if w[len(w)-2] == '\r' { // Windows style newlines
			state.Wordlist = append(state.Wordlist, w[0:len(w)-2])
		} else { // Everyone else's style of newlines
			state.Wordlist = append(state.Wordlist, w[0:len(w)-1])
		}
	}
	initWordState(&state)
	return state
}

func contains(b []byte, v byte) bool {
	for _, x := range b {
		if x == v {
			return true
		}
	}
	return false
}

func checkWinLoss(s *State) {
	if len(s.IncorrectLetters) > 5 {
		fmt.Println("Sorry, you lost!")
		fmt.Println("The correct answer was:", s.Wordlist[s.Word])
		s.Losses++
		s.Lost = true
		return
	}
	if !contains(s.WordState, byte('_')) {
		fmt.Println("YOU WON!!")
		fmt.Println("You correctly guessed the word:", s.Wordlist[s.Word])
		s.Wins++
		s.Won = true
	}
}

func main() {
	build := 157
	fmt.Printf("Go Hangman by Matt Weidner Build %d"+newline.STR, build)
	fmt.Println("Feb. 2018")
	s := initState()
	for {
		drawBoard(&s)
		playerInput(&s)
		checkWinLoss(&s)
		if s.Won || s.Lost {
			fmt.Printf("Wins: %d\tLosses: %d"+newline.STR, s.Wins, s.Losses)
			fmt.Print("Play again? (Y/N) ")
			r := bufio.NewReader(os.Stdin)
			l, _ := r.ReadString('\n')
			m := string(l[0])
			if m == "y" || m == "Y" {
				s.Won = false
				s.Lost = false
				s.IncorrectLetters = []byte{}
				initWordState(&s)
			} else {
				if s.Losses > s.Wins {
					fmt.Println("Here's your easter egg loser! Can you crack the shell? TGVlZWVyb29vb29vb29vb29vb295eXl5eSBKZWVubm5ubmtpbm5ubm5ubm5zc3Nz")
				} else {
					fmt.Println("Here's your winning easter egg! Can you crack the shell? aHR0cHM6Ly95b3V0dS5iZS9BM1ltSFo5SE1Qcw==")
				}
				return
			}
		}

	}
}
