package main


//func main() {
//
//	gofk.Ignite().
//		Beans(db.NewGormAdapter(), cache.NewRedis()).
//		Attach(middlewares.NewUserMid()).
//		Mount("v1",
//			classes.NewUserClass(),
//			classes.NewArticleCalss()).
//		Cron("0/3 * * * * *", func() {
//			log.Println("每隔3秒执行一下子")
//		}).
//		Cron("0/2 * * * * *", func() {
//			log.Println("每隔2秒执行一下子")
//		}).
//		Launch()
//
//}
