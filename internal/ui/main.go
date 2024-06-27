package ui

import (
	"embed"
	_ "embed"
	"fmt"
	"html/template"
	"io"

	"cfichtmueller.com/htmx-game/internal/client"
	"cfichtmueller.com/htmx-game/internal/engine"
)

var (
	//go:embed js/*
	js embed.FS
	//go:embed style.css
	Css []byte
	//go:embed html/*
	htmlFiles embed.FS
	//go:embed img/*
	imgFiles  embed.FS
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

func RenderImage(w io.Writer, name string) error {
	b, err := imgFiles.ReadFile("img/" + name)
	if err != nil {
		return fmt.Errorf("unable to read image %s: %v", name, err)
	}
	if _, err = w.Write(b); err != nil {
		return fmt.Errorf("unable to write image %s: %v", name, err)
	}
	return nil
}

func RenderShellStart(w io.Writer) error {
	return renderTemplate(w, "ShellStart", nil)
}

func RenderShellEnd(w io.Writer) error {
	return renderTemplate(w, "ShellEnd", nil)
}

type playerModel struct {
	ID string
}

type indexPageModel struct {
	Engine *engine.Engine
	Player playerModel
}

func RenderIndexPage(w io.Writer, e *engine.Engine, p string) error {
	return renderTemplate(w, "IndexPage", indexPageModel{
		Engine: e,
		Player: playerModel{ID: p},
	})
}

func RenderField(w io.Writer, s *client.State) error {
	return renderTemplate(w, "Field", s)
}

type osdModel struct {
	Player engine.Entity
}

func RenderOsd(w io.Writer, p engine.Entity) error {
	return renderTemplate(w, "Osd", osdModel{Player: p})
}

func renderTemplate(w io.Writer, name string, data any) error {
	if err := templates.ExecuteTemplate(w, name, data); err != nil {
		return fmt.Errorf("unable to render template %s: %v", name, err)
	}
	return nil
}
