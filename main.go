package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
)

var storageLogsBucket, bucketSizesFile *string

func main() {

	flag.Parse()

	ctx := context.Background()

	// Creates a client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	bucket := client.Bucket(*storageLogsBucket)
	fmt.Println(bucket)
	//deleteDataUsageObjects(bucket)
	deleteObjectsWithSubstring("_usage", bucket)

	createDataSourceFile(bucket)
}

func createDataSourceFile(bucket *storage.BucketHandle) {
	ctx := context.Background()
	object := bucket.Object(*bucketSizesFile)
	substring := "_storage"

	recordsTitle := []string{"project", "bucket", "size", "timestamp"}
	ow := object.NewWriter(ctx)

	defer ow.Close()

	w := csv.NewWriter(ow)
	w.Write(recordsTitle)

	objects := bucket.Objects(ctx, nil)

	for {
		o, err := objects.Next()
		if err != nil {
			println(err.Error())
			break
		}

		if strings.Contains(o.Name, substring) {
			fmt.Println("Processing file ", o.Name)
			storageObject := bucket.Object(o.Name)
			reader, _ := storageObject.NewReader(ctx)
			csvReader := csv.NewReader(reader)

			for {
				record, err := csvReader.Read()
				if err == io.EOF {
					break
				}
				if err != nil {
					fmt.Println(err.Error())
				}

				if len(record) == 2 {
					storageRow := ParseStorageRow(record[0], o.Name, record[1])
					w.Write(storageRow)
				}
			}
		}
	}
	w.Flush()

}

//ebay_n_ai_data_storage_2018_07_16_07_00_00_01a3a_v0
func ParseStorageRow(bucket string, objectName string, size string) []string {

	//_ai_data_storage_
	bucketWithUnderscores := "_" + strings.Replace(bucket, "-", "_", -1) + "_storage_"

	data := strings.Split(objectName, bucketWithUnderscores)

	if len(data) < 2 {
		return []string{}
	}

	project := strings.Replace(data[0], "_", "-", -1)

	timestamp, _ := func() (time.Time, error) {
		timestampTemp := strings.Split(
			strings.Replace(data[1], "_", "-", 5),
			"_")[0]
		return time.Parse(
			"2006-01-02-15-04-05",
			timestampTemp)

	}()

	return []string{
		project,
		bucket,
		size,
		timestamp.Format("2006-01-02 15:04:05"),
	}
}

func deleteObjectsWithSubstring(substring string, bucket *storage.BucketHandle) {
	ctx := context.Background()
	objects := bucket.Objects(ctx, nil)
	var wg sync.WaitGroup

	for {
		o, err := objects.Next()
		if err != nil {
			println(err.Error())
			break
		}
		wg.Add(1)

		go func() {
			defer wg.Done()
			if strings.Contains(o.Name, substring) {
				object := bucket.Object(o.Name)
				attrs, err := object.Attrs(ctx)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Printf("Deleting object %s\n", attrs.Name)
					object.Delete(ctx)
				}
			}
		}()

	}
	// Wait for all subroutines to complete.
	wg.Wait()
}

func init() {
	storageLogsBucket = flag.String("storageLogsBucket", "n-storage-logs-bucket", "The Bucket where the storage logs are stored")
	bucketSizesFile = flag.String("bucketSizesFile", "storage-buckets-size.csv", "The CSV file to which the bucket sizes are dumped")
}
