package main

import (
    "log"
    "net/http"
    "github.com/googollee/go-socket.io"
)

type player struct {
    id int;
    socket socketio.Socket;
}

func main() {
    server, err := socketio.NewServer(nil)
    if err != nil {
        log.Fatal(err)
    }
    
    var MAX_NUMBERS_OF_PLAYERS int = 5

    var players [5]player;
    
    var main_client socketio.Socket;
    
    for i := 0; i < MAX_NUMBERS_OF_PLAYERS; i++ {
        players[i].id = -1
    }
    
    server.On("connection", func(so socketio.Socket) {
        log.Println("on connection")
        so.Join("spacejunk")
        
        so.On("spacejunk player", func(msg string) {
            var freeId int = -1
                    
            for i := 0; i < MAX_NUMBERS_OF_PLAYERS; i++ {
                if players[i].id == -1 {
                    freeId = i                
                    break
                }
            }
            
            if freeId == -1 {
                so.Emit("disconnect")
                return
            }
            
            players[freeId].id = freeId;
            players[freeId].socket = so;
            
 
            
            so.Emit("spacejunk accept", freeId)
            
            
            for i := 0; i < MAX_NUMBERS_OF_PLAYERS; i++ {
                if players[i].id != -1 {
                    so.Emit("spacejunk newPlayer", i)              
                }
            }
            
            so.BroadcastTo("spacejunk", "spacejunk newPlayer", freeId)
            
            so.On("spacejunk shootRocket", func(msg string) {
                so.BroadcastTo("spacejunk", "spacejunk shootRocket", msg)
            })
            
            so.On("spacejunk shootPlayer", func(msg string) {
                log.Println(msg)
                so.Emit(msg)
                so.BroadcastTo("spacejunk", "spacejunk shootPlayer", msg)
            })
            
            so.On("disconnection", func() {
                log.Println("on disconnect")
                players[freeId].id = -1
                so.BroadcastTo("spacejunk", "spacejunk playerLeft", freeId)
            })
        })
        
        so.On("spacejunk server", func(msg string) {
            for i := 0; i < MAX_NUMBERS_OF_PLAYERS; i++ {
                if players[i].id != -1 {
                    so.Emit("spacejunk newPlayer", i)              
                }
            }
            
            main_client = so
        })
    })
    
    server.On("error", func(so socketio.Socket, err error) {
        log.Println("error:", err)
    })

    http.Handle("/socket.io/", server)
    http.Handle("/", http.FileServer(http.Dir("./asset")))
    log.Println("Serving at localhost:8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}