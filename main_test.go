package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseStorageRow(t *testing.T) {

	bucket := "sample-bucket"
	objectName := "n_project_sample_bucket_storage_2018_07_16_07_00_00_01a3a_v0.dms"
	size := "9939136570712692"

	row := ParseStorageRow(bucket, objectName, size)

	assert.Equal(t, 4, len(row))

	assert.Equal(t, "n-project", row[0])
	assert.Equal(t, "sample-bucket", row[1])
	assert.Equal(t, size, row[2])
	assert.Equal(t, "2018-07-16 07:00:00", row[3])
}
