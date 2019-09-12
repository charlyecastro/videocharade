# Info 441 Final Project
> Team members: Allen Cho, Brandon Chen, Charlye Castro

Charades441 is a charades app that lets users play charades over the internet with real time video chat. 

## Note
webRTC is supported by mozilla fire fox and google chrome. recently webRTC released new methods for hanlding video chat calls,meaning some functions are becoming deprecated. Video Chat does not work on all browsers, so if it doesnt work on yours either try on firefox or chrome. Maybe try an incognito window. Video chat can also fail due to chrome extensions interferring.

## Why

Our app is intended for responsible adults who want to have fun playing charades on the internet with like minded individuals. Users may act out or guess random words to earn internet points and climb the leaderboards. 

## What 

We’re creating a charades game, but with the real time video chat. Users get to compete with all the competitive users in the world. In user onboarding and periodically throughout the time that the user is interacting with the app, we will show the user statistics of number of wins, losts, and tries. Also, the world ranking will be displayed. 

# Architectural Diagram
![Architectural Diagram](/images/charadesDiagram.png?raw=true "Architectural Diagram")



## User Stories

| Priority       | User           | Description  |
| ------------- |-------------| -----|
| p0 | As an actor | I want to receive a random word to act out |
| p0 | As a guesser | I want to see the actors live video feed |
| p0 |  As a guesser | I want to guess words with live feedback |
| p0 | As a user | I want to see the leaderboards for most correct guesses |

## Technical Implementation Strategy

#### Charades Game

The game will be handled by the Charade microservice. Users may begin a game by sending a request which will be forwarded to the service. The microservice will send game state updates to a containerized **RabbitMQ** message queue, which are then sent to the user through a **websocket** connection. The charade microservice will be a containerized **Golang** server that will serve random words stored on disk, and track leaderboards in a containerized **MySQL** database. 

#### Video Feed

The guesser will be able to see the actor’s live video feed using **WebRTC**. Connecting the actor and guesser will be handled by a websocket connection in the api gateway. The Chat handler will be appended to the our gateway since we already established a **websocket** handler. We will append a feature where the gateway will expect a **websocket** message for **requesting** and **accepting** video chats. Through these chat interaction the server will share user information (SDP, ICE Candidate, etc) over a **websocket** connection, thus initiating webRTC peer to peer communication.

## Appendix

#### Database Schemas
##### Leaderboard Table
```sql
-- Leaderboards
CREATE TABLE leaderboard (
	id INT AUTO_INCREMENT,
	firstUserID INT NOT NULL,
	secondUserID INT NOT NULL,
	numCorrect INT NOT NULL,
	numPlayed INT NOT NULL,
	numSkipped INT NOT NULL,
	numSkipped INT NOT NULL
)
``` 

### Charades Microservice API Documentation

##### POST `/v1/charades`
Starts a new game

Request Body:
```javascript
  {
    "FirstUserID": 1,
    "SecondUserID": 2
  }
```

Response Codes:
- 201: Game created successfully
- 400: Bad request body. Must be JSON, FirstUserID, SecondUserID must be set
- 405: Method other than POST
- 409: Users already in game
- 415: Request header Content-Type must be `application/json`

Response Body:
```javascript
{
    "actorID": number,
    "guesserID": number,
    "timeStart": datetime,
    "numGuessed": number,
    "numPlayed": number,
    "numRight": number,
    "numSkipped": number,
    "currWord": string
}
```

WebSocket Messages:
- Clients in the userlist will receive this websocket message
  - Guesser will see `"currWord": ""` instead of the word

At start of game
```javascript
{
  "type": "game-start",
  "userList": number[],
  "data": {
    "actorID": number,
    "guesserID": number,
    "timeStart": datetime,
    "numGuessed": number,
    "numPlayed": number,
    "numRight": number,
    "numSkipped": number,
    "currWord": string
  }
}
```
At end of game
```javascript
{
  "type":"game-end",
  "userList": number[],
  "data": {
    "actorID": number,
    "guesserID": number,
    "numGuessed": number,
    "numPlayed": number,
    "numRight": number,
    "numSkipped": number
  }
}
```

##### POST `/v1/charades/guess`
Make a guess for the current game

Request Body:
```javascript
  {
    "FirstUserID": 1,
    "SecondUserID": 2,
    "Guess": "word"
  }
```

Response Codes:
- 200: Guess made successfully
- 400: Bad request body. Must be JSON, FirstUserID, SecondUserID, Guess must be set
- 404: Active game not found. Game may be expired.
- 405: Method other than POST

Response Body:
```javascript
{
    "correct": boolean,
    "guessed": string,
    "GameState": {
        "actorID": number,
        "guesserID": number,
        "timeStart": datetime,
        "numGuessed": number,
        "numPlayed": number,
        "numRight": number,
        "numSkipped": number,
        "currWord": string
    }
}
```

WebSocket Messages:
- Clients in the userlist will receive this websocket message
  - Guesser will see `"currWord": ""` instead of the word
```javascript
{
  "type": "guess",
  "userList": number[],
  "data": {
    "correct": boolean,
    "guessed": string,
    "GameState": {
      "actorID": number,
      "guesserID": number,
      "timeStart": datetime,
      "numGuessed": number,
      "numPlayed": number,
      "numRight": number,
      "numSkipped": number,
      "currWord": string
    }
  }
}
```


###### POST `/v1/charades/skip`
Skips the word for the current game

Request Body:
```javascript
  {
    "FirstUserID": 1,
    "SecondUserID": 2
  }
```

Response Codes:
- 200: Skipped word successfully
- 400: Bad request body. Must be JSON, FirstUserID, SecondUserID must be set
- 404: Active game not found. Game may be expired.
- 405: Method other than POST

Response Body:
```javascript
{
    "actorID": number,
    "guesserID": number,
    "timeStart": datetime,
    "numGuessed": number,
    "numPlayed": number,
    "numRight": number,
    "numSkipped": number,
    "currWord": string
}
```

WebSocket Messages:
- Clients in the userlist will receive this websocket message
  - Guesser will see `"currWord": ""` instead of the word
```javascript
{
  "type": "skip",
  "userList": number[],
  "data": {
    "actorID": number,
    "guesserID": number,
    "timeStart": datetime,
    "numGuessed": number,
    "numPlayed": number,
    "numRight": number,
    "numSkipped": number,
    "currWord": string
  }
}
```

###### GET `/v1/leaderboards`
Gets the top 10 games played. Shows username for actor and guesser, including number played, and number right

- 200: Successfully fetched leaderboard data
- 405: Method other than GET
- 500: Server could not find user data for games

Response Body:
```javascript
// JSON array of top 10 games
[
  {
    "actorID": string,    // the username of the actor
    "guesserID": string,  // the username of the guesser
    "numPlayed": number,
    "numRight": number
  }
]
```

###### POST `/v1/handleOffer`
Pushes Video Chat request to message queue which is then returned through websockets, to all users involved

Request Body:
```javascript
  {
    "type": "video-invitation",
    "data": "{}",
    "userList": "[1]"
  }
```

Response Codes:
- 200: Offer Received
