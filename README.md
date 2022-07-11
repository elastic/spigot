# Spigot

A small utility to write synthetic logs to different destinations.

Currently supported log formats are:

- AWS Firewall
- AWS vpcflow
- Common Log Format
- Cisco ASA
- Citrix CEF
- Fortinet Firewall
- Generic CEF
- Windows Event XML (winlog)

Currently supported destinations are:

- Local file
- AWS S3 bucket
- Syslog (TCP or UDP)
- Rally (ndjson to local file)

## Command Line Flags

- `-c` Path to configuration.  Default "./spigot.yml"
- `-r` Seed random number generator with current time.  Default false.


## Config file

A configuration file is required.  The configuration file is a list of
runner configurations.  Runner configurations consist of:

- generator object.  This contains the configuration for the
  generator.  See godoc for each generator for config options.

- output object.  This contains the configuration for the output.  See
  godoc for each output for config options.

- records.  An integer, which is the number of records to write each
  interval.

- interval (Optional)  A golang duration.  Which specifies the time
  between writing records.  If omitted then the runner is executed
  once.
  
Example:

```yaml
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

