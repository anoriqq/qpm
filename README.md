# qpm
📦 Qanat Package Manager

## Installation
```sh
# Using go
go install github.com/anoriqq/qpm@latest

# Using sh
sh -c "$(curl -fsSL https://raw.githubusercontent.com/anoriqq/qpm/main/install)"
```

If needed, you can download binaries from [the releases page](https://github.com/anoriqq/qpm/releases).

## Usage
```txt
Qanat Package Manager

Usage:
  qpm [command]

Available Commands:
  aquifer     Manage aquifer
  completion  Generate the autocompletion script for the specified shell
  config      Manage qpm config
  help        Help about any command
  install     Install specifc package
  uninstall   Unnstall specifc package

Flags:
  -h, --help      help for qpm
  -v, --version   version for qpm

Use "qpm [command] --help" for more information about a command.

```

## Aquifer
Aquifer defines the package installation and uninstallation plans.  
It has the following directory structure.

```txt
.
└── pkgName.yml
```

Please refer to [the template](https://github.com/anoriqq/qpm/blob/6408dc7/testdata/foo.yml) for the yml file.
