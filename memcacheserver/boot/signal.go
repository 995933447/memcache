package boot

import (
	"sync"
	"syscall"
	"time"

	"github.com/995933447/easymicro/node"
	"github.com/995933447/easymicro/sysmon"
	"github.com/995933447/fastlog"
	"github.com/995933447/memcache/memcacheserver/config"
)

const (
	SignalAliasStop        = "stop"
	SignalAliasInterrupt   = "interrupt"
	SignalAliasSamplePProf = "sample_pprof"
	SignalAliasPrintMem    = "print_mem"
)

var (
	sig   *sysmon.OsSignal
	sigMu sync.RWMutex
)

func InitSignal() (*sysmon.OsSignal, error) {
	sigMu.Lock()
	defer sigMu.Unlock()

	if sig != nil {
		return sig, nil
	}

	sig = sysmon.NewOsSignal()
	sig.AliasSignal(syscall.SIGTERM, SignalAliasStop)
	sig.AliasSignal(syscall.SIGINT, SignalAliasInterrupt)
	sig.AliasSignal(syscall.SIGUSR1, SignalAliasSamplePProf)
	sig.AliasSignal(syscall.SIGUSR2, SignalAliasPrintMem)

	err := sig.AppendSignalCallbackByAlias(SignalAliasSamplePProf, func() {
		cpuFile := "./profile_cpu"
		heapFile := "./profile_mem"
		if node.GetName() != "" {
			cpuFile += "_" + node.GetName()
			heapFile += "_" + node.GetName()
		}
		var sampleTimeSecLong int
		config.SafeReadServerConfig(func(c *config.ServerConfig) {
			sampleTimeSecLong = c.SamplePProfTimeLongSec
		})
		if sampleTimeSecLong == 0 {
			sampleTimeSecLong = 20
		}
		if err := sysmon.DumpPProfiles(cpuFile, heapFile, time.Second*time.Duration(sampleTimeSecLong)); err != nil {
			fastlog.Fatal(err)
		}
	})
	if err != nil {
		return nil, err
	}

	err = sig.AppendSignalCallbackByAlias(SignalAliasPrintMem, func() {
		sysmon.PrintMemStats()
	})
	if err != nil {
		return nil, err
	}

	return sig, nil
}
