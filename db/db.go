package db

var GlobalStore *BoxerStore = nil

func InitializeStore() {
	if GlobalStore == nil {
		GlobalStore = &BoxerStore{}
		GlobalStore.Container = make(map[BoxerKey]BoxerValue)
	}
}

func (s *BoxerStore) Get(Key BoxerKey) (BoxerValue, error) {
	if val, ok := s.Container[Key]; ok {
		return val, nil
	}
	return BoxerValue{}, ErrorKeyNotFound(Key)
}

func (s *BoxerStore) Put(Key BoxerKey, Value BoxerValue) {
	s.Container[Key] = Value
}

func (s *BoxerStore) Delete(Key BoxerKey) error {
	if _, ok := s.Container[Key]; ok {
		delete(s.Container, Key)
		return nil
	}
	return ErrorKeyNotFound(Key)
}
