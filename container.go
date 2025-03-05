package mappath

// Container stores data and updates it only if change operations have been performed successfully.
type Container struct {
	Data any
}

func (c *Container) Get(key string) (any, error) {
	return Get(c.Data, key)
}

func (c *Container) Put(key string, val any) error {
	data, err := Put(c.Data, key, val)
	if err != nil {
		return err
	}

	c.Data = data
	return nil
}

func (c *Container) Delete(key string) error {
	data, err := Delete(c.Data, key)
	if err != nil {
		return err
	}

	c.Data = data
	return nil
}

func (c *Container) Clone() *Container {
	cc := &Container{}
	cc.Data = Clone(c.Data)
	return cc
}
