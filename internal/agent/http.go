package agent
import("encoding/json";"fmt";"net/http";"sync";"time")
type api struct{mu sync.RWMutex;snapshot Snapshot;events []AttackEvent}
func newAPI()*api{return &api{events:make([]AttackEvent,0,100)}}
func(a *api)update(s Snapshot,e *AttackEvent){a.mu.Lock();defer a.mu.Unlock();a.snapshot=s;if e!=nil{a.events=append(a.events,*e);if len(a.events)>100{a.events=a.events[len(a.events)-100:]}};a.snapshot.Events=append([]AttackEvent(nil),a.events...)}
func(a *api)routes(m *http.ServeMux){m.HandleFunc("/health",func(w http.ResponseWriter,r *http.Request){writeJSON(w,map[string]any{"ok":true,"service":"PureAntiDDoS","time":time.Now().UTC()})});m.HandleFunc("/api/v1/status",a.status);m.HandleFunc("/api/v1/events",a.eventsHandler);m.HandleFunc("/api/v1/metrics",a.metrics)}
func(a *api)status(w http.ResponseWriter,r *http.Request){a.mu.RLock();defer a.mu.RUnlock();writeJSON(w,a.snapshot)}
func(a *api)eventsHandler(w http.ResponseWriter,r *http.Request){a.mu.RLock();defer a.mu.RUnlock();writeJSON(w,a.events)}
func(a *api)metrics(w http.ResponseWriter,r *http.Request){a.mu.RLock();defer a.mu.RUnlock();fmt.Fprintf(w,"puredos_attack %d\npuredos_pps %f\npuredos_bps %f\npuredos_unique_sources %d\n",boolInt(a.snapshot.Attack),a.snapshot.Network.PacketsPerSecond,a.snapshot.Network.BitsPerSecond,a.snapshot.Network.UniqueSources)}
func writeJSON(w http.ResponseWriter,v any){w.Header().Set("Content-Type","application/json");_=json.NewEncoder(w).Encode(v)}
func boolInt(v bool)int{if v{return 1};return 0}
