package handlers

import (
	"database/sql"
	"encoding/json"

	"final-project-crew/videoCharade/servers/charades/models"

	"github.com/streadway/amqp"
)

const sqlInsertUsers = "insert into leaderboards (actorID,guesserID,numGuessRight,numGuessPlayed) values (?,?,?,?)"

// Context is a struct that holds all of the game states,
// the rabbitmq channel and queue
// and a connection to the leaderboards
type Context struct {
	Games   map[int64]*models.GameState // The games
	Db      *sql.DB                     // The database for leaderboard
	Channel *amqp.Channel               // The channel for mq
	Queue   amqp.Queue                  // The queue for mq
}

// GameSignal is a struct holding a userid and game state
type GameSignal struct {
	Type    string      `json:"type"`
	UserIDs []int64     `json:"userList"`
	State   interface{} `json:"data"`
}

// StartSignal is a string representing start
const StartSignal = "game-start"

// EndSignal is a string representing end
const EndSignal = "game-end"

// GuessSignal is a string representing guess
const GuessSignal = "guess"

// SkipSignal is a string representing skip
const SkipSignal = "skip"

// SendStartGameSignal notifies users that the game is starting
func (ctx *Context) SendStartGameSignal(firstID int64, secondID int64, state *models.GameState) {
	stateObj := *state
	firstMessageObj := GameSignal{Type: StartSignal, UserIDs: []int64{firstID}, State: stateObj}
	firstMessage, err := json.Marshal(firstMessageObj)
	if err != nil {
		println("there was an error marshalling first message")
		return
	}
	stateObj.CurrWord = ""
	secondMessageObj := GameSignal{Type: StartSignal, UserIDs: []int64{secondID}, State: stateObj}
	secondMessage, err := json.Marshal(secondMessageObj)
	if err != nil {
		println("There was an error marshalling second message")
		return
	}
	body := firstMessage
	err = ctx.Channel.Publish(
		"",             // exchange
		ctx.Queue.Name, // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		println("err was not nil " + err.Error())
	}
	body = secondMessage
	err = ctx.Channel.Publish(
		"",             // exchange
		ctx.Queue.Name, // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		println("err was not nil " + err.Error())
	}
}

// SendEndGameSignal sends an end game signal to both players
func (ctx *Context) SendEndGameSignal(firstID int64, secondID int64, results models.Results) {
	gameSignal := GameSignal{Type: EndSignal, UserIDs: []int64{firstID, secondID}, State: results}
	message, err := json.Marshal(gameSignal)
	if err != nil {
		println("error marshalling game signal for end game")
		return
	}
	body := message
	err = ctx.Channel.Publish(
		"",             // exchange
		ctx.Queue.Name, // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		println("err was not nil " + err.Error())
	}
}

// SendGuessSignal notifies users that a guess had been made
// Telling whether the guess was right
func (ctx *Context) SendGuessSignal(firstID int64, secondID int64, guessResults models.GuessResult) {
	println("The guess was " + guessResults.Guessed)
	println("The word was " + guessResults.State.CurrWord)
	println("Whether the guess was Correct: ")
	println(guessResults.Correct)

	result := guessResults
	state := *result.State

	firstMessageObj := GameSignal{Type: GuessSignal, UserIDs: []int64{firstID}, State: guessResults}
	message, err := json.Marshal(firstMessageObj)
	if err != nil {
		println("error marshalling game signal for end game")
		return
	}
	body := message
	err = ctx.Channel.Publish(
		"",             // exchange
		ctx.Queue.Name, // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		println("err was not nil " + err.Error())
	}

	state.CurrWord = ""
	result.State = &state
	secondMessageObj := GameSignal{Type: GuessSignal, UserIDs: []int64{secondID}, State: result}
	message, err = json.Marshal(secondMessageObj)
	if err != nil {
		println("error marshalling game signal for end game")
		return
	}
	body = message
	err = ctx.Channel.Publish(
		"",             // exchange
		ctx.Queue.Name, // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		println("err was not nil " + err.Error())
	}
}

// SendSkipSignal notifies the users that the word has been skipped
func (ctx *Context) SendSkipSignal(firstID int64, secondID int64, state *models.GameState) {
	println("Word skipped. New word is" + state.CurrWord)
	stateObj := *state
	firstMessageObj := GameSignal{Type: SkipSignal, UserIDs: []int64{firstID}, State: stateObj}
	firstMessage, err := json.Marshal(firstMessageObj)
	if err != nil {
		println("there was an error marshalling first message")
		return
	}
	stateObj.CurrWord = ""
	secondMessageObj := GameSignal{Type: SkipSignal, UserIDs: []int64{secondID}, State: stateObj}
	secondMessage, err := json.Marshal(secondMessageObj)
	if err != nil {
		println("There was an error marshalling second message")
		return
	}
	body := firstMessage
	err = ctx.Channel.Publish(
		"",             // exchange
		ctx.Queue.Name, // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		println("err was not nil " + err.Error())
	}
	body = secondMessage
	err = ctx.Channel.Publish(
		"",             // exchange
		ctx.Queue.Name, // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		println("err was not nil " + err.Error())
	}
}

// UpdateLeaderboard updates the leaderboard
func (ctx *Context) UpdateLeaderboard(results models.Results) {
	println("TODO: implement updateLeaderboard")
	println("updateleaderboard: numplayed: ")
	println(results.NumPlayed)
	// Implement UpdateLeaderboard use the results, update the leaderboard
	insert, err := ctx.Db.Exec(sqlInsertUsers, results.ActorID, results.GuesserID, results.NumRight, results.NumPlayed)
	if err != nil {
		println("err was not nil  " + err.Error())
	}
	println(insert)
}

//EndGame ends the game
func (ctx *Context) EndGame(key int64) {
	println("The game was ended")
	currGame, foundGame := ctx.Games[key]
	if !foundGame { // game not found
		return
	}

	delete(ctx.Games, key) // delete the gamestate from the games list
	results := currGame.EndGame()
	first := results.ActorID
	second := results.GuesserID
	println("results was good")
	resString, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		println("there was an error: " + err.Error())
		return
	}
	println("Result string was " + string(resString))

	ctx.SendEndGameSignal(first, second, results)
	ctx.UpdateLeaderboard(results)
}

// newEndGameFunc returns a function for ending the game
func (ctx *Context) newEndGameFunc(id int64) func() {
	return func() {
		ctx.EndGame(id)
	}
}
