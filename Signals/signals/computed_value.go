package signals

type ComputedValue[T any] struct {
	value T
	dirty bool
	mapFn func() T
}

func MakeComputedValue[T any](mapFn func() T, dependsOn ...signalSender) *ComputedValue[T] {
	if mapFn == nil || len(dependsOn) == 0 {
		return nil
	}

	cmpValue := &ComputedValue[T]{
		mapFn: mapFn,
		dirty: true,
	}

	for _, dep := range dependsOn {
		dep.AddBelowDependency(cmpValue)
	}

	return cmpValue
}

func (cs *ComputedValue[T]) DependencyChanged() {
	cs.dirty = true
}

func (cs *ComputedValue[T]) TriggerEvent() {}

func (cs *ComputedValue[T]) Get() T {
	if !cs.dirty {
		return cs.value
	}

	cs.value = cs.mapFn()
	cs.dirty = false
	return cs.value
}
