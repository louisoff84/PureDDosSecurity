package agent

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type testOverride struct {
	mu      sync.RWMutex
	active  bool
	stats   PacketStats
	until   time.Time
	created time.Time
}

type testAttackRequest struct {
	PacketsPerSecond float64 `json:"packets_per_second"`
	BitsPerSecond    float64 `json:"bits_per_second"`
	SynPerSecond     float64 `json:"syn_per_second"`
	UdpPerSecond     float64 `json:"udp_per_second"`
	IcmpPerSecond    float64 `json:"icmp_per_second"`
	UniqueSources    int     `json:"unique_sources"`
	DurationSeconds  int     `json:"duration_seconds"`
}

func newTestOverride() *testOverride { return &testOverride{} }
func (t *testOverride) current(now time.Time) (PacketStats, bool) {
	t.mu.Lock(); defer t.mu.Unlock()
	if !t.active || !now.Before(t.until) { t.active=false; return PacketStats{}, false }
	return t.stats, true
}
func (t *testOverride) set(req testAttackRequest, now time.Time) {
	if req.DurationSeconds <= 0 { req.DurationSeconds=30 }
	t.mu.Lock(); defer t.mu.Unlock(); t.active=true
	t.stats=PacketStats{PacketsPerSecond:req.PacketsPerSecond,BitsPerSecond:req.BitsPerSecond,SynPerSecond:req.SynPerSecond,UdpPerSecond:req.UdpPerSecond,IcmpPerSecond:req.IcmpPerSecond,UniqueSources:req.UniqueSources,TcpPackets:uint64(req.PacketsPerSecond),UdpPackets:uint64(req.UdpPerSecond),IcmpPackets:uint64(req.IcmpPerSecond)}
	t.created=now; t.until=now.Add(time.Duration(req.DurationSeconds)*time.Second)
}
func (t *testOverride) clear(){t.mu.Lock();t.active=false;t.mu.Unlock()}
func (t *testOverride) status(now time.Time) map[string]any {
	t.mu.RLock(); defer t.mu.RUnlock(); enabled:=t.active&&now.Before(t.until); remaining:=int64(0); if enabled{remaining=int64(t.until.Sub(now).Seconds())}
	return map[string]any{"enabled":enabled,"remaining_secs":remaining,"created":t.created.UTC(),"until":t.until.UTC(),"warning":"Synthetic test data only; no packets are generated."}
}
func (a *Agent) registerTestRoutes(m *http.ServeMux){m.HandleFunc("/api/v1/test/attack",a.testAttack);m.HandleFunc("/api/v1/test/clear",a.testClear);m.HandleFunc("/api/v1/test/status",a.testStatus)}
func (a *Agent) testProtected(w http.ResponseWriter) bool { if a.cfg.APIAuthToken=="" { http.Error(w,"test API disabled: set PUREDDOS_API_TOKEN first",http.StatusForbidden); return false }; return true }
func (a *Agent) testAttack(w http.ResponseWriter,r *http.Request){
	if !a.testProtected(w){return}; if r.Method!=http.MethodPost{http.Error(w,"method not allowed",http.StatusMethodNotAllowed);return}
	req:=testAttackRequest{PacketsPerSecond:100000,BitsPerSecond:1000000000,SynPerSecond:80000,UniqueSources:2000,DurationSeconds:30}
	if r.ContentLength>0 { if err:=json.NewDecoder(r.Body).Decode(&req);err!=nil{http.Error(w,"invalid JSON body",http.StatusBadRequest);return} }
	if req.PacketsPerSecond<0||req.BitsPerSecond<0||req.SynPerSecond<0||req.UdpPerSecond<0||req.IcmpPerSecond<0||req.UniqueSources<0 {http.Error(w,"test values must be non-negative",http.StatusBadRequest);return}
	a.test.set(req,time.Now()); writeJSON(w,map[string]any{"ok":true,"test_mode":true,"duration_seconds":req.DurationSeconds,"stats":req,"warning":"Synthetic test only. No network traffic is generated."})
}
func (a *Agent) testClear(w http.ResponseWriter,r *http.Request){
	if !a.testProtected(w){return}; if r.Method!=http.MethodPost{http.Error(w,"method not allowed",http.StatusMethodNotAllowed);return}
	a.test.clear();a.detector.reset();a.testWasActive=false;a.attackMu.Lock();a.attackActive=false;a.attackMu.Unlock();writeJSON(w,map[string]any{"ok":true,"test_mode":false})
}
func (a *Agent) testStatus(w http.ResponseWriter,r *http.Request){if !a.testProtected(w){return};writeJSON(w,a.test.status(time.Now()))}
