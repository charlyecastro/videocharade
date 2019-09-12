package handlers

import (
	"database/sql"
	"encoding/json"
	"final-project-crew/videoCharade/servers/charades/models"
	"net/http"
	"time"
)

// Sql statements for retrieving users information to update the leaderboard
const sqlSelectTop10 = "select top 10 * from leaderboards order by numGuessRight desc"

// ExpireTime is the length of a match
const ExpireTime = 60 * time.Second

// LBResult is a stuct representing what will be displayed on the leaderboard
type LBResult struct {
	ID        int64 `json:"id"`
	ActorID   int64 `json:"actorID"`
	GuesserID int64 `json:"guesserID"`
	NumPlayed int   `json:"numPlayed"`
	NumRight  int   `json:"numRight"`
}

// LBReturn is a stuct representing what will be displayed on the leaderboard
type LBReturn struct {
	ActorUserName   string `json:"actorID"`
	GuesserUserName string `json:"guesserID"`
	NumPlayed       int    `json:"numPlayed"`
	NumRight        int    `json:"numRight"`
}

// GamesHandler handles requests for a game
func (ctx *Context) GamesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	// Start a new game
	case "POST":
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Must be JSON", http.StatusUnsupportedMediaType)
			return
		}
		ids := &models.GameRequest{FirstUserID: -1, SecondUserID: -1, Guess: ""}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(ids); err != nil {
			http.Error(w, "Must be JSON", http.StatusBadRequest)
			return
		}
		if ids.FirstUserID == -1 || ids.SecondUserID == -1 {
			http.Error(w, "Both UserIDs must not be empty", http.StatusBadRequest)
			return
		}

		if _, foundGame := ctx.Games[ids.FirstUserID]; foundGame {
			http.Error(w, "User 1 already in game", http.StatusConflict)
			return
		}
		if _, foundGame := ctx.Games[ids.SecondUserID]; foundGame {
			http.Error(w, "User 2 already in game", http.StatusConflict)
			return
		}

		// A game can be started for the users
		ctx.Games[ids.FirstUserID] = models.NewGameState(ids.FirstUserID, ids.SecondUserID)

		// Send end game signal after time
		f := ctx.newEndGameFunc(ids.FirstUserID)
		// Invokes go routine after time to clean up gamestate, notify users
		time.AfterFunc(ExpireTime, f)

		w.Header().Set("Content-Type", "application/json")
		data, err := json.Marshal(ctx.Games[ids.FirstUserID])
		if err == nil {
			w.WriteHeader(http.StatusCreated)
			w.Write(data)
		} else {
			w.Write([]byte("[]"))
		}
		ctx.SendStartGameSignal(ids.FirstUserID, ids.SecondUserID, ctx.Games[ids.FirstUserID])
	default:
		http.Error(w, "Method must be POST", http.StatusMethodNotAllowed)
		return
	}
}

// GuessHandler handles requests for guessing
func (ctx *Context) GuessHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		ids := &models.GameRequest{FirstUserID: -1, SecondUserID: -1, Guess: ""}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(ids); err != nil {
			http.Error(w, "Must be JSON", http.StatusBadRequest)
			return
		}
		if ids.FirstUserID == -1 || ids.SecondUserID == -1 {
			http.Error(w, "Both UserIDs must not be empty", http.StatusBadRequest)
			return
		}

		if ids.Guess == "" {
			http.Error(w, "Guess must not be empty", http.StatusBadRequest)
			return
		}

		// Use IDs to find current Game
		currGame, found := ctx.Games[ids.FirstUserID]
		if !found {
			http.Error(w, "Game not found", http.StatusNotFound)
			return
		}

		correct := currGame.Guess(ids.Guess)
		guessResult := models.GuessResult{Correct: correct, Guessed: ids.Guess, State: currGame}

		// Current game is found, marshal it
		data, err := json.Marshal(guessResult)

		// Write response
		w.Header().Set("Content-Type", "application/json")
		if err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write(data)
		} else {
			w.Write([]byte("[]"))
		}

		// Notifies both users that a guess has been made.
		ctx.SendGuessSignal(ids.FirstUserID, ids.SecondUserID, guessResult)
	default:
		http.Error(w, "Method must be POST", http.StatusMethodNotAllowed)
		return
	}
}

// SkipHandler handles requests for skips
func (ctx *Context) SkipHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// Get new word, update gamestate
	case "POST":
		ids := &models.GameRequest{FirstUserID: -1, SecondUserID: -1, Guess: ""}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(ids); err != nil {
			http.Error(w, "Must be JSON", http.StatusBadRequest)
			return
		}
		if ids.FirstUserID == -1 || ids.SecondUserID == -1 {
			http.Error(w, "Both UserIDs must not be empty", http.StatusBadRequest)
			return
		}
		currGame, found := ctx.Games[ids.FirstUserID]
		if !found {
			http.Error(w, "Game not found", http.StatusNotFound)
			return
		}

		currGame.Skip()

		data, err := json.Marshal(currGame)
		// Write response
		w.Header().Set("Content-Type", "application/json")
		if err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write(data)
		} else {
			w.Write([]byte("[]"))
		}

		// Notifies both users that a word has been skipped
		ctx.SendSkipSignal(ids.FirstUserID, ids.SecondUserID, currGame)
	default:
		http.Error(w, "Method must be POST", http.StatusMethodNotAllowed)
		return
	}
}

// LeaderboardHandler handles requests for the leaderboard
func (ctx *Context) LeaderboardHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		//query top 10 user pairs by comparing the number of guess they got right
		// rows, err := ctx.Db.Query(sqlSelectTop10)
		rows, err := ctx.Db.Query("SELECT * FROM leaderboards ORDER BY numGuessRight DESC Limit 10")

		if err != nil {
			http.Error(w, "Could not get leaderboard results", http.StatusInternalServerError)
			return
		}

		// Array to contain top 10 result rows
		var list []*LBReturn

		// Iterate and append each row into the array
		for rows.Next() {
			log := &LBResult{}
			rows.Scan(&log.ID, &log.ActorID, &log.GuesserID, &log.NumRight, &log.NumPlayed)
			var actorUser string
			var guesserUser string
			err = ctx.Db.QueryRow("SELECT user_name FROM users WHERE id = ?", int(log.ActorID)).Scan(&actorUser)
			err = ctx.Db.QueryRow("SELECT user_name FROM users WHERE id = ?", int(log.GuesserID)).Scan(&guesserUser)
			if err != nil && err != sql.ErrNoRows {
				http.Error(w, "Could not get user details for leaderboard", http.StatusInternalServerError)
				return // proper error handling instead of panic in your app
			}
			lb := &LBReturn{
				ActorUserName:   actorUser,
				GuesserUserName: guesserUser,
				NumPlayed:       log.NumPlayed,
				NumRight:        log.NumRight,
			}
			list = append(list, lb)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, "Could not get user details for leaderboard", http.StatusInternalServerError)
			return
		}

		data, err := json.Marshal(list)
		// Write response
		w.Header().Set("Content-Type", "application/json")
		if err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write(data)
		} else {
			w.Write([]byte("[]"))
		}
	default:
		http.Error(w, "Method must be GET", http.StatusMethodNotAllowed)
		return
	}
}
