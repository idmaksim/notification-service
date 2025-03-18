package consumer

func (c *Consumer) Stop() {
	close(c.done)
	if c.ch != nil {
		c.ch.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}
