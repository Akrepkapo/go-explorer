
func (r *RedisModel) Str() string {
	return fmt.Sprintf("%s:%d", r.Address, r.Port)
}

func (r *RedisModel) Initer() error {
	rc = redis.NewClient(&redis.Options{
		Addr:     r.Str(),
		Password: r.Password,
		DB:       0,
	})
	_, err := rc.Ping(ctx).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisModel) Conn() *redis.Client {
	return rc
}
func (l *RedisModel) Close() error {
	return rc.Close()
}
