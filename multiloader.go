package multiconfig

type multiLoader []Loader

// MultiLoader creates a loader that executes the loaders one by one in order
// and returns on the first error.
func MultiLoader(loader ...Loader) Loader {
	return multiLoader(loader)
}

// Load loads the source into the config defined by struct s
func (m multiLoader) Load(s interface{}) error {
	// just set it once here
	if err := setDefaults(s); err != nil {
		return err
	}

	for _, loader := range m {
		if err := loader.Load(s); err != nil {
			return err
		}
	}

	return nil
}
