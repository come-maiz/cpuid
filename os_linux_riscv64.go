// Copyright (c) 2025 Klaus Post, released under MIT License. See LICENSE file.

// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file located
// here https://github.com/golang/sys/blob/master/LICENSE

package cpuid

import _ "unsafe" // needed for go:linkname
import (
	"golang.org/x/sys/unix"
	"runtime"
	"syscall"
	"unsafe"
)

const (
	// Copied from golang.org/x/sys/unix/zsysnum_linux_riscv64.go.
	sys_RISCV_HWPROBE = 258
)

// Copied from golang.org/x/sys/unix/ztypes_linux_riscv64.go.
type riscvHWProbePairs struct {
	key   int64
	value uint64
}

func detectOS(c *CPUInfo) bool {
	c.LogicalCores = runtime.NumCPU()
	// For now assuming 1 thread per core...
	c.ThreadsPerCore = 1
	c.LogicalCores = c.PhysicalCores

	pairs := []riscvHWProbePairs{
		{unix.RISCV_HWPROBE_KEY_MVENDORID, 0},
		{unix.RISCV_HWPROBE_KEY_MARCHID, 0},
	}
	if !riscvHWProbe(pairs, 0) {
		return false
	}
	if pairs[0].key != -1 && pairs[0].value == 0 && pairs[1].key != -1 {
		// https://github.com/riscv/riscv-isa-manual/blob/main/marchid.md
		if pairs[1].value == 42 {
			// There is an open patch for QEMU to return the right value:
			// https://lore.kernel.org/all/20240131182430.20174-1-palmer@rivosinc.com/
			c.VendorID = QEMU
		}
	}

	return true
}

// Copied from golang.org/x/sys/cpu/cpu_linux_riscv64.go
func riscvHWProbe(pairs []riscvHWProbePairs, flags uint) bool {
	var _zero uintptr
	var p0 unsafe.Pointer
	if len(pairs) > 0 {
		p0 = unsafe.Pointer(&pairs[0])
	} else {
		p0 = unsafe.Pointer(&_zero)
	}

	_, _, e1 := syscall.Syscall6(sys_RISCV_HWPROBE, uintptr(p0), uintptr(len(pairs)), uintptr(0), uintptr(0), uintptr(flags), 0)
	return e1 == 0
}
