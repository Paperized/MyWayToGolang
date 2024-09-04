package main

import (
	"fmt"
	"signals/signals"
	"strings"
)

func main() {
	words := []string{"hellos", "friendos", "ornitorinco"}
	/* 	Declare a signal, a signal has a value that can be set or retrieved, listeners can be attached to it
	It's a "primary" type of Signal meaning that it cannot depend on any other signal
	*/
	wordIndexSig := signals.MakeSignal(0)

	id2, _ := wordIndexSig.Listen(func(s int, bs *signals.BaseSignal[int]) {
		fmt.Printf("New fruit selected : %v!\n", s)
	})

	/*  Declare a computed signal, it has a value but cannot be changed manually since it's computed based on a provided function.
	 	It's a "secondary" type of Signal, the listeners are triggered based on some signals dependencies.
		When the values of the top change the value gets re-computed.
	*/
	wordCmpSig := signals.MakeComputedSignal(func() string {
		index := wordIndexSig.Get()
		return words[index]
	}, wordIndexSig)

	id, _ := wordCmpSig.Listen(func(s string, bs *signals.BaseSignal[string]) {
		fmt.Printf("New fruit selected : %v!\n", s)
	})

	// Not really a signal, its a computed value, in simple words a lazy loaded value based on a dirty flag. Optimized.
	// We listen to fruitsNameCmpSig to split the calculated string from fruitsNameCmpSig and split by "o"
	splitWordCmpValue := signals.MakeComputedValue(func() []string {
		return strings.Split(wordCmpSig.Get(), "o")
	}, wordCmpSig)

	// ComputedValue skips the index = 0 notification, the flag is set to dirty but nothing is computed until .Get() is called
	wordIndexSig.Set(0, true)

	wordIndexSig.Set(1)
	fmt.Println("Index: 1, ", splitWordCmpValue.Get())

	wordIndexSig.Set(2)
	fmt.Println("Index: 2, ", splitWordCmpValue.Get())

	wordCmpSig.UnlistenById(id)
	wordIndexSig.UnlistenById(id2)
}
