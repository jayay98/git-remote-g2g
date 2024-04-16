package main

import (
	"image/color"
	"io"
	"log"
	"os"
	"path"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/util"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Table Widget")

	var err error
	mfs := memfs.New()
	r, err := git.Clone(memory.NewStorage(), mfs, &git.CloneOptions{
		URL: "file:///Users/blank/Documents/sandbox/g2g",
	})
	if err != nil {
		log.Fatal(err)
	}

	filepaths := make(map[string][]string)
	latestCommits := make(map[string]*object.Commit)

	util.Walk(mfs, ".", func(p string, info os.FileInfo, err error) error {
		if p == "." {
			return nil
		}

		dir := path.Dir(p)
		filepaths[dir] = append(filepaths[dir], p)
		if info.IsDir() {
			filepaths[p] = []string{}
		} else {
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
