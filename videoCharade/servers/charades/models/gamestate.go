package models

import (
	"math/rand"
	"strings"
	"time"
)

// GameRequest is a struct holding the ID of two players
// And the guessed word
type GameRequest struct {
	FirstUserID  int64
	SecondUserID int64
	Guess        string
}

// GameState is a struct representing the game state for a game
type GameState struct {
	ActorID    int64     `json:"actorID"`
	GuesserID  int64     `json:"guesserID"`
	TimeStart  time.Time `json:"timeStart"`
	NumGuessed int       `json:"numGuessed"`
	NumPlayed  int       `json:"numPlayed"`
	NumRight   int       `json:"numRight"`
	NumSkipped int       `json:"numSkipped"`
	CurrWord   string    `json:"currWord"`
}

// NewGameState returns a new GameState with the given actor and guesser IDs
func NewGameState(actorID int64, guesserID int64) *GameState {
	return &GameState{
		ActorID:    actorID,
		GuesserID:  guesserID,
		TimeStart:  time.Now(),
		NumGuessed: 0,
		NumPlayed:  0,
		NumRight:   0,
		NumSkipped: 0,
		CurrWord:   getWord(),
	}
}

// Skip skips the current word, incrementing total number played and skipped
// Returns the new word
func (gs *GameState) Skip() {
	gs.NumSkipped++
	gs.NumPlayed++
	gs.NewWord()
}

// Guess checks whether the guessed word is correct, updating the word
// and updating the gamestate numbers if it is
func (gs *GameState) Guess(guess string) bool {
	found := false
	if strings.ToLower(guess) == strings.ToLower(gs.CurrWord) { // correct
		found = true
	}
	if found { // correct, set new word, increment num right
		gs.NewWord()
		gs.NumRight++
	}
	gs.NumGuessed++
	gs.NumPlayed++
	return found
}

// NewWord sets a new word for the actor to act
func (gs *GameState) NewWord() {
	gs.CurrWord = getWord()
}

// EndGame ends the current game, returning the results
func (gs *GameState) EndGame() Results {
	return Results{
		ActorID:    gs.ActorID,
		GuesserID:  gs.GuesserID,
		NumGuessed: gs.NumGuessed,
		NumPlayed:  gs.NumPlayed,
		NumRight:   gs.NumRight,
		NumSkipped: gs.NumSkipped,
	}
}

// getWord returns a random word
func getWord() string {
	wordSlice := strings.Split(words, ",")
	return wordSlice[rand.Intn(len(wordSlice))]
}

// Results is a struct for game results
type Results struct {
	ActorID    int64 `json:"actorID"`
	GuesserID  int64 `json:"guesserID"`
	NumGuessed int   `json:"numGuessed"`
	NumPlayed  int   `json:"numPlayed"`
	NumRight   int   `json:"numRight"`
	NumSkipped int   `json:"numSkipped"`
}

// GuessResult is a struct representing whether the guess
// was correct, holding the game state
type GuessResult struct {
	Correct bool       `json:"correct"`
	Guessed string     `json:"guessed"`
	State   *GameState `json:"GameState"`
}

const words = "Airplane,Ears,Piano,Angry,Elephant,Pinch,Baby,Fish,Reach,Ball,Flick,Remote,Baseball,Football,Roll,Basketball,Fork,Sad,Bounce,Giggle,Scissors,Cat,Golf,Skip,Chicken,Guitar,Sneeze,Chimpanzee,Hammer,Spin,Clap,Happy,Spoon,Cough,Horns,Stomp,Cry,Joke,Stop,Dog,Mime,Tail,Drink,Penguin,Toothbrush,Drums,hone,Wiggle,Duck,Photographer"
