package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type setRequest struct { Key string `json:"key"`; Value interface{} `json:"value"`; TTLSeconds int `json:"ttl_seconds"` }
func (s *Server) cacheHandler(w http.ResponseWriter, r *http.Request) { start:=time.Now(); s.metrics.RecordRequest(); defer s.metrics.RecordLatency(start); switch r.Method { case http.MethodGet: s.get(w,r); case http.MethodPost, http.MethodPut: s.set(w,r); case http.MethodDelete: s.delete(w,r); default: writeError(w,"method not allowed",http.StatusMethodNotAllowed) } }
func (s *Server) get(w http.ResponseWriter,r *http.Request) { key:=r.URL.Query().Get("key"); if key=="" { writeError(w,"missing key",400); return }; value,ok:=s.cache.Get(key); if !ok { s.metrics.RecordMiss(); writeError(w,"not found",404); return }; s.metrics.RecordHit(); writeJSON(w,map[string]interface{}{ "key":key,"value":value},200) }
func (s *Server) set(w http.ResponseWriter,r *http.Request) { var req setRequest; if r.Header.Get("Content-Type")=="application/json" { if err:=json.NewDecoder(r.Body).Decode(&req); err != nil { writeError(w,"invalid JSON",400); return } } else { req.Key=r.URL.Query().Get("key"); req.Value=r.URL.Query().Get("value"); req.TTLSeconds,_=strconv.Atoi(r.URL.Query().Get("ttl")) }; if req.Key=="" { writeError(w,"missing key",400); return }; expiry:=time.Time{}; if req.TTLSeconds>0 { expiry=time.Now().Add(time.Duration(req.TTLSeconds)*time.Second) }; s.cache.SetWithExpiry(req.Key,req.Value,expiry); if s.aof!=nil { if err:=s.aof.AppendSet(req.Key,req.Value,expiry); err!=nil { writeError(w,"persistence failure",500); return } }; writeJSON(w,map[string]interface{}{ "key":req.Key,"value":req.Value,"ttl_seconds":req.TTLSeconds},200) }
func (s *Server) delete(w http.ResponseWriter,r *http.Request) { key:=r.URL.Query().Get("key"); if key=="" { writeError(w,"missing key",400); return }; if !s.cache.Delete(key) { s.metrics.RecordMiss(); writeError(w,"not found",404); return }; if s.aof!=nil { if err:=s.aof.AppendDelete(key); err!=nil { writeError(w,"persistence failure",500); return } }; s.metrics.RecordHit(); writeJSON(w,map[string]string{"message":"deleted"},200) }
func (s *Server) metricsHandler(w http.ResponseWriter,r *http.Request) { snap:=s.metrics.Snapshot(); w.Header().Set("Content-Type","text/plain; version=0.0.4"); fmt.Fprintf(w,"cache_hits_total %d\ncache_misses_total %d\ncache_evictions_total %d\ncache_requests_total %d\n",snap.Hits,snap.Misses,snap.Evictions,snap.Requests); for i,b:=range snap.Bounds { fmt.Fprintf(w,"cache_request_duration_seconds_bucket{le=\"%g\"} %d\n",b,snap.Buckets[i]) } }
func writeJSON(w http.ResponseWriter,data interface{},status int) { w.Header().Set("Content-Type","application/json"); w.WriteHeader(status); _=json.NewEncoder(w).Encode(data) }
func writeError(w http.ResponseWriter,msg string,status int) { writeJSON(w,map[string]string{"error":msg},status) }
