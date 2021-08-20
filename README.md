# trdl

[![maintainability][maintainability-badge]][maintainability-link]
[![coverage][coverage-badge]][coverage-link]

[maintainability-badge]:    https://api.codeclimate.com/v1/badges/a95ed9e90acae45f40ee/maintainability
[maintainability-link]:     https://codeclimate.com/github/werf/trdl/maintainability
[coverage-badge]:           https://api.codeclimate.com/v1/badges/a95ed9e90acae45f40ee/test_coverage
[coverage-link]:            https://codeclimate.com/github/werf/trdl/test_coverage

## File structure

```shell
~/.trdl$ tree -a
.
├── config.yaml
├── .locks
│   ├── ...
│   └── repositories
│       ├── ...
│       └── <REPO>
├── logs
│   ├── ...
│   └── repositories
│       ├── ...
│       └── <REPO>
├── .tmp
│   ├── ...
│   └── repositories
│       ├── ...
│       └── <REPO>
└── repositories
    ├── ...
    └── <REPO>
        ├── channels
        │   ├── ...
        │   └── <GROUP>
        │       ├── ...
        │       └── <CHANNEL>
        ├── .meta
        ├── releases
        │   ├── ...
        │   └── <RELEASE_VERSION>
        │       └── <OS>_<ARCH>
        │           └── ...
        └── scripts
            ├── ...
            └── <GROUP>-<CHANNEL>
                ├── ...
                └── source_script[.<ext>]
```

```shell
~/.trdl$ cat config.yaml
repositories:
  ...
  - name: <REPO>
    url: <URL>
    defaultChannel: <CHANNEL>
```