package conf

// 性能优化相关常量
const (
	// 数据库相关
	DefaultPageSize        = 15
	MaxSearchKeywordLength = 100
	MongoConnectTimeout    = 10000 // 10秒
	MongoMaxPoolSize       = 100   // 最大连接池大小
	MongoMinPoolSize       = 10    // 最小连接池大小

	// 缓存相关
	ImageCacheDir         = "upload/image/"
	DefaultCacheTTL       = 300 // 5分钟
	RandomContentCacheTTL = 300 // 5分钟

	// HTTP相关
	HTTPRequestTimeout = 30               // 30秒
	MaxResponseSize    = 10 * 1024 * 1024 // 10MB

	// 并发相关
	MaxConcurrentRequests = 100
	WorkerPoolSize        = 10

	// 健康检查相关
	HealthCheckInterval = 25 // 25秒
)

// 性能优化配置结构
type PerformanceConfig struct {
	EnableGzip         bool
	EnableBrotli       bool
	EnableCache        bool
	CacheSize          int
	MaxMemoryUsage     int64
	GCTargetPercentage int
}

// 默认性能配置
var DefaultPerformanceConfig = PerformanceConfig{
	EnableGzip:         false, // 当前禁用以避免解码问题
	EnableBrotli:       false, // 当前禁用以避免解码问题
	EnableCache:        true,
	CacheSize:          1000,
	MaxMemoryUsage:     512 * 1024 * 1024, // 512MB
	GCTargetPercentage: 100,
}
