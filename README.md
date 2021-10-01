# Spigot

A small utility to write synthetic AWS vpcflow logs to an S3 bucket.

The files are named "vpcflow_XXX_XXX.gz" where the Xs are numbers.

number of lines per file varies between 1024 and 131072 by powers of
2.

## Options

- `-b` name of the S3 bucket to write to.  REQUIRED
- `-r` AWS regions.  REQUIRED
- `-d` a Go duration, this is the interval to wait between generating
  files. default 60s
- `n` number of files to generate per interval. default 1
- `w` number of workers to generate files. default 1


## Asumptions

- Either aws credentials file or environment variables
  (AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY) are set.
- Credentials have rights to put an S3 object into the bucket.
