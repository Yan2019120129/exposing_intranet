package gorm

// Config gorm 配置
type Config struct {
	QueryFields   bool `yaml:"query-fields"`   // 是否开启全字段匹配查询
	SingularTable bool `yaml:"singular-table"` // 是否关闭数据库表复数s
}

func (c *Config) GetQueryFields() bool {
	return c.QueryFields
}
func (c *Config) SetQueryFields(queryFields bool) *Config {
	c.QueryFields = queryFields
	return c
}

func (c *Config) GetSingularTable() bool {
	return c.SingularTable
}
func (c *Config) SetSingularTable(singularTable bool) *Config {
	c.SingularTable = singularTable
	return c
}
