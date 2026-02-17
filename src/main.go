package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// Labels
var lblTime, lblIter, lblStatus *widget.Label

var queen fyne.Resource

func read_file(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	g = []string{}
	for sc.Scan() {
		cnt += 1
		line := sc.Text()
		if len(line) > 0 {
			g = append(g, line)
		}
	}

	// TODO : add validation

	n = len(g)
	ans = make([]int, n)
	for i := range ans {
		ans[i] = -1
	}
	return true
}

func refresh_grid() {
	if n == 0 {
		return
	}
	for r := 0; r < n; r++ {
		for c := 0; c < n; c++ {
			if ans[r] == c {
				texts[r][c].Show()
			} else {
				texts[r][c].Hide()
			}
			texts[r][c].Refresh()
		}
	}
}

func run_solver(mode int, live bool, w *fyne.Window) {
	if n == 0 {
		dialog.ShowError(fmt.Errorf("Load file first"), *w)
		return
	}

	// reset
	found = false
	cnt = 0
	stop = false
	ans = make([]int, n)
	for i := range ans {
		ans[i] = -1
	}

	lblStatus.SetText("Status: Running...")
	start := time.Now()
	go func() {
		if mode == 1 {
			solve1(0, live)
		} else {
			mask := make(map[byte]bool)
			solve2(0, live, mask)
		}

		dur := time.Since(start)

		refresh_grid()
		lblStatus.SetText("Status: found a solution!")
		lblTime.SetText(fmt.Sprintf("Time: %d ms", dur.Milliseconds()))
		lblIter.SetText(fmt.Sprintf("Iterations: %d", cnt))
	}()

}

func build_grid() {
	rects = make([][]*canvas.Rectangle, n)
	texts = make([][]*canvas.Image, n)
	objects := []fyne.CanvasObject{}

	for r := 0; r < n; r++ {
		rects[r] = make([]*canvas.Rectangle, n)
		texts[r] = make([]*canvas.Image, n)
		for c := 0; c < n; c++ {
			rect := canvas.NewRectangle(Colors[g[r][c]-'A'])
			rect.StrokeColor = Colors[g[r][c]-'A']

			img := canvas.NewImageFromResource(queen)
			img.FillMode = canvas.ImageFillContain
			img.Hide()

			rects[r][c] = rect
			texts[r][c] = img

			stack := container.NewStack(rect, img)
			objects = append(objects, stack)
		}
	}
	grid.Layout = layout.NewGridLayout(n)
	grid.Objects = objects
	grid.Refresh()
}

func main() {
	fmt.Println("working")
	n = 0
	queen, _ = fyne.LoadResourceFromPath("../assets/queen.png")

	a := app.New()
	w := a.NewWindow("Tucil 1")
	w.Resize(fyne.NewSize(1200, 600))

	grid = container.New(layout.NewGridLayout(1))
	lblTime = widget.NewLabel("Time: 0ms")
	lblIter = widget.NewLabel("Iterations: 0")
	lblStatus = widget.NewLabel("Status: Waiting for input")

	btn1 := widget.NewButton("Solver 1 (Live without pruning)", func() { run_solver(1, true, &w) })
	btn2 := widget.NewButton("Solver 2 (Live with Pruning)", func() { run_solver(2, true, &w) })
	btn3 := widget.NewButton("Solver 3 (No live updates)", func() { run_solver(1, false, &w) })
	btnLoad := widget.NewButton("Load Input", func() {
		dialog.ShowFileOpen(func(r fyne.URIReadCloser, err error) {
			if r == nil {
				return
			}

			if read_file(r.URI().Path()) {
				build_grid()
				lblStatus.SetText("Status: File Loaded")
			}
		}, w)
	})

	lblSlider := widget.NewLabel("Update speed slider (ms): ")
	slider := widget.NewSlider(1, 500)
	slider.SetValue(50)
	slider.OnChanged = func(v float64) { delay = int(v) }

	uiChan = make(chan bool)
	go func() {
		for range uiChan {
			refresh_grid()
		}
	}()

	sidepanel := container.NewVBox(
		btnLoad,
		btn1, btn2, btn3,
		lblSlider,
		slider,
		layout.NewSpacer(),
		lblStatus, lblTime, lblIter,
	)

	split := container.NewHSplit(sidepanel, container.NewPadded(grid))

	w.SetContent(split)
	w.ShowAndRun()
}
