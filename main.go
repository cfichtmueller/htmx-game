package main

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"

	"cfichtmueller.com/htmx-game/internal/client"
	"cfichtmueller.com/htmx-game/internal/state"
	"cfichtmueller.com/htmx-game/internal/ui"
)

func main() {

	s := state.New(1000, 600)

	loopTicker := time.NewTicker(30 * time.Millisecond)
	spawnTicker := time.NewTicker(100 * time.Millisecond)

	go func() {
		last := time.Now().UnixMilli()
		for {
			<-loopTicker.C
			now := time.Now().UnixMilli()
			delta := float64(now-last) / 1000
			last = now
			s.Update(delta)
		}
	}()

	go func() {
		for {
			<-spawnTicker.C
			if len(s.Cells) > 100 {
				continue
			}
			x := rand.Float64() * 10
			if x > 7 {
				s.AddCell(state.NewVelocityPowerUpCell(
					s.Width*rand.Float64(),
					s.Height*rand.Float64(),
				))
			} else {
				s.AddCell(state.NewBulletCell())
			}
		}
	}()

	http.HandleFunc("/js/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		w.Header().Add("Content-Type", "application/javascript")
		ui.RenderScript(w, name)
	})

	http.HandleFunc("/css/style.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/css")
		w.Write(ui.Css)
	})

	http.HandleFunc("/player/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		p := s.PlayerWithId(id)
		if p == nil {
			w.WriteHeader(404)
			return
		}

		if r.Method == "GET" {
			w.Header().Set("Cache-Control", "no-store")
			includeShell := r.Header.Get("Hx-Request") != "true"
			if includeShell && !must("render shell start", ui.RenderShellStart(w)) {
				return
			}

			if !must("render index", ui.RenderIndexPage(w, s, p)) {
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

		for _, cmd := range input.Commands {
			switch cmd.M {
			case "setVelocity":
				p.Agent.SetVelocity(p.Agent.MaxVelocity * cmd.V)
			case "setRotation":
				p.Agent.Direction = cmd.V
			case "respawn":
				p.Die()
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
		p := s.PlayerWithId(id)
		if p == nil {
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
		p := s.SpawnPlayer()
		w.Header().Set("Location", "/player/"+p.ID)
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

		cstate.Update(s)

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
