use anyhow::{bail, Result};
use std::str::FromStr;
use tracing::info;

use crate::{
    path::settings_file,
    settings::SettingsFile,
    toolchain::{DistToolchainDescription, Toolchain},
    toolchain_override::ToolchainOverride,
};

pub fn default(toolchain: Option<String>) -> Result<()> {
    let current_toolchain = Toolchain::from_settings()?;

    let toolchain = match toolchain {
        Some(toolchain) => toolchain,
        None => {
            let mut message = String::new();

            if let Some(toolchain_override) = ToolchainOverride::from_project_root() {
                let override_name = toolchain_override
                    .description()
                    .map(|description| description.to_string())
                    .unwrap_or_else(|_| toolchain_override.cfg.toolchain.channel.to_string());
                message.push_str(&format!("{override_name} (override), "));
            }

            message.push_str(&format!("{} (default)", current_toolchain.name));
            info!("{message}");
            return Ok(());
        }
    };

    let new_default = match DistToolchainDescription::from_str(&toolchain) {
        Ok(desc) => Toolchain::from_path(&desc.to_string()),
        Err(_) => Toolchain::from_path(&toolchain),
    };

    if !new_default.exists() {
        bail!("Toolchain with name '{}' does not exist", &new_default.name);
    };

    let settings = SettingsFile::new(settings_file());
    settings.with_mut(|s| {
        s.default_toolchain = Some(new_default.name.clone());
        Ok(())
    })?;
    info!("default toolchain set to '{}'", new_default.name);

    Ok(())
}
