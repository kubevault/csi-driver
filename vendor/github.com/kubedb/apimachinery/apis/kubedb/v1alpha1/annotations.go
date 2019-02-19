package v1alpha1

const (
	DatabaseNamePrefix = "kubedb"

	GenericKey = "kubedb.com"

	LabelDatabaseKind = GenericKey + "/kind"
	LabelDatabaseName = GenericKey + "/name"

	PostgresKey         = ResourceSingularPostgres + "." + GenericKey
	ElasticsearchKey    = ResourceSingularElasticsearch + "." + GenericKey
	MySQLKey            = ResourceSingularMySQL + "." + GenericKey
	MongoDBKey          = ResourceSingularMongoDB + "." + GenericKey
	RedisKey            = ResourceSingularRedis + "." + GenericKey
	MemcachedKey        = ResourceSingularMemcached + "." + GenericKey
	EtcdKey             = ResourceSingularEtcd + "." + GenericKey
	SnapshotKey         = ResourceSingularSnapshot + "." + GenericKey
	LabelSnapshotStatus = SnapshotKey + "/status"

	AnnotationInitialized = GenericKey + "/initialized"
	AnnotationJobType     = GenericKey + "/job-type"

	PrometheusExporterPortNumber = 56790
	PrometheusExporterPortName   = "prom-http"

	JobTypeBackup  = "backup"
	JobTypeRestore = "restore"

	ElasticsearchRestPort     = 9200
	ElasticsearchRestPortName = "http"
	ElasticsearchNodePort     = 9300
	ElasticsearchNodePortName = "transport"

	RedisShardKey   = RedisKey + "/shard"
	RedisNodePort   = 6379
	RedisGossipPort = 16379
)
