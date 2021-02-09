# gocloudproxy

Small service for proxying HTTP to cloud storage using [gocloud](https://gocloud.dev/)

- Handles `If-None-Match` and `If-Modified-Since` headers, so browser caching works.
- Only does proxying, leaving auth, HTTPS, etc. to other components.

## Environment Variables

See also [gocloud blob docs](https://gocloud.dev/howto/blob/#services)

| Environment Variables | Description                                   | Required | Default |
| --------------------- | --------------------------------------------- | -------- | ------- |
| STORAGE_URL           | gocloud url for the storage                   | \*       |         |
| PORT                  | The port number to be assigned for listening. |          | 8080    |
| AWS_REGION            | The AWS `region` where the S3 bucket exists.  |          |         |
| AWS_ACCESS_KEY_ID     | AWS `access key` for S3.                      |          |         |
| AWS_SECRET_ACCESS_KEY | AWS `secret key` for S3.                      |          |         |
| AZURE_STORAGE_ACCOUNT | Storage account for Azure                     |          |         |
| AZURE_STORAGE_KEY     | Storage key for Azure                         |          |         |
