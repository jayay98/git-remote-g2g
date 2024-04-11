package main

import (
	"image/color"
	"io"
	"io/fs"
	"log"
	"path"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Table Widget")
	repoPath := "."

	var err error
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		log.Fatal(err)
	}

	filepaths := make(map[string][]string)
	latestCommits := make(map[string]*object.Commit)
	filepath.WalkDir(repoPath, func(p string, d fs.DirEntry, err error) error {
		if p == "." || p == ".git" || strings.HasPrefix(p, ".git/") {
			return nil
		}

		dir := path.Dir(p)
		_, ok := filepaths[dir]
		if !ok {
			filepaths[dir] = []string{}
		}
		filepaths[dir] = append(filepaths[dir], p)

		if !d.IsDir() {
			commits, _ := r.Log(&git.LogOptions{FileName: &p})
			commit, err := commits.Next()
			if err != nil {
				return nil
			}
			latestCommits[p] = commit
		}

		return nil
	})

	topbar := canvas.NewText("top bar", color.Black)
	bindtext := binding.NewString()
	multiline := widget.NewMultiLineEntry()
	multiline.Wrapping = fyne.TextWrapBreak
	multiline.Bind(bindtext)
	treewidget := widget.NewTree(
		func(id widget.TreeNodeID) []widget.TreeNodeID {
			pid := id
			if pid == "" {
				pid = "."
			}
			child, ok := filepaths[pid]
			if ok {
				return child
			}
			return []string{}
		},
		func(id widget.TreeNodeID) bool {
			_, ok := filepaths[id]
			return ok || id == ""
		},
		func(branch bool) fyne.CanvasObject {
			if branch {
				return widget.NewLabel("Branch template")
			}
			return widget.NewButton("Leaf template", func() {})
		},
		func(id widget.TreeNodeID, branch bool, o fyne.CanvasObject) {
			if !branch {
				text := path.Base(id)
				button := o.(*widget.Button)
				button.SetText(text)
				button.OnTapped = func() {
					bindtext.Set(getBlobContent(id, latestCommits[id]))
				}
				button.SetIcon(theme.FileIcon())
			} else {
				text := path.Base(id)
				o.(*widget.Label).SetText(text)
			}
		},
	)

	content := container.NewBorder(topbar, nil, treewidget, nil, multiline)
	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}

func getBlobContent(path string, commit *object.Commit) string {
	file, _ := commit.File(path)
	reader, _ := file.Reader()
	blob, _ := io.ReadAll(reader)
	return string(blob)
}
