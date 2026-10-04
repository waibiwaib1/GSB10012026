# Overrides

`fuelup` automatically determines which [toolchain] to use when an installed
command such as `forc` is executed.

You can override the installed default toolchain by adding a
`fuel-toolchain.toml` file to a project directory. `fuelup` searches for this
file in the current directory and each parent directory.

## The Toolchain File

The `fuel-toolchain.toml` file lets a project select a distributed Fuel
toolchain. Only distributed toolchains are currently supported.

`latest` and `nightly` must include a date so the selected toolchain remains
reproducible:

```toml
[toolchain]
channel = "nightly-2023-01-09"
```

A beta channel can be selected directly:

```toml
[toolchain]
channel = "beta-2"
```

When the selected toolchain is not installed, the component proxies install it
automatically before executing the requested component.

Run `fuelup show` inside the project to display the active override and the
path to its `fuel-toolchain.toml` file.

[toolchain]: concepts/toolchains.md
