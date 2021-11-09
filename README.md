# Spigot

A small utility to write synthetic logs to different destinations.

Currently supported log formats are:

- AWS vpcflow
- Cisco ASA

Currently supported destinations are:

- Local file
- AWS S3 bucket
- Syslog (TCP or UDP)


## Config file

Configuration file is required.  Here are the available options:

```yaml
---
# workers: Number of go routines generating logs.  Each worker has a separate
# destination.  example worker 2 would generate 2 log files for file output.
workers: 2

# records: number of log entries per output.  example records 5 would
# output 5 lines to file output.
records: 3

# interval: go Duration, time in between running workers.
interval: 10s

# ASA.  include_timestamp mimics cisco timestamp in message behavior
generator_asa:
  enabled: true
  include_timestamp: true

generator_vpcflow:
  enabled: false

# File.  Directory is where logs will be written. Pattern is for
# filename see os.CreateTemp.  Delimiter is string between records,
# normally newline
output_file:
  enabled: true
  directory: "/var/tmp"
  pattern: "spigot_*.log"
  delimiter: "\n"

# S3. bucket is name of S3 bucket to write to.  Region is region
# bucket is in.  Delimiter is string between records, normally
# newline.  Prefix is for generating random key name, after name is
# time in seconds & nano seconds.
output_s3:
  enabled: false
  bucket: "leh-test-spigot"
  region: "us-west-2"
  delimiter: "\n"
  prefix: "vpcflow"

# Syslog. See log.syslog for valid facility and severities.  see
# syslog.Dial for meaning of tag.  Network is udp or tcp.  host is
# hostname to connect to.  port is port the syslog server is listen
# on.
  
output_syslog:
  enabled: false
  facility: "LOG_LOCAL0"
  severity: "LOG_INFO"
  tag: "test"
  network: "udp"
  host: localhost
  port: "1234"

```


## Running

```
spigot
```

All config is taken from config file which is assumed to be
`spigot.yml` in same directory as executable.

# Assumptions

## S3 output
- Either aws credentials file or environment variables
  (AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY) are set.
- Credentials have rights to put an S3 object into the bucket.
