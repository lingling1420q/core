// +build integration

package ots3

import (
	"context"
	"fmt"
	"os"
)

func createBucket() {

}

func Example() {
	file, err := os.Open("./testdata/basn0g01-30.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	uploader := NewManager("Q3AM3UQ867SPQQA43P2F", "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG", "https://play.minio.io:9000", "asia", "mybucket")
	_ = uploader.CreateBucket(context.Background(), "mybucket")
	url, _ := uploader.Upload(context.Background(), "foo", file)
	fmt.Println(url)

	// Output:
	// https://play.minio.io:9000/mybucket/foo.png
}
