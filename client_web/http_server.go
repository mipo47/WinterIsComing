package main

import (
	"../client/gamer"
	"../core"
	"net/http"
	"log"
	"strconv"
	"encoding/json"
	"bytes"
	"fmt"
	"net"
	"math/rand"
	"time"
)

type HttpServer struct {
	sessions map[int]*HttpSession
}

func (s *HttpServer) Start(port int)  {
	s.sessions = make(map[int]*HttpSession)

	fmt.Println("Open http://localhost:"+ strconv.Itoa(port))

	http.Handle("/", http.FileServer(http.Dir("./client_web/html")))
	http.HandleFunc("/start", s.start)
	http.HandleFunc("/status", s.status)
	http.HandleFunc("/shoot", s.shoot)

	if err := http.ListenAndServe(":" + strconv.Itoa(port), nil); err != nil {
		log.Fatalf("Failed listen port 8080: %v", err)
	}
}

func sendJson(w http.ResponseWriter, result interface{})  {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(result)
	w.Write(b.Bytes())
}

func (s *HttpServer) getSession(w http.ResponseWriter, r *http.Request) *HttpSession {
	q := r.URL.Query()
	sessionString := q.Get("SessionID")
	sessionID, err := strconv.Atoi(sessionString)
	if err != nil {
		panic("SessionID in wrong format: " + sessionString)
	}
	session, found := s.sessions[sessionID]
	if !found {
		panic("Session not found: " + sessionString)
	}
	return session
}

func (s *HttpServer) start(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	useAI := q.Get("AI") == "true"

	if core.LOG_INFO {
		log.Println("Connecting to localhost:", core.TCP_PORT)
	}

	conn, err := net.Dial("tcp", "localhost:" + strconv.Itoa(core.TCP_PORT))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	ioServer := core.StartCommandIO(conn, "CLIENT")

	var ai gamer.AI
	if useAI {
		ai = new(gamer.AI_Closest)
	}
	session := CreateHttpSession(ai, *ioServer)
	sessionID := rand.Intn(1000000)
	s.sessions[sessionID] = session
	if core.LOG_INFO {
		log.Println("New session:", sessionID)
	}

	go session.TrackCommands()
	go func() {
		session.gamer.Play(*ioServer)
		ioServer.Close()

		time.Sleep(10 * time.Second)
		delete(s.sessions, sessionID)
	}()

	sendJson(w, dtoStart{ SessionID: sessionID })
}

func (s *HttpServer) status(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			if core.LOG_ERROR {
				log.Println("Can't update status:", r)
			}
			sendJson(w, dtoStart { Error: fmt.Sprintf("%v", r) })
		}
	}()

	session := s.getSession(w, r)
	sendJson(w, dtoStatus{
		Width:  core.BOARD_WIDTH,
		Height: core.BOARD_HEIGHT,
		Commands: session.newCommands,
		Game: *session.gamer,
	})
	session.newCommands = session.newCommands[:0]
}

func (s *HttpServer) shoot(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			if core.LOG_ERROR {
				log.Println("Can't shoot to zombie:", r)
			}
			sendJson(w, dtoStart { Error: fmt.Sprintf("%v", r) })
		}
	}()

	session := s.getSession(w, r)

	q := r.URL.Query()
	xString := q.Get("X")
	x, err := strconv.Atoi(xString)
	if err != nil {
		panic("Shoot X in wrong format: " + xString)
	}
	yString := q.Get("Y")
	y, err := strconv.Atoi(yString)
	if err != nil {
		panic("Shoot Y in wrong format: " + yString)
	}

	if core.LOG_INFO {
		log.Println("Shoot to zombie:", x, y)
	}
	session.io.SendCommand("SHOOT", x, y)
}