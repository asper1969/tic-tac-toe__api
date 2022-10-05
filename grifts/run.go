package grifts

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"tic-tac-toe__api/models"
	"time"

	. "github.com/gobuffalo/grift/grift"
	socketio "github.com/googollee/go-socket.io"
)

const (
	port = ":5000"

	gracefulDelay = 3 * time.Second
)

var _ = Namespace("socket_server", func() {

	Desc("run", "Run socket server")
	Add("run", func(c *Context) error {
		server := socketio.NewServer(nil)

		server.OnConnect("/", func(s socketio.Conn) error {
			url := s.URL()
			gameCode := url.Query().Get("game_code")
			teamID := url.Query().Get("team_id")
			s.SetContext("")
			roomName := fmt.Sprintf("%s:%s", gameCode, teamID)
			s.Join(roomName)

			return nil
		})

		server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
			fmt.Println("notice:", msg)
			fmt.Println("Socket#id:", s.ID())
			fmt.Println("Rooms: ", s.Rooms())

			/**
			* Test broadcast to specific room
			**/
			roomName := "eDF2n:1"
			server.BroadcastToRoom("/", roomName, "room", "Room: eDF2n:1")

			/**
			* Tes db query
			**/
			var result models.Session
			err := models.DB.RawQuery("SELECT * FROM sessions WHERE game_pass = ?", "eDF2n").First(&result)

			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(result)

			s.Emit("reply", "have "+msg)
		})

		server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
			s.SetContext(msg)
			return "recv " + msg
		})

		server.OnEvent("/", "bye", func(s socketio.Conn) string {
			last := s.Context().(string)
			s.Emit("bye", last)
			s.Close()
			return last
		})

		server.OnError("/", func(s socketio.Conn, e error) {
			fmt.Println("meet error:", e)
		})

		server.OnDisconnect("/", func(s socketio.Conn, reason string) {
			fmt.Println("closed", reason)
		})

		go server.Serve()
		defer server.Close()

		http.Handle("/socket.io/", server)
		// http.Handle("/", http.FileServer(http.Dir("./asset")))

		log.Printf("Serving at 127.0.0.1:%s", port)
		log.Fatal(http.ListenAndServe(port, nil))

		//graceful-shutdown
		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			if err := server.Serve(); err != nil {
				log.Fatalf("socketio listen error: %s\n", err)
			}
		}()

		go func() {
			if err := http.ListenAndServe(port, nil); err != nil {
				log.Fatalf("http listen error: %s\n", err)
			}
		}()

		log.Printf("server started by %v", port)

		<-done

		//shutdown delay
		log.Printf("graceful delay: %v\n", gracefulDelay)

		time.Sleep(gracefulDelay)

		log.Println("server stopped")

		if err := server.Close(); err != nil {
			log.Fatalf("server shutdown failed: %s\n", err)
		}

		log.Println("server is shutdown")

		return nil
	})

})
