package agent
import("fmt";"sync";"time")
type detector struct{mu sync.Mutex;window time.Duration;pps,bps,synRatio []float64;unique []int}
func newDetector(w time.Duration)*detector{return &detector{window:w}}
func(d *detector)evaluate(now time.Time,n PacketStats,cfg Config,interfaces []string,host string)(AttackEvent,bool){
 d.mu.Lock();defer d.mu.Unlock();d.pps=append(d.pps,n.PacketsPerSecond);d.bps=append(d.bps,n.BitsPerSecond);ratio:=0.0;if n.TcpPackets>0{ratio=n.SynPerSecond/n.PacketsPerSecond};d.synRatio=append(d.synRatio,ratio);d.unique=append(d.unique,n.UniqueSources)
 limit:=int(d.window.Seconds());if limit<1{limit=1};for len(d.pps)>limit{d.pps=d.pps[1:]};for len(d.bps)>limit{d.bps=d.bps[1:]};for len(d.synRatio)>limit{d.synRatio=d.synRatio[1:]};for len(d.unique)>limit{d.unique=d.unique[1:]}
 ap,ab,as:=avg(d.pps),avg(d.bps),avg(d.synRatio);mu:=maxInt(d.unique);reasons:=[]string{};if ap>=float64(cfg.AlertPPS){reasons=append(reasons,fmt.Sprintf("pps %.0f >= %d",ap,cfg.AlertPPS))};if ab>=float64(cfg.AlertBPS){reasons=append(reasons,fmt.Sprintf("bps %.0f >= %d",ab,cfg.AlertBPS))};if as>=cfg.AlertSYNRatio&&ap>=float64(cfg.AlertPPS)/5{reasons=append(reasons,fmt.Sprintf("syn-ratio %.1f%%",as*100))};if mu>=cfg.AlertUniqueSrc{reasons=append(reasons,fmt.Sprintf("unique-sources %d >= %d",mu,cfg.AlertUniqueSrc))}
 if len(reasons)==0{return AttackEvent{},false};sev:="medium";if ap>=float64(cfg.AlertPPS)*4||ab>=float64(cfg.AlertBPS)*4{sev="critical"}else if ap>=float64(cfg.AlertPPS)*2||ab>=float64(cfg.AlertBPS)*2{sev="high"}
 return AttackEvent{ID:fmt.Sprintf("%d",now.UnixNano()),Time:now.UTC(),HostID:host,Severity:sev,Active:true,Reason:join(reasons),PacketsPerSec:ap,BitsPerSec:ab,UniqueSources:mu,SynRatio:as,Interfaces:interfaces},true}
func avg(v []float64)float64{if len(v)==0{return 0};var x float64;for _,n:=range v{x+=n};return x/float64(len(v))}
func maxInt(v []int)int{m:=0;for _,n:=range v{if n>m{m=n}};return m}
func join(v []string)string{if len(v)==0{return ""};o:=v[0];for _,s:=range v[1:]{o+="; "+s};return o}
