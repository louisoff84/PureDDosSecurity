package agent
import("context";"encoding/json";"net/http";"strings";"time")
type collector struct{url,token string;http *http.Client}
func newCollector(c Config)*collector{return &collector{url:strings.TrimRight(c.CollectorURL,"/"),token:c.CollectorToken,http:&http.Client{Timeout:5*time.Second}}}
func(c *collector)send(ctx context.Context,path string,payload any){if c.url==""{return};b,e:=json.Marshal(payload);if e!=nil{return};r,e:=http.NewRequestWithContext(ctx,http.MethodPost,c.url+path,strings.NewReader(string(b)));if e!=nil{return};r.Header.Set("Content-Type","application/json");if c.token!=""{r.Header.Set("Authorization","Bearer "+c.token)};resp,e:=c.http.Do(r);if e==nil{resp.Body.Close()}}
