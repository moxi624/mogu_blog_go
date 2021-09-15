package admin

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"io/ioutil"
	"mogu-go-v2/common"
	"mogu-go-v2/controllers/base"
	"mogu-go-v2/models/serverInfo"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/12 11:06 上午
 * @version 1.0
 */

type ServerMonitorRestApi struct {
	base.BaseController
}

func (c *ServerMonitorRestApi) GetInfo() {
	var cpuModel serverInfo.Cpu
	count, _ := cpu.Counts(false)
	cpuModel.CpuNum = count
	used, _ := cpu.Percent(0, false)
	times, _ := cpu.Times(false)
	var total float64
	var system float64
	var nice float64
	var idle float64
	var iowait float64
	for _, t := range times {
		total += cpu.TimesStat.Total(t)
		system += t.System
		nice += t.Nice
		idle += t.Idle
		iowait += t.Iowait
	}
	cpuModel.Total = common.Round(total, 2)
	cpuModel.Sys = common.Round(system/total*100, 2)
	cpuModel.Used = common.Round(used[0], 2)
	cpuModel.Wait = common.Round(iowait/total*100, 2)
	cpuModel.Free = common.Round(idle/total*100, 2)

	var memModel serverInfo.Mem
	m, _ := mem.VirtualMemory()
	memModel.Total = common.Round(float64(m.Total/1024/1024/1024), 2)
	memModel.Used = common.Round(float64(m.Used/1024/1024/1024), 2)
	memModel.Free = common.Round(float64(m.Free/1024/1024/1024), 2)
	memModel.Usage = common.Round(m.UsedPercent, 2)

	var sys serverInfo.Sys
	h, _ := host.Info()
	sys.ComputerName = h.Hostname
	sys.ComputerIp = common.IpUtils.GetOutboundIP().String()
	sys.UserDir = ""
	sys.OsName = h.OS
	sys.OsArch = h.KernelArch

	var sysFile serverInfo.SysFile
	disks, _ := ioutil.ReadDir("/")
	path := []string{"/"}
	var sysFiles []serverInfo.SysFile
	for _, d := range disks {
		if d.IsDir() {
			path = append(path, "/"+d.Name())
		}
	}
	for _, p := range path {
		usage, _ := disk.Usage(p)
		patisions, _ := disk.Partitions(true)
		for _, partition := range patisions {
			if partition.Fstype == usage.Fstype {
				sysFile.TypeName = partition.Device
			}
			break
		}
		sysFile.DirName = usage.Path
		sysFile.SysTypeName = usage.Fstype
		sysFile.Total = convertFileSize(usage.Total)
		sysFile.Free = convertFileSize(usage.Free)
		sysFile.Used = convertFileSize(usage.Used)
		sysFile.Usage = common.Round(usage.UsedPercent, 2)
		sysFiles = append(sysFiles, sysFile)
	}

	result := map[string]interface{}{
		"cpu":      cpuModel,
		"mem":      memModel,
		"sys":      sys,
		"sysFiles": sysFiles,
	}
	c.SuccessWithData(result)
}

func convertFileSize(size uint64) string {
	s := float64(size)
	const kb float64 = 1024
	const mb = kb * 1024
	const gb = mb * 1024
	if s >= gb {
		return fmt.Sprintf("%.1f GB", s/gb)
	} else if s >= mb {
		f := s / mb
		if f > 100 {
			return fmt.Sprintf("%.0f MB", f)
		}
		return fmt.Sprintf("%.1f MB", f)
	} else if s >= kb {
		f := s / kb
		if f > 100 {
			return fmt.Sprintf("%.0f KB", f)
		}
		return fmt.Sprintf("%.1f KB", f)
	}
	return fmt.Sprintf("%d B", size)
}
