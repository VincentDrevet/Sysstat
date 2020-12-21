package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
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

func GetHostInfo(paragraphhost *widgets.Paragraph) {
	host, err := host.Info()
	if err != nil {
		log.Fatalln(err)
	}

	paragraphhost.Text = fmt.Sprintf("Hostname : %s\nUptime : %s\nBootTime : %s\nNumber of processes : %d\nOS : %s\nPlatform : %s\nPlatform Version : %s\nKernel Version : %s\nKernel Arch : %s\nHost ID : %s", host.Hostname, ConvertUptime(host.Uptime), time.Unix(int64(host.BootTime), 0).Format(time.UnixDate), host.Procs, host.OS, host.Platform, host.PlatformVersion, host.KernelVersion, host.KernelArch, host.HostID)
	ui.Render(paragraphhost)
}

func ConvertUptime(uptime uint64) string {
	// Convert Second to minute
	mins := uptime / 60
	rseconds := uptime % 60

	// Convert minutes to hours
	hours := mins / 60
	rmins := mins % 60

	// Convert Hours To days
	days := hours / 24
	rhours := hours % 60

	return fmt.Sprintf("%d days %d:%d:%d", days, rhours, rmins, rseconds)

}

func GetUserStat(paragraphuserstat *widgets.Paragraph) {

	usersstats, err := host.Users()
	if err != nil {
		log.Fatalln(err)
	}

	var paramstring string = ""
	for index, stat := range usersstats {
		if index == 0 {
			paramstring = fmt.Sprintf("User : %s\nTerminal : %s\nHost : %s\nStarted : %s\n---\n", stat.User, stat.Terminal, stat.Host, time.Unix(int64(stat.Started), 0).Format(time.UnixDate))
		} else {
			paramstring = fmt.Sprintf("%sUser : %s\nTerminal : %s\nHost : %s\nStarted : %s\n---\n", paramstring, stat.User, stat.Terminal, stat.Host, time.Unix(int64(stat.Started), 0).Format(time.UnixDate))
		}
	}
	paragraphuserstat.Text = paramstring
	ui.Render(paragraphuserstat)
}

func GetNICStat(tablenicstat *widgets.Table) {
	nics, err := net.Interfaces()
	if err != nil {
		log.Fatalln(err)
	}
	var tablecontent [][]string
	for index, nic := range nics {
		addrs, err := nic.Addrs()
		var alladdress string
		for indexaddress, addr := range addrs {
			if indexaddress == 0 {
				alladdress = fmt.Sprintf("%s", addr)
			} else {
				alladdress = fmt.Sprintf("%s\n%s", alladdress, addr)
			}
		}
		if err != nil {
			log.Fatalln(err)
		}
		if index == 0 {
			tablecontent = [][]string{[]string{"Index", "MTU", "Name", "HardwareAddr", "Flags", "Addrs"},
				[]string{strconv.Itoa(nic.Index), strconv.Itoa(nic.MTU), nic.Name, nic.HardwareAddr.String(), nic.Flags.String(), alladdress}}
		} else {
			tablecontent = append(tablecontent, []string{strconv.Itoa(nic.Index), strconv.Itoa(nic.MTU), nic.Name, nic.HardwareAddr.String(), nic.Flags.String(), alladdress})
		}
	}

	tablenicstat.Rows = tablecontent
	ui.Render(tablenicstat)

}
