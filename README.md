# cloud-storage-analyzer

Use this tool to monitor the sizes of the Google Cloud Storage Buckets over time.


# Details

## Access Logs & Storage Logs
A way to get a daily report of your bucket’s statistics is the [Access Logs & Storage Logs](https://cloud.google.com/storage/docs/access-logs) for Google Cloud Storage.
Google Cloud Storage offers access logs and storage logs in the form of CSV files that you can download and view.
Access logs provide information for all of the requests made on a specified bucket and are created hourly,
while the daily storage logs provide information about the storage consumption of that bucket for the last day.

### Enable Access & Storage Logs

This is a sample command to enable acccess & storage logs on `sample-bucket` bucket in project `n-project`

```
gsutil logging set on -b gs://storage-logs-bucket -o n_project_sample_bucket gs://sample-bucket
```

The above commands stores the csv files in the `storage-logs-bucket` bucket.

## Cloud Storage Analyzer Job

Once the access & storage logs have been enabled on a bucket
* csv files with `_usage` substring are added to the storage bucket on every access
* csv files with `_storage` substring are added at the end of each day with the bucket size details.

The usage files are really big (each being 100s of MBs) and we dont really need them for determining the size.

We are only interested in the storage files which contain the bucket size.
The storage file looks like this:
```
"bucket","storage_byte_hours"
"sample-bucket","9939136570712692"
```

To get the actual size (in bytes), divide storage_byte_hours by 24.

### What this job does:
On each run,

* It deletes all the `_usage` access files from the storage bucket.
* It reads the `_storage` storage files and pushes a summary csv file in the same storage bucket.
  * This summary file looks like this:
  ```
    project,bucket,size,timestamp
  n-project,sample-bucket,414130690443239,2018-07-17 07:00:00
  n-project,sample-bucket,414130690442298,2018-07-18 07:00:00
  n-project,sample-bucket,414130690446362,2018-07-16 07:00:00
  ```
  * Note that the job goes through all storage files each time and recreates the summary file each time. This can be improved.
* The job runs every 30 mins in our case as a kubernetes cron job

## Cloud Data Studio

You can create a dashboard in cloud data studio and add a datasource that points to the summary file.
This will help to graph the sizes of the monitored buckets over time.

## More Reading

This setup is inspired from [this](https://blog.doit-intl.com/carefully-calculating-google-cloud-storage-buckets-size-with-cloud-functions-pub-sub-89ac34676629) blog.
Please read through it to understand why we can not use stackdriver to monitor the bucket sizes.






