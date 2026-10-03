package agent
import("net";"time";gnet "github.com/shirou/gopsutil/v4/net")
type ifaceSample struct{rxBytes,txBytes,rxPackets,txPackets uint64;at time.Time}
func interfaceStats(prev map[string]ifaceSample)([]InterfaceStats,map[string]ifaceSample){
 now:=time.Now();counters,_:=gnet.IOCounters(true);out:=make([]InterfaceStats,0,len(counters));next:=make(map[string]ifaceSample,len(counters))
 for _,c:=range counters{up:=false;if i,e:=net.InterfaceByName(c.Name);e==nil{up=i.Flags&net.FlagUp!=0};s:=ifaceSample{c.BytesRecv,c.BytesSent,c.PacketsRecv,c.PacketsSent,now};x:=InterfaceStats{Name:c.Name,Up:up,RxBytes:c.BytesRecv,TxBytes:c.BytesSent,RxPackets:c.PacketsRecv,TxPackets:c.PacketsSent}
  if old,ok:=prev[c.Name];ok{d:=now.Sub(old.at).Seconds();if d>0{x.RxPPS=float64(c.PacketsRecv-old.rxPackets)/d;x.TxPPS=float64(c.PacketsSent-old.txPackets)/d;x.RxMbps=float64(c.BytesRecv-old.rxBytes)*8/d/1e6;x.TxMbps=float64(c.BytesSent-old.txBytes)*8/d/1e6}}
  if i,e:=net.InterfaceByName(c.Name);e==nil{x.Addresses=interfaceAddresses(i)};out=append(out,x);next[c.Name]=s}
 return out,next}
