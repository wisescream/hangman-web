package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	words          = []string{"hangman", "computer", "programming", "golang", "ynov", "coding", "html", "css", "javascript", "python", "java", "ruby", "php", "database", "network", "algorithm", "git", "server", "client", "framework"}
	wordToGuess    string
	lettersFound   []string
	chancesLeft    int
	gameOver       bool
	win            bool
	alphabet       = "abcdefghijklmnopqrstuvwxyz"
	guessedLetters []string
	hangmanParts   = []string{
		"   ____\n  |    |\n       |\n       |\n       |\n       |\n_______|",
		"   ____\n  |    |\n  O    |\n       |\n       |\n       |\n_______|",
		"   ____\n  |    |\n  O    |\n  |    |\n       |\n       |\n_______|",
		"   ____\n  |    |\n  O    |\n /|    |\n       |\n       |\n_______|",
		"   ____\n  |    |\n  O    |\n /|\\   |\n       |\n       |\n_______|",
		"   ____\n  |    |\n  O    |\n /|\\   |\n /     |\n       |\n_______|",
		"   ____\n  |    |\n  O    |\n /|\\   |\n / \\   |\n       |\n_______|",
	}
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("index.html").Funcs(template.FuncMap{"inSlice": inSlice}).ParseFiles("index.html", "menu.html"))
	data := struct {
		ChancesLeft    int
		GameOver       bool
		Win            bool
		Letters        []string
		Alphabet       []string
		GuessedLetters []string
		HangmanPart    string
	}{
		ChancesLeft:    chancesLeft,
		GameOver:       gameOver,
		Win:            win,
		Letters:        lettersFound,
		Alphabet:       strings.Split(alphabet, ""),
		GuessedLetters: guessedLetters,
		HangmanPart:    hangmanParts[6-chancesLeft],
	}
	tmpl.Execute(w, data)
}

func guessHandler(w http.ResponseWriter, r *http.Request) {
	if gameOver {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	r.ParseForm()
	guess := strings.ToLower(r.FormValue("guess"))
	if guess == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if !inSlice(guess, guessedLetters) {
		guessedLetters = append(guessedLetters, guess)
	}

	letterFound := false
	for i, letter := range wordToGuess {
		if string(letter) == guess {
			lettersFound[i] = string(letter)
			letterFound = true
		}
	}

	if !letterFound {
		chancesLeft--
	}

	if chancesLeft <= 0 || strings.Join(lettersFound, "") == wordToGuess {
		gameOver = true
		if strings.Join(lettersFound, "") == wordToGuess {
			win = true
			wordToGuess = ""
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func playAgainHandler(w http.ResponseWriter, r *http.Request) {
	gameOver = false
	win = false
	guessedLetters = nil
	resetGame()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func menuHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("menu.html"))
	tmpl.Execute(w, nil)
}

func resetGame() {
	rand.Seed(time.Now().UnixNano())
	wordToGuess = words[rand.Intn(len(words))]
	lettersFound = make([]string, len(wordToGuess))
	for i := range lettersFound {
		lettersFound[i] = "_"
	}
	chancesLeft = 6
}

func inSlice(needle string, haystack []string) bool {
	for _, value := range haystack {
		if value == needle {
			return true
		}
	}
	return false
}

func main() {
	resetGame()

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/guess", guessHandler)
	http.HandleFunc("/play-again", playAgainHandler)
	http.HandleFunc("/menu", menuHandler)

	fmt.Println("Hangman game is running at http://localhost:8080")
	log.SetOutput(os.Stdout)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
