package prometheus

import (
	Context "context"
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
	METRICS = "METRICS"
)

func NewMetricsCtx(context Context.Context, metrics MetricsInterface) Context.Context {
	return ctx.Update(context, METRICS, metrics)
}

func GetMetricsCtx(context Context.Context) (bool, MetricsInterface) {
	metrics, ok := ctx.Get(context, METRICS).(MetricsInterface)
	return ok, metrics

}
