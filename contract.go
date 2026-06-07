package main

import "unsafe"

// NEAR host functions
//go:wasmimport env log_utf8
func log_utf8(len uint64, ptr uint64)

//go:wasmimport env value_return
func value_return(value_len uint64, value_ptr uint64)

//go:export get_message
func get_message() {
	msg := []byte("Hello from Go on NEAR!")
	ptr := uint64(uintptr(unsafe.Pointer(&msg[0])))
	log_utf8(uint64(len(msg)), ptr)
	value_return(uint64(len(msg)), ptr)
}

//go:export set_message
func set_message() {
	msg := []byte("Message was updated!")
	ptr := uint64(uintptr(unsafe.Pointer(&msg[0])))
	log_utf8(uint64(len(msg)), ptr)
}

func main() {}
