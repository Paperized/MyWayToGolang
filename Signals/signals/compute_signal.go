package signals

import (
	"errors"
	"fmt"
)

type ComputedSignal[T comparable] struct {
	value T
	mapFn func() T

	listeners         map[string]ListenerWrapper[T]
	belowDependencies []signalReceiver
}

func MakeComputedSignal[T comparable](mapFn func() T, dependsOn ...signalSender) *ComputedSignal[T] {
	if mapFn == nil || len(dependsOn) == 0 {
		return nil
	}

	cmpSignal := &ComputedSignal[T]{
		mapFn:             mapFn,
		listeners:         map[string]ListenerWrapper[T]{},
		belowDependencies: []signalReceiver{},
	}

	for _, dep := range dependsOn {
		dep.AddBelowDependency(cmpSignal)
	}

	return cmpSignal
}

func (cs *ComputedSignal[T]) AddBelowDependency(sr signalReceiver) {
	cs.belowDependencies = append(cs.belowDependencies, sr)
}

func (cs *ComputedSignal[T]) DependencyChanged() {
	cs.value = cs.mapFn()

	for _, dep := range cs.belowDependencies {
		dep.DependencyChanged()
	}
}

func (cs *ComputedSignal[T]) TriggerEvent() {
	var bs BaseSignal[T] = cs

	// wrap in goroutine
	for _, lsWrapper := range cs.listeners {
		if lsWrapper.isAsync {
			lsWrapper.listener(cs.value, &bs)
		}
	}

	for _, lsWrapper := range cs.listeners {
		if !lsWrapper.isAsync {
			lsWrapper.listener(cs.value, &bs)
		}
	}

	for _, dep := range cs.belowDependencies {
		dep.TriggerEvent()
	}
}

func (cs *ComputedSignal[T]) Get() T {
	return cs.value
}

func (cs *ComputedSignal[T]) ListenByEvent(listener *ListenerEvent[T], id ...string) (string, error) {
	if listener == nil {
		return "", errors.New("listener is null")
	}

	var newId string
	if len(id) > 0 {
		newId = id[0]
	} else {
		newId = fmt.Sprintf("%v", listener)
	}

	cs.listeners[newId] = ListenerWrapper[T]{listener: *listener, isAsync: false}
	return newId, nil
}

func (cs *ComputedSignal[T]) Listen(listener func(T, *BaseSignal[T]), id ...string) (string, error) {
	if listener == nil {
		return "", errors.New("listener is null")
	}

	event := MakeEventListener(listener)
	return cs.ListenByEvent(event, id...)
}

func (cs *ComputedSignal[T]) ListenAsyncByEvent(listener *ListenerEvent[T], id ...string) (string, error) {
	if listener == nil {
		return "", errors.New("listener is null")
	}

	var newId string
	if len(id) > 0 {
		newId = id[0]
	} else {
		newId = fmt.Sprintf("%v", listener)
	}

	cs.listeners[newId] = ListenerWrapper[T]{listener: *listener, isAsync: true}
	return newId, nil
}

func (cs *ComputedSignal[T]) ListenAsync(listener func(T, *BaseSignal[T]), id ...string) (string, error) {
	if listener == nil {
		return "", errors.New("listener is null")
	}

	event := MakeEventListener(listener)
	return cs.ListenAsyncByEvent(event, id...)
}

func (cs *ComputedSignal[T]) Unlisten(listener *ListenerEvent[T]) {
	id := fmt.Sprintf("%v", listener)
	delete(cs.listeners, id)
}

func (cs *ComputedSignal[T]) UnlistenById(id string) {
	delete(cs.listeners, id)
}

func (cs *ComputedSignal[T]) UnlistenAll() {
	cs.listeners = map[string]ListenerWrapper[T]{}
}
