---
title: Csi-Vault Init
menu:
  product_csi_vault_0.1.0-alpha.1:
    identifier: csi-vault-init
    name: Csi-Vault Init
    parent: reference
product_name: pharmer
menu_name: product_csi_vault_0.1.0-alpha.1
section_menu_id: reference
---
## csi-vault init

Initializes the driver.

### Synopsis

Initializes the driver.

```
csi-vault init [flags]
```

### Options

```
      --endpoint string   CSI endpoint (default "unix:///var/lib/kubelet/plugins/com.vault.csi.vaultdbs/csi.sock")
  -h, --help              help for init
      --node string       Linode Hostname (default "pc")
      --token string      Vault access token
      --url string        Vault API URL (default "https://api.vault.com/")
```

### Options inherited from parent commands

```
      --alsologtostderr                  log to standard error as well as files
      --log_backtrace_at traceLocation   when logging hits line file:N, emit a stack trace (default :0)
      --log_dir string                   If non-empty, write log files in this directory
      --logtostderr                      log to standard error instead of files
      --stderrthreshold severity         logs at or above this threshold go to stderr (default 2)
  -v, --v Level                          log level for V logs
      --vmodule moduleSpec               comma-separated list of pattern=N settings for file-filtered logging
```

### SEE ALSO

* [csi-vault](/docs/reference/csi-vault.md)	 - Vault CSI by Appscode - Start farms

