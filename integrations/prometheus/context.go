package prometheus

import (
	Context "context"
	"github.com/go-rod/rod"
	"github.com/vortex14/gotyphoon/ctx"
)

const (
	TypeSummaryVec = "summaryVec"
	TypeSummary    = "summary"
	TypeCounterVec = "counterVec"
	TypeCounter    = "counter"
	TypeGaugeVec   = "gaugeVec"
	TypeGauge      = "gauge"
)

const (
	RuntimeCpuNum           = "proxy_runtime_cpu_num"
	RuntimeCpuGoroutine     = "proxy_runtime_cpu_goroutine"
	RuntimeCpuCgoCall       = "proxy_runtime_cpu_cgo_call"
	RuntimeMemAlloc         = "runtime_mem_alloc"
	RuntimeMemAllocTotal    = "runtime_mem_alloc_total"
	RuntimeMemSys           = "runtime_mem_sys"
	RuntimeMemOthersys      = "runtime_mem_othersys"
	RuntimeMemLookups       = "runtime_mem_lookups"
	RuntimeMemMalloc        = "runtime_mem_malloc"
	RuntimeMemFrees         = "runtime_mem_frees"
	RuntimeMemHeapAlloc     = "runtime_mem_heap_alloc"
	RuntimeMemHeapSys       = "runtime_mem_heap_sys"
	RuntimeMemHeapIdle      = "runtime_mem_heap_idle"
	RuntimeMemHeapInuse     = "runtime_mem_heap_inuse"
	RuntimeMemHeapReleased  = "runtime_mem_heap_released"
	RuntimeMemHeapObjects   = "runtime_mem_heap_objects"
	RuntimeMemStackInuse    = "runtime_mem_stack_inuse"
	RuntimeMemStackSys      = "runtime_mem_stack_sys"
	RuntimeMemMspanInuse    = "runtime_mem_mspan_inuse"
	RuntimeMemMspanSys      = "runtime_mem_mspan_sys"
	RuntimeMemMcacheInuse   = "runtime_mem_mcache_inuse"
	RuntimeMemMcacheSys     = "runtime_mem_mcache_sys"
	RuntimeMemGCSys         = "runtime_mem_gc_sys"
	RuntimeMemGCNext        = "runtime_mem_gc_next"
	RuntimeMemGCLast        = "runtime_mem_gc_last"
	RuntimeMemGCPauseTotal  = "runtime_mem_gc_pause_total"
	RuntimeMemGCPause       = "runtime_mem_gc_pause"
	RuntimeMemGCNum         = "runtime_mem_gc_num"
	RuntimeMemGCCount       = "runtime_mem_gc_count"
	RuntimeMemGCCPUFraction = "runtime_mem_gc_cpu_fraction"
)

const (
	METRICS = "metrics"
)

func NewBrowserCtx(context Context.Context, browser *rod.Browser) Context.Context {
	return ctx.Update(context, METRICS, browser)
}

//func NewBodyResponse(context Context.Context, body *string) Context.Context {
//	return ctx.Update(context, RESPONSE, body)
//}
//
//func GetPageResponse(context Context.Context) (bool, *string) {
//	body, ok := ctx.Get(context, RESPONSE).(*string)
//	return ok, body
//}
//
//func GetBrowserCtx(context Context.Context) (bool, *rod.Browser) {
//	browser, ok := ctx.Get(context, BROWSER).(*rod.Browser)
//	return ok, browser
//}
//
//func NewPageCtx(context Context.Context, page *rod.Page) Context.Context {
//	return ctx.Update(context, PAGE, page)
//}
//
//func GetPageCtx(context Context.Context) (bool, *rod.Page) {
//	page, ok := ctx.Get(context, PAGE).(*rod.Page)
//	return ok, page
//}

