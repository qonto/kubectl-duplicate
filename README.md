# kubectl duplicate

This app is a plugin for `kubectl`, it allows you to duplicate a running Pod and auto-exec into. The list of Pods is filterable, and you can select the namespace you want.

You can also set this parameters for customization of the duplicata:
 - `cpu`
 - `memory`
 - `ttl`
 - `shell`

Already created duplicatas remain 4h (by default) and you can exec into them as long they're running.

## Requirements

### For build

- `GoLang v1.16`

### For Usage

- `kubectl`

## Usage

### Help

```shell
usage: kubectl-duplicate [<flags>]

Flags:
      --help                 Show context-sensitive help (also try --help-long and --help-man).
      --all-namespaces       All Namespace
  -t, --ttl=14400            Time to live of pods is seconds
  -n, --namespace="default"  Namespace
  -p, --pod=POD              Pod
  -s, --shell="sh"           Shell to use
  -c, --cpu=CPU              cpu
  -m, --memory=MEMORY        Memory
  -k, --kubeconfig=$HOME/.kube/config  
                             Kube config file (override by env var KUBECONFIG
  -v, --version              Print version
```

### Install

Download latest release from https://github.com/qonto/kubectl-duplicate/releases and extract `kubectl-duplicate` into your `/usr/loca/bin` and run `chmod -x /usr/local/bin/kubectl-duplicate`.

### Build

```shell
git clone https://github.com/qonto/kubectl-qonto.git
cd ./kubectl-duplicate
go build
mv kubectl-duplicate /usr/local/bin/
```

For `MacOSC` user:
```shell
xattr -d com.apple.quarantine /usr/local/bin/kubectl-duplicate
```

### Run 

* List pods:

    ```shell
    $ kubectl duplicate
    Search:
    ? Pods: 
      > falcosidekick-5f44cb5bff-94sqc
        falcosidekick-5f44cb5bff-jh9wk
    ```

* List pods with already created duplicatas:

    ```shell
    Search: █
    ? Pods: 
        falcosidekick-duplicata-nglvr-f569r [duplicata]
      > falcosidekick-duplicata-kzx9z-kpjh6 [duplicata]
        falcosidekick-duplicata-mtb9x-lb29p [duplicata]
        falcosidekick-5f44cb5bff-94sqc
        falcosidekick-5f44cb5bff-jh9wk
        falcosidekick-ui-867f5d6f7-76lfx

    End: 2021-03-08 01:06:25
    ```

* Filter:

    ```shell
    Search: 94█
    ? Pods: 
      > falcosidekick-5f44cb5bff-94sqc
    ```

## Author

[Thomas Labarussias](https://github.com/Issif)