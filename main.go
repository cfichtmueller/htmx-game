package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"cfichtmueller.com/htmx-game/internal/client"
	"cfichtmueller.com/htmx-game/internal/engine"
	"cfichtmueller.com/htmx-game/internal/ui"
)

func main() {

	game := engine.New(1000, 600)

	game.Start()

	http.HandleFunc("/js/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		w.Header().Add("Content-Type", "application/javascript")
		ui.RenderScript(w, name)
	})

	http.HandleFunc("/img/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		w.Header().Add("Content-Type", "image/png")
		ui.RenderImage(w, name)
	})

	http.HandleFunc("/css/style.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/css")
		w.Write(ui.Css)
	})

	http.HandleFunc("/player/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		game.Lock()
		defer game.Unlock()
		p, ok := game.PlayerWithId(id)
		if !ok {
			w.WriteHeader(404)
			return
		}

		if r.Method == "GET" {
			w.Header().Set("Cache-Control", "no-store")
			includeShell := r.Header.Get("Hx-Request") != "true"
			if includeShell && !must("render shell start", ui.RenderShellStart(w)) {
				return
			}

			if !must("render index", ui.RenderIndexPage(w, game, id)) {
				return
			}
			if includeShell && !must("render shell end", ui.RenderShellEnd(w)) {
				return
			}
			return
		}

		if r.Method != "POST" {
			w.WriteHeader(405)
			return
		}

		b, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("unable to read body: %v", err)
			w.WriteHeader(500)
			return
		}

		var input PlayerInput
		json.Unmarshal(b, &input)

		if engine.IsEntityDead(game.World, p) {
			return
		}

		for _, cmd := range input.Commands {
			switch cmd.M {
			case "setVelocity":
				engine.SetEntityVelocity(game.World, p, cmd.V)
			case "setRotation":
				engine.SetEntityDirection(game.World, p, cmd.V)
			case "respawn":
				engine.KillEntity(game.World, p)
				return
			}
		}
	})

	http.HandleFunc("/player/{id}/osd", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(405)
			return
		}
		id := r.PathValue("id")
		game.Lock()
		defer game.Unlock()
		p, ok := game.PlayerWithId(id)
		if !ok {
			w.WriteHeader(404)
			return
		}
		must("render osd", ui.RenderOsd(w, p))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			w.WriteHeader(404)
			return
		}
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "GET")
			w.WriteHeader(204)
			return
		}
		game.Lock()
		p := game.SpawnPlayer()
		w.Header().Set("Location", "/player/"+p)
		game.Unlock()
		w.Header().Set("Cache-Control", "no-store")
		w.WriteHeader(302)
	})

	http.HandleFunc("/field", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		width := q.Get("w")
		height := q.Get("h")

		cstate, err := client.NewState(width, height)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		game.Lock()
		cstate.Update(game)
		game.Unlock()

		if !must("render field", ui.RenderField(w, cstate)) {
			return
		}
	})

	log.Fatal(http.ListenAndServe("127.0.0.1:3000", nil))
}

func must(what string, err error) bool {
	if err != nil {
		log.Printf("unable to %s: %v", what, err)
		return false
	}
	return true
}

type PlayerInput struct {
	Commands []PlayerCommand `json:"commands"`
}

type PlayerCommand struct {
	M string  `json:"m"`
	V float64 `json:"v"`
}
