package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

// Labels
var lblTime, lblIter, lblStatus *widget.Label

var queen fyne.Resource

func valid_grid() bool {
	if n <= 0 {
		return false
	}
	if len(g) != n {
		return false
	}

	unique := make(map[byte]bool)

	for r := 0; r < n; r++ {
		if len(g[r]) != n {
			return false
		}

		for c := 0; c < n; c++ {
			str := g[r][c]

			if str < byte('A') || str > byte('Z') {
				return false
			}
			unique[str] = true
		}
	}

	if len(unique) != n {
		fmt.Println("DEBUG : FALSE")
		fmt.Println(len(unique))
		return false
	}

	return true
}

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
	n = len(g)

	if !valid_grid() {
		lblStatus.SetText("Status: Invalid input! make sure the number of row, columns and colour are the same")

		g = []string{}
		n = 0
		found = false
		cnt = 0
		build_grid()
		refresh_grid()

		return false
	}

	ans = make([]int, n)
	for i := range ans {
		ans[i] = -1
	}

	lblTime = widget.NewLabel("Time: -")
	lblIter = widget.NewLabel("Iterations: -")

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

	lblIter.SetText("Iterations: -")
	lblTime.SetText("Time: -")
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
		if found {
			lblStatus.SetText("Status: found a solution!")
		} else {
			lblStatus.SetText("Status: No solution found :(")
		}
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

func save_to_image(filename string) {
	if n == 0 {
		return
	}

	sz := 130
	var q image.Image

	decoded, err := png.Decode(bytes.NewReader(queen.Content()))
	if err == nil {
		q = decoded
		sz = q.Bounds().Dx()
	}

	width := n * sz
	height := n * sz
	ul := image.Point{0, 0}
	lr := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{ul, lr})

	for r := 0; r < n; r++ {
		for c := 0; c < n; c++ {
			x := c * sz
			y := r * sz

			color := Colors[g[r][c]-'A']
			outline := image.Rect(x, y, x+sz, y+sz)
			draw.Draw(img, outline, image.Black, image.Point{0, 0}, draw.Src)
			bgRect := image.Rect(x+3, y+3, x+sz-3, y+sz-3)
			draw.Draw(img, bgRect, &image.Uniform{color}, image.Point{0, 0}, draw.Over)

			if ans[r] == c {
				offset := image.Point{
					X: (sz - q.Bounds().Dx()) / 2,
					Y: (sz - q.Bounds().Dy()) / 2,
				}
				targetRect := image.Rect(x+offset.X, y+offset.Y, x+offset.X+q.Bounds().Dx(), y+offset.Y+q.Bounds().Dy())
				draw.Draw(img, targetRect, q, q.Bounds().Min, draw.Over)
			}
		}
	}

	filename += ".png"

	outputPath := filepath.Join("../test", filename)

	f, err := os.Create(outputPath)
	if err != nil {
		lblStatus.SetText("Status: failed to save image.")
		return
	}
	defer f.Close()

	err2 := png.Encode(f, img)
	if err2 != nil {
		fmt.Println(err2)
	}
}

func save_to_txt(filename string) {
	if n == 0 {
		return
	}

	filename += ".txt"

	outputPath := filepath.Join("../test", filename)

	f, err := os.Create(outputPath)
	if err != nil {
		lblStatus.SetText("Status: failed to save txt.")
		return
	}
	defer f.Close()

	for r := 0; r < n; r++ {
		var line string
		for c := 0; c < n; c++ {
			if ans[r] == c {
				line += "#"
			} else {
				line += string(g[r][c])
			}
		}

		if _, err := f.WriteString(line + "\n"); err != nil {
			lblStatus.SetText("Status: failed to save txt.")
			return
		}
	}

}

type MappedColor struct {
	R, G, B uint32
	Char    byte
}

func map_image_input(path string, sz int) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return false
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width != height {
		lblStatus.SetText("Status: input image isn't valid")
	}

	step := width / sz

	n = sz
	g = make([]string, n)
	mp := []MappedColor{}
	cur := byte('A')

	for r := 0; r < n; r++ {
		rowStr := ""
		for c := 0; c < n; c++ {
			cx := (c * step) + step/2
			cy := (r * step) + step/2

			r, g, b, _ := img.At(cx, cy).RGBA()
			curColor := MappedColor{R: r, G: g, B: b}

			ch := byte(0)
			for _, k := range mp {
				if curColor.R == k.R && curColor.G == k.G && curColor.B == k.B {
					ch = k.Char
					break
				}
			}

			if ch == 0 {
				ch = cur
				curColor.Char = ch
				mp = append(mp, curColor)
				cur++
			}

			rowStr += string(ch)
		}
		g[r] = rowStr
	}

	for i := range mp {
		fmt.Println(mp[i].R, mp[i].G, mp[i].B)
	}
	ans = make([]int, n)
	for i := range ans {
		ans[i] = -1
	}

	return true
}

func main() {
	fmt.Println("working")
	n = 0
	found = false
	queen = resourceQueenPng

	a := app.New()
	w := a.NewWindow("Tucil 1")
	w.Resize(fyne.NewSize(1200, 600))

	grid = container.New(layout.NewGridLayout(1))
	lblTime = widget.NewLabel("Time: -")
	lblIter = widget.NewLabel("Iterations: -")
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
	btnSaveImg := widget.NewButton("Save answer as image", func() {
		if !found {
			lblStatus.SetText("Status: Unable to save solution, no solution found!")
			return
		}
		save_to_image("sol-img-" + time.Now().Format("20060102150405"))
		lblStatus.SetText("Status: Solution saved as image!")
	})
	btnSaveTxt := widget.NewButton("Save answer as txt", func() {
		if !found {
			lblStatus.SetText("Status: Unable to save solution, no solution found!")
			return
		}
		save_to_txt("sol-txt-" + time.Now().Format("20060102150405"))
		lblStatus.SetText("Status: Solution saved as text!")
	})

	entryN := widget.NewEntry()
	entryN.SetPlaceHolder("Enter the number of row or column")
	btnLoadImg := widget.NewButton("Load Image Input", func() {
		inputSize, err := strconv.Atoi(entryN.Text)
		if err != nil || inputSize <= 0 {
			dialog.ShowError(fmt.Errorf("Please enter a valid number for N first"), w)
			return
		}

		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			defer reader.Close()

			cek := map_image_input(reader.URI().Path(), inputSize)
			if cek {
				lblStatus.SetText("Status: Image Loaded.")
				if !valid_grid() {
					lblStatus.SetText("Status: Image unable to be loaded.")
					n = 0
					found = false
					cnt = 0
				} else {
					build_grid()
				}
			} else {
				fmt.Println("Eror on image process")
			}
		}, w)

		fd.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg"}))
		fd.Show()
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
		widget.NewLabel("Or... upload as a file:"),
		entryN,
		btnLoadImg,
		btn1, btn2, btn3,
		lblSlider,
		slider,
		layout.NewSpacer(),
		btnSaveImg,
		btnSaveTxt,
		lblStatus, lblTime, lblIter,
	)

	split := container.NewHSplit(sidepanel, container.NewPadded(grid))

	w.SetContent(split)
	w.ShowAndRun()
}
