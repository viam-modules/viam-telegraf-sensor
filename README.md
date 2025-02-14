# `telegraf-sensor` modular resource

A Viam `sensor` implementation in Go that reads the output from [Telegraf](https://github.com/influxdata/telegraf).
Currently, this sensor executes telegraf as a client and collect the metrics enabled on [viam-telegraf.conf](viam-telegraf.conf). 

> [!NOTE]
> Before configuring your `telegraf-sensor`, you must [create a machine](https://docs.viam.com/manage/fleet/machines/#add-a-new-machine).

Navigate to the **Config** tab of your machine's page in [the Viam app](https://app.viam.com/).
Click on the **Components** subtab and click **Create component**.
Select the `sensor` type, then select the `viam:viam-sensor:telegrafsensor` model.
Click **Add module**, then enter a name for your sensor and click **Create**.


> [!NOTE]
> For more information, see [Configure a Machine](https://docs.viam.com/configure/).

## Build and run

To use this module, follow the instructions to [add a module from the Viam Registry](https://docs.viam.com/registry/configure/#add-a-modular-resource-from-the-viam-registry) and select the `viam:viam-sensor:telegrafsensor` model from the [`viam-telegraf-sensor` module](https://app.viam.com/module/viam/viam-telegraf-sensor).

This sensor will attempt to automatically setup Telegraf on your device using `apt-get` on Linux or `homebrew` on Mac OS. It has been tested on the following devices 
* Raspberry 4/5 running Debian Bookworm
* Raspberry 4 running Debian Bullseye
* Orange Pi Zero 3 running OrangeArch 23.07
* Orange Pi Zero 3 running Ubuntu Jammy
* Jetson Orin Nano running Ubuntu Focal

## Configure your `telegraf-sensor`
On the new component panel, copy and paste the following attribute template into your JSON configuration:

```json
{
    "disable_cpu": false,
    "disable_disk": false,
    "disable_disk_io": false,
    "disable_kernel": false,
    "disable_mem": false,
    "disable_net": false,
    "disable_netstat": false,
    "disable_processes": false,
    "disable_swap": false,
    "disable_system": false,
    "disable_temp": true,
    "disable_wireless": true
}
```

Addjust your configuration and save your config.



### Attributes

The following metrics are enabled by default CPU, Disk, Disk IO, Kernel, Mem, Net, Netstat, Processes, Swap, System. Wireless and Temp are disabled by default. You can change this configuration by setting the following attributes accordingly:

| Name | Type | Inclusion | Default |
|---|---|---|---|
| disable_cpu | boolean | Optional | false |
| disable_disk | boolean | Optional | false |
| disable_disk_io | boolean | Optional | false |
| disable_kernel | boolean | Optional | false |
| disable_mem | boolean | Optional | false |
| disable_net | boolean | Optional | false |
| disable_netstat | boolean | Optional | false |
| disable_processes | boolean | Optional | false |
| disable_swap | boolean | Optional | false |
| disable_system | boolean | Optional | false |
| disable_temp | boolean | Optional | true |
| disable_wireless | boolean | Optional | true | 

### Example configuration

```json
{
  "name": "myststemsor",
  "model": "viam:viam-sensor:telegrafsensor",
  "type": "sensor",
  "namespace": "rdk",
  "attributes": {
    "disable_kernel": true
  },
  "depends_on": [],
  "service_configs": [
    {
      "type": "data_manager",
      "attributes": {
        "capture_methods": [
          {
            "method": "Readings",
            "additional_params": {},
            "capture_frequency_hz": 0.2
          }
        ]
      }
    }
  ]
}
```

## Local Development

To use the `viam-telegraf-sensor` module, clone this repository to your
machine’s computer, navigate to the `module` directory, and run:

```go
go build
```

On your robot’s page in the [Viam app](https://app.viam.com/), enter
the [module’s executable
path](/registry/create/#prepare-the-module-for-execution), then click
**Add module**.
The name must use only lowercase characters.
Then, click **Save config**.

## Next Steps

1. To test your sensor, go to the [**Control** tab](https://docs.viam.com/manage/fleet/robots/#control) and test that you are getting readings.
   <div class="highlight highlight-source-json notranslate position-relative overflow-auto" dir="auto">
   <details> 
    <summary>Example reading captured by the sensor</summary>
    <pre><code lang="json">
      {
        "readings": {
          "host": "raspi5agent",
          "diskio": {
            "write_time": 2608023,
            "io_time": 1889772,
            "write_bytes": 5211975680,
            "name": "mmcblk0p2",
            "iops_in_progress": 0,
            "merged_reads": 1998,
            "merged_writes": 637945,
            "read_bytes": 413373440,
            "weighted_io_time": 2640492,
            "read_time": 30239,
            "reads": 9984,
            "writes": 453041,
            "timestamp": 1711641414
          },
          "netstat": {
            "tcp_established": 68,
            "tcp_fin_wait2": 0,
            "tcp_syn_sent": 0,
            "tcp_time_wait": 0,
            "tcp_syn_recv": 0,
            "tcp_close_wait": 0,
            "tcp_close": 0,
            "udp_socket": 12,
            "tcp_last_ack": 0,
            "tcp_fin_wait1": 0,
            "tcp_none": 27,
            "tcp_closing": 0,
            "tcp_listen": 4,
            "timestamp": 1711641414
          },
          "system": {
            "n_cpus": 4,
            "uptime_format": "8 days, 20:05",
            "n_users": 0,
            "load5": 0.35,
            "timestamp": 1711641414,
            "n_unique_users": 0,
            "load15": 0.47,
            "uptime": 763516,
            "load1": 0.73
          },
          "net": [
            {
              "packets_sent": 0,
              "packets_recv": 0,
              "bytes_recv": 0,
              "bytes_sent": 0,
              "drop_out": 0,
              "drop_in": 0,
              "speed": -1,
              "err_out": 0,
              "timestamp": 1711641414,
              "interface": "eth0",
              "err_in": 0
            },
            {
              "err_out": 0,
              "bytes_recv": 5272812086,
              "packets_recv": 32067452,
              "err_in": 0,
              "drop_in": 2784928,
              "speed": -1,
              "interface": "wlan0",
              "drop_out": 0,
              "packets_sent": 1759907,
              "timestamp": 1711641414,
              "bytes_sent": 174888403
            }
          ],
          "swap": {
            "used": 0,
            "used_percent": 0,
            "in": 0,
            "out": 0,
            "timestamp": 1711641414,
            "free": 104808448,
            "total": 104808448
          },
          "disk": {
            "total": 30825463808,
            "inodes_total": 1849536,
            "inodes_free": 1789077,
            "used": 2455629824,
            "device": "mmcblk0p2",
            "fstype": "ext4",
            "timestamp": 1711641414,
            "used_percent": 8.314264523637558,
            "inodes_used_percent": 3.268873922973113,
            "path": "/",
            "free": 27079512064,
            "inodes_used": 60459
          },
          "wireless": {
            "interface": "wlan0",
            "timestamp": 1711641414,
            "status": 0,
            "retry": 5320,
            "noise": -256,
            "beacon": 0,
            "misc": 0,
            "link": 56,
            "crypt": 0,
            "frag": 0,
            "level": -54,
            "nwid": 0
          },
          "mem": {
            "swap_free": 104808448,
            "high_free": 0,
            "vmalloc_used": 17268736,
            "huge_pages_total": 0,
            "shared": 5226496,
            "swap_total": 104808448,
            "available_percent": 94.28952795578138,
            "sunreclaim": 31768576,
            "write_back": 0,
            "free": 7323566080,
            "mapped": 207929344,
            "used_percent": 4.489377016484977,
            "total": 8444952576,
            "cached": 656261120,
            "low_total": 0,
            "vmalloc_total": 69818585710592,
            "vmalloc_chunk": 0,
            "huge_pages_free": 0,
            "sreclaimable": 45105152,
            "commit_limit": 4327276544,
            "high_total": 0,
            "timestamp": 1711641414,
            "write_back_tmp": 0,
            "huge_page_size": 0,
            "swap_cached": 0,
            "dirty": 180224,
            "inactive": 353271808,
            "buffered": 85999616,
            "page_tables": 5341184,
            "committed_as": 1116209152,
            "active": 637059072,
            "low_free": 0,
            "available": 7962705920,
            "used": 379125760,
            "slab": 76873728
          },
          "kernel": {
            "context_switches": 2426899428,
            "entropy_avail": 256,
            "interrupts": 1279151542,
            "processes_forked": 24867,
            "timestamp": 1711641414,
            "boot_time": 1710877898
          },
          "processes": {
            "idle": 58,
            "paging": 0,
            "blocked": 0,
            "running": 0,
            "dead": 0,
            "total": 141,
            "timestamp": 1711641414,
            "stopped": 0,
            "total_threads": 255,
            "sleeping": 83,
            "unknown": 0,
            "zombies": 0
          },
          "cpu": {
            "usage_iowait": 0,
            "usage_nice": 0,
            "timestamp": 1711641414,
            "usage_steal": 0,
            "usage_guest": 0,
            "usage_softirq": 0,
            "usage_system": 1.0050251257723701,
            "usage_guest_nice": 0,
            "usage_irq": 0,
            "usage_idle": 98.49246232167037,
            "usage_user": 0.5025125628861851
          },
          "temp": [
            {
              "timestamp": 1711641414,
              "temp": 55.1,
              "sensor": "cpu_thermal"
            },
            {
              "temp": 50.823,
              "sensor": "rp1_adc",
              "timestamp": 1711641414
            }
          ]
        }
      } 
    </code></pre>
    </details> </div>

2. Once you can obtain your machine's performance metrics, configure the data manager to [capture](https://docs.viam.com/data/capture/) and [sync](https://docs.viam.com/data/cloud-sync/) the data from all of your machines.
3. To retrieve data captured with the data manager, you can [query data with SQL or MQL](https://docs.viam.com/data/query/) or [visualize it with tools like Grafana](https://docs.viam.com/data/visualize/).

## License
Copyright 2021-2023 Viam Inc. <br>
Apache 2.0
