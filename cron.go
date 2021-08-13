package main
//
//import (
//	"github.com/rcrick/blog-api/models"
//	"github.com/robfig/cron"
//	"log"
//	"time"
//)
//
//-api
//
//import (
//	"github.com/rcrick/blog-api/models"
//	"github.com/robfig/cron"
//	"log"
//	"time"
//)
//
//func s() {
//	c := cron.New()
//
//	c.AddFunc("* * * * * *", func() {
//		log.Println("Run models.CleanAllTag...")
//		models.CleanDeletedTag()
//	})
//	c.AddFunc("* * * * * *", func() {
//		log.Println("Run models.CleanAllArticle...")
//		models.CleanDeletedArticle()
//	})
//
//	c.Start()
//
//	t1 := time.NewTimer(time.Second * 10)
//	for {
//		select {
//		case <-t1.C:
//			t1.Reset(time.Second * 10)
//		}
//	}
//}
