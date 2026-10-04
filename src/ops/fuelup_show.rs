use anyhow::Result;
use component::{self, Components};
use semver::Version;
use std::{io::Write, path::Path};
use std::str::FromStr;
use tracing::{error, info};

use crate::{
    config::Config,
    fmt::{bold, print_header},
    path::fuelup_dir,
    target_triple::TargetTriple,
    toolchain::{DistToolchainDescription, Toolchain},
    toolchain_override::ToolchainOverride,
};

fn exec_show_version(component_executable: &Path) -> Result<()> {
    match std::process::Command::new(component_executable)
        .arg("--version")
        .output()
    {
        Ok(o) => {
            let output = String::from_utf8_lossy(&o.stdout).into_owned();
            match output.split_whitespace().last() {
                Some(v) => {
                    let version = Version::parse(v)?;
                    info!(" : {}", version);
                }
                None => {
                    error!(" : Error getting version string");
                }
            };
        }
        Err(e) => {
            print!(" - ");
            if component_executable.exists() {
                error!("execution error - {}", e);
            } else {
                error!("not found");
            }
        }
    }

    Ok(())
}

pub fn show() -> Result<()> {
    bold(|s| write!(s, "Default host: "));
    info!("{}", TargetTriple::from_host()?);

    bold(|s| write!(s, "fuelup home: "));
    info!("{}\n", fuelup_dir().display());

    print_header("installed toolchains");
    let cfg = Config::from_env()?;
    let mut active_toolchain = Toolchain::from_settings()?;
    let toolchain_override = ToolchainOverride::from_project_root();
    let override_name = toolchain_override
        .as_ref()
        .map(|toolchain_override| {
            DistToolchainDescription::from_str(
                &toolchain_override.cfg.toolchain.channel.to_string(),
            )
            .map(|description| description.to_string())
            .unwrap_or_else(|_| toolchain_override.cfg.toolchain.channel.to_string())
        });

    for toolchain in cfg.list_toolchains()? {
        let mut message = toolchain.clone();

        if toolchain == active_toolchain.name {
            message.push_str(" (default)");
        }
        if Some(&toolchain) == override_name.as_ref() {
            message.push_str(" (override)");
        }

        info!("{message}");
    }

    print_header("\nactive toolchain");

    let active_toolchain_message = match toolchain_override {
        Some(toolchain_override) => {
            let override_name = override_name.expect("override name is available");
            active_toolchain = Toolchain::from_path(&override_name);
            format!(
                "{override_name} (override), path: {}",
                toolchain_override.path.display()
            )
        }
        None => format!("{} (default)", active_toolchain.name),
    };

    info!("{active_toolchain_message}");

    for component in Components::collect_exclude_plugins()? {
        bold(|s| write!(s, "  {}", &component.name));
        let component_executable = active_toolchain.bin_path.join(&component.name);
        exec_show_version(component_executable.as_path())?;

        if component.name == component::FORC {
            for plugin in Components::collect_plugins()? {
                bold(|s| write!(s, "    - {}", &plugin.name));
                if !plugin.is_main_executable() {
                    info!("");
                    for executable in plugin.executables.iter() {
                        bold(|s| write!(s, "      - {}", &executable));
                        let plugin_executable = active_toolchain.bin_path.join(executable);
                        exec_show_version(plugin_executable.as_path())?;
                    }
                } else {
                    let plugin_executable = active_toolchain.bin_path.join(&plugin.name);
                    exec_show_version(plugin_executable.as_path())?;
                }
            }
        }
    }

    Ok(())
}
