package agent

import("os";"strconv";"time")
type Config struct {
 ListenAddr string
 CollectorURL string
 CollectorToken string
 HostID string
 SampleInterval time.Duration
 Window time.Duration
 AlertPPS uint64
 AlertBPS uint64
 AlertSYNRatio float64
 AlertUniqueSrc int
 PacketSnapLen int
 APIAuthToken string
 CORSOrigin string
}
func LoadConfig() Config { return Config{
 ListenAddr:env("PUREDDOS_LISTEN","0.0.0.0:2456"),
 CollectorURL:os.Getenv("PUREDDOS_COLLECTOR"),
 CollectorToken:os.Getenv("PUREDDOS_TOKEN"),
 HostID:env("PUREDDOS_HOST_ID",hostname()),
 SampleInterval:envDuration("PUREDDOS_INTERVAL",time.Second),
 Window:envDuration("PUREDDOS_WINDOW",10*time.Second),
 AlertPPS:envUint64("PUREDDOS_ALERT_PPS",25000),
 AlertBPS:envUint64("PUREDDOS_ALERT_BPS",100000000),
 AlertSYNRatio:envFloat("PUREDDOS_ALERT_SYN_RATIO",0.70),
 AlertUniqueSrc:envInt("PUREDDOS_ALERT_UNIQUE_SRC",500),
 PacketSnapLen:128,
 APIAuthToken:os.Getenv("PUREDDOS_API_TOKEN"),
 CORSOrigin:os.Getenv("PUREDDOS_CORS_ORIGIN"),
}}
func env(k,d string)string{if v:=os.Getenv(k);v!=""{return v};return d}
func envDuration(k string,d time.Duration)time.Duration{if v:=os.Getenv(k);v!=""{if x,e:=time.ParseDuration(v);e==nil{return x}};return d}
func envUint64(k string,d uint64)uint64{if v:=os.Getenv(k);v!=""{if x,e:=strconv.ParseUint(v,10,64);e==nil{return x}};return d}
func envInt(k string,d int)int{if v:=os.Getenv(k);v!=""{if x,e:=strconv.Atoi(v);e==nil{return x}};return d}
func envFloat(k string,d float64)float64{if v:=os.Getenv(k);v!=""{if x,e:=strconv.ParseFloat(v,64);e==nil{return x}};return d}
func hostname()string{h,e:=os.Hostname();if e!=nil{return"unknown"};return h}
