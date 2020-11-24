package storage

func (c *CentrifugoConfig) Initer() error {
	if c.Enable {
		publisher = gocent.New(gocent.Config{
			Addr: c.URL,
			Key:  c.Key,
		})
	}
	return nil
}

func (c *CentrifugoConfig) Conn() *gocent.Client {
	return publisher
}

func (l *CentrifugoConfig) Close() error {
	return nil
}
