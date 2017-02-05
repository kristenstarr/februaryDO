package data

// Implementation of IndexStore interface to be used for testing purposes.
type TestStore struct {
	canAdd bool
	errAdd error
	canRemove bool
	errRemove error
	canHas bool
	errHas error
	canParents bool
	errParents error

}

func (t *TestStore) AddLibrary(name string, deps []string) (added bool, err error) {
	return t.canAdd, t.errAdd
}

func (t *TestStore) RemoveLibrary(name string) (removed bool, err error) {
	return t.canRemove, t.errRemove
}

func (t *TestStore) HasLibrary(name string) (exists bool, err error) {
	return t.canHas, t.errHas
}

func (t *TestStore) HasParents(name string) (hasParents bool, err error) {
	return t.canParents, t.errParents
}

// Creates a new IndexStore to be used for testing.
func NewTestStore(canAdd bool, errAdd error, canRemove bool, errRemove error, canHas bool, errHas error, canParents bool, errParents error) (IndexStore) {
	return &TestStore{
		canAdd, errAdd, canRemove, errRemove, canHas, errHas, canParents, errParents,
	}
}
