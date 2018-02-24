package main

// Go dev tools: https://golang.org/dl/
// compile using: go build gohangman.go
import (
	"bufio"
	"fmt"
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
	fmt.Println(s.Wordlist[s.Word])
	fmt.Println(string(s.WordState[:]))
	fmt.Print("Wrong Letters: ")
	for _, x := range s.IncorrectLetters {
		fmt.Print(string(x) + " ")
	}
	fmt.Print("\n\n")
	fmt.Print("Guess a letter: ")
}

func playerInput(s *State) {
	r := bufio.NewReader(os.Stdin)
	l, _ := r.ReadString('\n')
	if !((l[0] > 0x40 && l[0] < 0x5b) || (l[0] > 0x60 && l[0] < 0x7b)) {
		fmt.Println("\n*** Please only use alphabetic characters.")
		return
	}
	c := strings.ToLower(string(l[0]))
	if contains(s.IncorrectLetters, byte(c[0])) {
		fmt.Println("\n*** You already guessed that letter!")
		return
	}
	if !strings.Contains(s.Wordlist[s.Word], c) {
		// word DOES NOT contain the guessed letter.
		s.IncorrectLetters = append(s.IncorrectLetters, byte(c[0]))
		fmt.Println("\n*** Incorrect!")
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
		state.Wordlist = append(state.Wordlist, w[0:len(w)-1])
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
		s.Losses++
		s.Lost = true
		return
	}
	if !contains(s.WordState, byte('_')) {
		s.Wins++
		s.Won = true
	}
}

func main() {
	build := 121
	fmt.Printf("Go Hangman by Matt Weidner Build %d\n", build)
	fmt.Println("Feb. 2018")
	s := initState()
	for {
		drawBoard(&s)
		playerInput(&s)
		checkWinLoss(&s)
		if s.Won || s.Lost {
			fmt.Printf("Wins: %d\tLosses: %d\n", s.Wins, s.Losses)
			fmt.Print("Play again? (Y/N) ")
			r := bufio.NewReader(os.Stdin)
			l, _ := r.ReadString('\n')
			if l == "y\n" || l == "Y\n" {
				s.Won = false
				s.Lost = false
				s.IncorrectLetters = []byte{}
				initWordState(&s)
			} else {
				if time.Time.UnixNano(time.Now())%2 == 0 {
					fmt.Println("TGVlZWVyb29vb29vb29vb29vb295eXl5eSBKZWVubm5ubmtpbm5ubm5ubm5zc3NzZWNobwo=")
				} else {
					fmt.Println("aHR0cHM6Ly93d3cueW91dHViZS5jb20vd2F0Y2g/dj1vSGc1U0pZUkhBMAo=")
				}
				return
			}
		}

	}
}
