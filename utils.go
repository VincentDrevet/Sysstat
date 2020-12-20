package main

import (
	"fmt"
	"log"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/v3/mem"

	ui "github.com/gizak/termui/v3"

	"github.com/gizak/termui/v3/widgets"
)

func UpdateMemoryPercentage(memgauge *widgets.Gauge, swapgauge *widgets.Gauge) {

	//Checking memory
	memory, err := mem.VirtualMemory()
	if err != nil {
		log.Fatalln(err)
	}
	switch {
	case memory.UsedPercent <= 40:
		memgauge.BarColor = ui.ColorGreen
	case (memory.UsedPercent > 40) && (memory.UsedPercent <= 70):
		memgauge.BarColor = ui.ColorYellow
	default:
		memgauge.BarColor = ui.ColorRed
	}
	memgauge.Percent = int(memory.UsedPercent)

	//Checking Swap

	swap, err := mem.SwapMemory()
	if err != nil {
		log.Fatalln(err)
	}
	switch {
	case swap.UsedPercent <= 40:
		swapgauge.BarColor = ui.ColorGreen
	case (swap.UsedPercent > 40) && (swap.UsedPercent <= 70):
		swapgauge.BarColor = ui.ColorYellow
	default:
		swapgauge.BarColor = ui.ColorRed
	}
	ui.Render(memgauge, swapgauge)
}

func UpdateTime(paragraph *widgets.Paragraph, tickertime time.Time) {
	paragraph.Text = tickertime.Format(time.UnixDate)
	ui.Render(paragraph)
}

func GetCPUInfo(paragraphcpuinfo *widgets.Paragraph) {
	cpuinfo, err := cpu.Info()
	if err != nil {
		log.Fatalln(err)
	}
	paragraphcpuinfo.Text = fmt.Sprintf("Vendor : %s\nFamily : %s\nModel : %s\nFrequency : %f\nCache Size : %d\nMicrocode : %s", cpuinfo[0].VendorID, cpuinfo[0].Model, cpuinfo[0].ModelName, cpuinfo[0].Mhz, cpuinfo[0].CacheSize, cpuinfo[0].Microcode)
	ui.Render(paragraphcpuinfo)
}

func GetMemInfo(paragraphmem *widgets.Paragraph) {
	memory, err := mem.VirtualMemory()
	if err != nil {
		log.Fatalln(err)
	}

	swap, err := mem.SwapMemory()
	if err != nil {
		log.Fatalln(err)
	}
	paragraphmem.Text = fmt.Sprintf("Total memory : %d Bytes\nMemory Used : %d Bytes\nPercentage Used : %f %%\n\nTotal Swap : %d Bytes\nSwap Used : %d Bytes", memory.Total, memory.Used, memory.UsedPercent, swap.Total, swap.Used)
	ui.Render(paragraphmem)
}
