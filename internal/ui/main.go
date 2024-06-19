package ui

import (
	"embed"
	_ "embed"
	"fmt"
	"html/template"
	"io"

	"cfichtmueller.com/htmx-game/internal/client"
	"cfichtmueller.com/htmx-game/internal/state"
)

var (
	//go:embed js/*
	js embed.FS
	//go:embed style.css
	Css []byte
	//go:embed html/*
	htmlFiles embed.FS
	templates = template.Must(template.New("").ParseFS(htmlFiles, "html/*.html"))
)

func RenderScript(w io.Writer, name string) error {
	b, err := js.ReadFile("js/" + name)
	if err != nil {
		return fmt.Errorf("unable to read script %s: %v", name, err)
	}
	if _, err = w.Write(b); err != nil {
		return fmt.Errorf("unable to write script %s: %v", name, err)
	}
	return nil
}

func RenderShellStart(w io.Writer) error {
	return renderTemplate(w, "ShellStart", nil)
}

func RenderShellEnd(w io.Writer) error {
	return renderTemplate(w, "ShellEnd", nil)
}

type indexPageModel struct {
	State  *state.State
	Player *state.Player
}

func RenderIndexPage(w io.Writer, s *state.State, p *state.Player) error {
	return renderTemplate(w, "IndexPage", indexPageModel{
		State:  s,
		Player: p,
	})
}

func RenderField(w io.Writer, s *client.State) error {
	return renderTemplate(w, "Field", s)
}

type osdModel struct {
	Player *state.Player
}

func RenderOsd(w io.Writer, p *state.Player) error {
	return renderTemplate(w, "Osd", osdModel{Player: p})
}

func renderTemplate(w io.Writer, name string, data any) error {
	if err := templates.ExecuteTemplate(w, name, data); err != nil {
		return fmt.Errorf("unable to render template %s: %v", name, err)
	}
	return nil
}
