// Copyright (c) 2025 Klaus Post, released under MIT License. See LICENSE file.

//go:build riscv64 && !gccgo && !noasm && !appengine
// +build riscv64,!gccgo,!noasm,!appengine

package cpuid

func initCPU() {
	cpuid = func(uint32) (a, b, c, d uint32) { return 0, 0, 0, 0 }
}

func addInfo(c *CPUInfo, safe bool) {
	detectOS(c)
}

func getVectorLength() (vl, pl uint64) { return 0, 0 }
