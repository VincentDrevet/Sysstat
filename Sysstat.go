package main

import (
	"log"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	terminal "github.com/wayneashleyberry/terminal-dimensions"
)

func main() {

	x, _ := terminal.Width()
	y, _ := terminal.Height()

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	//TimeBar
	timebar := widgets.NewParagraph()
	timebar.SetRect(0, 0, int(x), 3)
	timebar.TextStyle.Fg = ui.ColorWhite
	timebar.BorderStyle.Fg = ui.ColorCyan

	// RAM gauge
	memorygauge := widgets.NewGauge()
	memorygauge.Title = "Percentage used memory"
	memorygauge.BarColor = ui.ColorGreen
	memorygauge.SetRect(0, 6, int(x)/2, 3)

	//Swap gauge
	swapgauge := widgets.NewGauge()
	swapgauge.Title = "Percentage used swap"
	swapgauge.BarColor = ui.ColorGreen
	swapgauge.SetRect(int(x)/2, 6, int(x), 3)

	// Paragraph InfoCPU
	cpuparagraph := widgets.NewParagraph()
	cpuparagraph.Title = "CPU Info"
	cpuparagraph.SetRect(0, 3, int(x)/2, 8)
	cpuparagraph.TextStyle.Fg = ui.ColorWhite
	cpuparagraph.BorderStyle.Fg = ui.ColorCyan

	// Paragraph InfoMemory
	memparagraph := widgets.NewParagraph()
	memparagraph.Title = "Memory Info"
	memparagraph.SetRect(int(x)/2, 3, int(x), 8)
	memparagraph.TextStyle.Fg = ui.ColorWhite
	memparagraph.BorderStyle.Fg = ui.ColorCyan

	// tabpanel
	Pane := widgets.NewTabPane()
	Pane.Title = "Sysstat"
	Pane.SetRect(0, int(y)-3, int(x), int(y))
	Pane.TabNames = []string{"General Info", "System metrics"}
	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "<Left>":
				Pane.FocusLeft()
				ui.Clear()
			case "<Right>":
				Pane.FocusRight()
				ui.Clear()
			}
		case timetick := <-ticker.C:
			switch Pane.ActiveTabIndex {
			case 0:
				GetCPUInfo(cpuparagraph)
				GetMemInfo(memparagraph)
				UpdateTime(timebar, timetick)
				ui.Render(Pane, timebar)
			case 1:
				UpdateMemoryPercentage(memorygauge, swapgauge)
				UpdateTime(timebar, timetick)
				ui.Render(Pane, timebar)
			}
			//UpdateTime(p, timeticker)
		}
	}

}
