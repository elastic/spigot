# Spigot

A small utility to write synthetic logs to different destinations.

Currently supported log formats are:

- AWS vpcflow
- Cisco ASA
- Fortinet Firewall

Currently supported destinations are:

- Local file
- AWS S3 bucket
- Syslog (TCP or UDP)


## Command Line Flags

- `-c` Path to configuration.  Default "./spigot.yml"
- `-r` Seed random number generator with current time.  Default false.


## Config file

A configuration file is required.  The configuration file is a list of
runner configurations.  Runner configurations consist of:

- generator object.  This contains the configuration for the generator.

- output object.  This contains the configuration for the output.

- records.  An integer, which is the number of records to write each
  interval.

- interval (Optional)  A golang duration.  Which specifies the time
  between writing records.  If omitted then the runner is executed
  once.
  
Example:

```yaml
---
---
runners:
  - generator:
      type: "cisco:asa"
      include_timestamp: false
    output:
      type: file
      directory: "/var/tmp"
      pattern: "spigot_asa_*.log"
      delimiter: "\n"
    interval: 5s
    records: 250
  - generator:
      type: "fortinet:firewall"
      include_timestamp: false
    output:
      type: file
      directory: "/var/tmp"
      pattern: "spigot_fortinet_firewall_*.log"
      delimiter: "\n"
    interval: 10s
    records: 2048
```


# Assumptions

## S3 output
- Either aws credentials file or environment variables
  (AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY) are set.
- Credentials have rights to put an S3 object into the bucket.
