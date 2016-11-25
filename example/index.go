package main

import (
	"flag"

	coral "github.com/coral"
	cache "github.com/coral/cache"
	config "github.com/coral/config"
	db "github.com/coral/db"
	log "github.com/coral/log"

	filter "github.com/coral/example/filter"
)

var _config config.Configer

/**
 * initRouter 方法，定义服务全部接口、参数校验、并指定过滤器链
 */
func initRouter(server *coral.Server) {
	r := coral.ParamChecker
	// /
	baseRouter := server.NewRouter("/", filter.Index)

	// /param?<params>
	paramRouter := baseRouter.NewRouter("param", filter.Param)
	// /param/check?a=1&b=2&c=1
	paramRouter.NewRouter(
		"check",
		r.Check(coral.V{
			"a": r.IsString,
			"b": r.IsInt,
			"c": r.IsBool}),
		filter.Param,
	)
	// TODO 更复杂的check

	// log
	baseRouter.NewRouter("log", filter.Log)

	// /mysql
	mysqlRouter := baseRouter.NewRouter("mysql", filter.Mysql)
	// /mysql/*
	mysqlRouter.NewRouter("select", filter.Select)
	mysqlRouter.NewRouter("insert", filter.Insert)
	mysqlRouter.NewRouter("update", filter.Update)
	mysqlRouter.NewRouter("transCommit", filter.TransCommit)
	mysqlRouter.NewRouter("transRollback", filter.TransRollback)

	// /redis
	redisRouter := baseRouter.NewRouter("redis", filter.Redis)
	redisRouter.NewRouter("set",
		r.Check(coral.V{
			"key": r.IsString,
			"val": r.IsInt}),
		filter.Set)
	redisRouter.NewRouter("get",
		r.Check(coral.V{
			"key": r.IsString}),
		filter.Get)
}

func initDB() {
	// add default db
	db.DB.AddDB(
		_config.Get("db.DEFAULT_DB"),
		_config.Get("db.DEFAULT_DB_DSN"),
		_config.Int("db.DEFAULT_DB_MAX_CONNECTION"),
		_config.Int("db.DEFAULT_DB_MAX_IDLE"))

	// add other db
	// ...
}

func initRedis() {
	// add default cache
	cache.Cache.AddRedis(
		_config.Get("cache.DEFAULT_REDIS"),
		_config.Get("cache.DEFAULT_REDIS_SERVER"),
		_config.Get("cache.DEFAULT_REDIS_AUTH"),
		_config.Int("cache.DEFAULT_REDIS_MAX_CONNECTION"),
		_config.Int("cache.DEFAULT_REDIS_MAX_IDLE"))

	// add other cache
	// ...
}

func initLog() {
	// add default logger
	log.Log.AddLogger(
		_config.Get("log.DEFAULT_LOG"),
		_config.Get("log.DEFAULT_LOG_PATH"),
		_config.Int("log.DEFAULT_LOG_MAX_NUMBER"),
		_config.Int64("log.DEFAULT_LOG_MAX_SIZE"),
		_config.Int("log.DEFAULT_LOG_MAX_LEVEL"),
		_config.Int("log.DEFAULT_LOG_MIN_LEVEL"))

	// add other logger
	// ...
}

func main() {
	conf := flag.String("ini", "", "your config file")
	flag.Parse()
	if *conf != "" {
		config.AddConfiger(config.INI, "config", *conf)
		_config = config.Use("config")

		// init log
		initLog()

		// init db
		initDB()

		// init redis
		initRedis()

		// new server
		server := coral.NewServer(_config.Get("server.HOST"))

		// new router
		initRouter(server)

		// start server
		server.Run()
	} else {
		panic("run with -h to find usage")
	}
}
