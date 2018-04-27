package main

import (
	"flag"
	"log"
	"strings"
	"time"
	"fmt"
	"gopkg.in/olivere/elastic.v3"
	"github.com/Pallinder/go-randomdata"
)

// format of messages sent to ElasticSearch
type MyType struct {
  Time    time.Time `json:"@timestamp"`
  Message string    `json:"message"`
}

func main() {

    // load command line arguments
    server := flag.String("server","","ElasticSearch server e.g. http://localhost:9200")
    flag.Parse()


    // configure connection to ES
    client, err := elastic.NewClient(elastic.SetURL(*server))
    if err != nil {
      panic(err)
    }
    log.Printf("client.running? %v",client.IsRunning())
    if ! client.IsRunning() {
      panic("Could not make connection, not running")
    }


    // check ElasticSearch version
    log.Println("-------ElasticSearch version--------")
    version,verr := client.ElasticsearchVersion(*server)
    if verr != nil {
      panic(verr)
    }
    // make sure this version of API is suited to ES backend
    log.Printf("olivere/elastic API version: %s",elastic.Version)
    log.Printf("ElasticSearch server version: %s",version)
    if version[0:2] > "2." {
      panic(fmt.Sprintf("This API oliver/elastic.v3 is meant for ElasticSearch 2.x, and you are using %s.  Import 'gopkg.in/olivere/elastic.v5' instead",version))
    }

    t := strings.ToLower(randomdata.SillyName())
    count := 0
		for {


        if count == 1000000 {
				   count = 0
					 t = strings.ToLower(randomdata.SillyName())
				}
		    //log.Println("-------ElasticSearch insert--------")
		    // insert row of data into index=myindex, type=mytype
				msg := fmt.Sprintf(`%s %s %s`,randomdata.Paragraph(), randomdata.Paragraph(), time.Now())
				//log.Println(msg)
		    row := MyType{
		      Time: time.Now(),
		      Message: msg,
		     }

		     _,ierr := client.Index().Index("myindex"+"-"+t).Type("mytype").BodyJson(row).Refresh(true).Do()
		    if ierr != nil {
		       panic(ierr)
		     }
				 count++
		     //log.Printf("Successfully inserted row of data into myindex/mydata: %+v",ires)
         //time.Sleep(5 * time.Millisecond)

	 }
    /*
    log.Println("-------ElasticSearch search--------")
    termQuery := elastic.NewTermQuery("message", "inserted")
    res, err := client.Search("myindex").
        Index("myindex").
        Query(termQuery).
        Sort("@timestamp", true).
        Do()
    if err != nil {
        return
    }
    fmt.Println("Rows found:")
    var l MyType
    for _, item := range res.Each(reflect.TypeOf(l)) {
        l := item.(MyType)
        fmt.Printf("time: %s message: %s\n", l.Time, l.Message)
    }
    log.Println("done with search")
    */

}
