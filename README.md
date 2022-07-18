# qpm
📦 Qanat Package Manager

## Installation
```sh
go install github.com/anoriqq/qpm@latest
```

## Usage
```txt
qpm

Usage:
  qpm [command]

Available Commands:
  aquifer     manage aquifer
  completion  Generate the autocompletion script for the specified shell
  config      update config
  help        Help about any command
  install     install packages
  version     show version info

Flags:
  -h, --help   help for qpm

Use "qpm [command] --help" for more information about a command.
```

## Aquifer
Aquifer defines the package installation and uninstallation plans.  
It has the following directory structure.

```txt
.
└── pkgname
    ├── v1.18.yml
    └── latest.yml
```

Please refer to [the template](https://github.com/anoriqq/qpm/blob/b600c503f98c4d68ac0428dc03e36505988c2826/template/pkgname/latest.yml) for the yml file.
