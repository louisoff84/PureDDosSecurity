package agent
import("runtime";"github.com/shirou/gopsutil/v4/cpu";"github.com/shirou/gopsutil/v4/disk";"github.com/shirou/gopsutil/v4/host";"github.com/shirou/gopsutil/v4/load";"github.com/shirou/gopsutil/v4/mem";"github.com/shirou/gopsutil/v4/process")
func systemStats()SystemStats{
 c,_:=cpu.Percent(0,false);v,_:=mem.VirtualMemory();l,_:=load.Avg();u,_:=host.Uptime();d,_:=disk.Usage("/");ps,_:=process.Pids()
 cp:=0.0;if len(c)>0{cp=c[0]}
 return SystemStats{CPUPercent:cp,MemoryPercent:v.UsedPercent,MemoryUsed:v.Used,MemoryTotal:v.Total,Load1:l.Load1,Load5:l.Load5,Load15:l.Load15,Goroutines:runtime.NumGoroutine(),UptimeSeconds:u,DiskUsedPercent:d.UsedPercent,Processes:uint64(len(ps))}
}
