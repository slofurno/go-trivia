package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

var lock = &sync.Mutex{}
var players = make([]*TriviaPlayer, 0)
var game = newTriviaGame()
var highscore = &TriviaPlayer{name: "player", points: 0}

func broadcast(msg string) {
	lock.Lock()
	defer lock.Unlock()

	for _, c := range players {
		c.conn.WriteS(msg)
	}

}

type TriviaQuestion struct {
	question string
	answer   string
}

type PlayerMessage struct {
	player  *TriviaPlayer
	message string
}

type TriviaPlayer struct {
	name   string
	points int
	conn   *WebSocket
}

type TriviaGame struct {
	Inbox           chan *PlayerMessage
	Questions       chan *TriviaQuestion
	CurrentQuestion *TriviaQuestion
	NeedQuestion    chan struct{}
}

func newTriviaGame() *TriviaGame {
	return &TriviaGame{
		Inbox:        make(chan *PlayerMessage, 1028),
		Questions:    make(chan *TriviaQuestion, 256),
		NeedQuestion: make(chan struct{}, 1),
	}
}

func lower(s string) string {
	return strings.ToLower(s)
}

func main() {

	go func() {

		timer := time.NewTimer(time.Second * 2)

		for {

			func() {

				select {

				case <-timer.C:
					fmt.Println("timer elapsed")
					if game.CurrentQuestion == nil {

						select {
						case game.CurrentQuestion = <-game.Questions:
							broadcast("Q: " + game.CurrentQuestion.question + "?")
						default:
							broadcast("there are no questions left to ask...")
							timer = time.NewTimer(time.Second * 10)
						}

					}

				case next := <-game.Inbox:
					msg := next.message

					if len(msg) < 1 {
						return
					}

					if msg[:1] == "/" {
						args := strings.SplitN(msg[1:], " ", 2)

						if len(args) != 2 {
							return
						}

						switch lower(args[0]) {
						case "name":
							next.player.name = args[1]
							next.player.conn.WriteS("you changed your name to " + args[1])

						case "add":
							qs := strings.Split(args[1], "?")

							if len(qs) == 2 {
								broadcast(next.player.name + " added a question")

								question := &TriviaQuestion{question: qs[0], answer: strings.TrimSpace(qs[1])}
								if game.CurrentQuestion == nil {
									game.CurrentQuestion = question
									broadcast("Q: " + game.CurrentQuestion.question + "?")
								} else {
									game.Questions <- question
								}

							}

						}

						return
					} else {

						broadcast(next.player.name + ": " + msg)

						if game.CurrentQuestion != nil && lower(game.CurrentQuestion.answer) == lower(msg) {

							broadcast(next.player.name + " got it right!")
							//broadcast("Q: " + game.CurrentQuestion.question + " A: " + game.CurrentQuestion.answer)
							next.player.points += 1

							if next.player.points > highscore.points {
								highscore.name = next.player.name
								highscore.points = next.player.points
							}

							broadcast(fmt.Sprintln("current highscore: ", highscore.name, "(", highscore.points, " points )"))

							game.CurrentQuestion = nil

							select {
							case game.CurrentQuestion = <-game.Questions:
								broadcast("Q: " + game.CurrentQuestion.question)
							default:
								timer = time.NewTimer(time.Second * 2)
							}

						}

					}

				}
			}()
		}
	}()

	http.HandleFunc("/question/", questionHandler)
	http.ListenAndServe(":80", nil)

}

func questionHandler(w http.ResponseWriter, req *http.Request) {

	ws := upgrade(w, req)
	player := &TriviaPlayer{name: "player", conn: ws}

	ws.WriteS("welcome!")
	ws.WriteS("here are the commands")
	ws.WriteS("/name {name}")
	ws.WriteS("/add {question}? {answer}")
	ws.WriteS("eg. /add is this cool? yes")

	func() {

		lock.Lock()
		defer lock.Unlock()
		players = append(players, player)

	}()

	func() {
		for {
			read, code, err := ws.Read()
			if err != nil || code == Close {
				return
			}

			game.Inbox <- &PlayerMessage{message: read, player: player}

		}
	}()

	func() {

		lock.Lock()
		defer lock.Unlock()

		for i, p := range players {
			if p == player {
				players = append(players[:i], players[i+1:]...)
				return
			}
		}
	}()

}
