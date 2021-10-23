# **telegraf-huawei-plugin**

## **Overview**
the huawei plugin for telegraf to collect and process information from huawei devices

## **Installation**
### **Prerequisites**

- OS : Ubuntu, CentOS, Suse, Windows, Red Hat
- Go : go1.17.2
- Telegraf : Telegraf (version 1.19 or later)
- Glibc : https://www.gnu.org/software/libc/libc.html
- Make : https://www.gnu.org/software/make/
- Protoc-gen-go :https://github.com/golang/protobuf

### **From Source**

$git clone https://github.com/HuaweiDatacomm/telegraf-huawei-plugin


### Build From Source

Telegraf requires Go version 1.17 or newer, the Makefile requires GNU make.

1. [Install Go](https://golang.org/doc/install) >=1.17 (1.17.2 recommended)
2. Clone the Telegraf repository:
   ```
   git clone https://github.com/influxdata/telegraf.git
   ```
3. Run install.sh(this file can prompt for users for the telegraf,users must enter the right path of telegraf)
   ```
   ./install.sh
4. Run `make` from the source directory
   ```
   cd telegraf
   make
   ```

## Getting Started

See usage with:

```shell
telegraf --help
```

#### Generate a telegraf config file:

```shell
telegraf config > telegraf.conf
```

#### Generate config with only cpu input & influxdb output plugins defined:

```shell
telegraf --section-filter agent:inputs:outputs --input-filter cpu --output-filter influxdb config
```

#### Run a single telegraf collection, outputting metrics to stdout:

```shell
telegraf --config telegraf.conf --test
```

#### Run telegraf with all plugins defined in config file:

```shell
telegraf --config telegraf.conf
```


