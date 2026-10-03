package agent
import "time"
type InterfaceStats struct {
 Name string `json:"name"`; Up bool `json:"up"`; RxBytes uint64 `json:"rx_bytes"`; TxBytes uint64 `json:"tx_bytes"`
 RxPackets uint64 `json:"rx_packets"`; TxPackets uint64 `json:"tx_packets"`; RxPPS float64 `json:"rx_pps"`; TxPPS float64 `json:"tx_pps"`
 RxMbps float64 `json:"rx_mbps"`; TxMbps float64 `json:"tx_mbps"`; Addresses []string `json:"addresses,omitempty"`
}
type SystemStats struct {
 CPUPercent float64 `json:"cpu_percent"`; MemoryPercent float64 `json:"memory_percent"`; MemoryUsed uint64 `json:"memory_used"`; MemoryTotal uint64 `json:"memory_total"`
 Load1 float64 `json:"load1"`; Load5 float64 `json:"load5"`; Load15 float64 `json:"load15"`; Goroutines int `json:"goroutines"`; UptimeSeconds uint64 `json:"uptime_seconds"`
}
type PacketStats struct {
 PacketsPerSecond float64 `json:"packets_per_second"`; BitsPerSecond float64 `json:"bits_per_second"`; SynPerSecond float64 `json:"syn_per_second"`
 UdpPerSecond float64 `json:"udp_per_second"`; IcmpPerSecond float64 `json:"icmp_per_second"`; UniqueSources int `json:"unique_sources"`
 TcpPackets uint64 `json:"tcp_packets"`; UdpPackets uint64 `json:"udp_packets"`; IcmpPackets uint64 `json:"icmp_packets"`
}
type AttackEvent struct {
 ID string `json:"id"`; Time time.Time `json:"time"`; HostID string `json:"host_id"`; Severity string `json:"severity"`; Active bool `json:"active"`
 Reason string `json:"reason"`; PacketsPerSec float64 `json:"packets_per_second"`; BitsPerSec float64 `json:"bits_per_second"`; UniqueSources int `json:"unique_sources"`
 SynRatio float64 `json:"syn_ratio"`; Interfaces []string `json:"interfaces"`
}
type Snapshot struct {
 Time time.Time `json:"time"`; HostID string `json:"host_id"`; Attack bool `json:"attack"`; System SystemStats `json:"system"`
 Interfaces []InterfaceStats `json:"interfaces"`; Network PacketStats `json:"network"`; Events []AttackEvent `json:"events,omitempty"`
}