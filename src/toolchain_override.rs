use anyhow::{bail, Result};
use semver::Version;
use serde::de::Error as DeError;
use serde::{Deserialize, Deserializer, Serialize};
use std::collections::HashMap;
use std::fmt;
use std::path::PathBuf;
use std::str::FromStr;
use time::Date;
use toml_edit::{de, value, Document};
use tracing::warn;

use crate::channel::{is_beta_toolchain, LATEST, NIGHTLY};
use crate::constants::{DATE_FORMAT, FUEL_TOOLCHAIN_TOML_FILE};
use crate::file;
use crate::path::get_fuel_toolchain_toml;
use crate::toolchain::{DistToolchainDescription, Toolchain};

#[derive(Debug, Deserialize, Serialize)]
pub struct ToolchainOverride {
    pub cfg: OverrideCfg,
    pub path: PathBuf,
}

#[derive(Debug, Deserialize, Serialize)]
pub struct OverrideCfg {
    pub toolchain: ToolchainCfg,
    pub components: Option<HashMap<String, Version>>,
}

#[derive(Debug, Deserialize)]
pub struct ToolchainCfg {
    #[serde(deserialize_with = "deserialize_channel")]
    pub channel: Channel,
}

#[derive(Debug, PartialEq, Eq)]
pub struct Channel {
    pub name: String,
    pub date: Option<Date>,
}

impl Serialize for ToolchainCfg {
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;

        let mut state = serializer.serialize_struct("ToolchainCfg", 1)?;
        state.serialize_field("channel", &self.channel.to_string())?;
        state.end()
    }
}

fn deserialize_channel<'de, D>(deserializer: D) -> Result<Channel, D::Error>
where
    D: Deserializer<'de>,
{
    let channel = String::deserialize(deserializer)?;
    channel.parse().map_err(|_| {
        D::Error::invalid_value(
            serde::de::Unexpected::Str(&channel),
            &"one of beta-1, beta-2, latest-YYYY-MM-DD, or nightly-YYYY-MM-DD",
        )
    })
}

impl fmt::Display for Channel {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self.date {
            Some(date) => write!(f, "{}-{}", self.name, date),
            None => write!(f, "{}", self.name),
        }
    }
}

impl FromStr for Channel {
    type Err = anyhow::Error;

    fn from_str(value: &str) -> Result<Self> {
        if is_beta_toolchain(value) {
            return Ok(Self {
                name: value.to_string(),
                date: None,
            });
        }

        match value.split_once('-') {
            Some((LATEST | NIGHTLY, date)) => Ok(Self {
                name: value[..value.len() - date.len() - 1].to_string(),
                date: Some(Date::parse(date, DATE_FORMAT)?),
            }),
            _ => bail!("invalid toolchain channel '{value}'"),
        }
    }
}

impl ToolchainOverride {
    pub(crate) fn from_path(path: PathBuf) -> Result<Self> {
        let contents = file::read_file(FUEL_TOOLCHAIN_TOML_FILE, &path)?;
        let cfg = OverrideCfg::from_toml(&contents)?;
        Ok(Self { cfg, path })
    }

    pub fn from_project_root() -> Option<Self> {
        let path = get_fuel_toolchain_toml()?;
        match Self::from_path(path) {
            Ok(toolchain_override) => Some(toolchain_override),
            Err(error) => {
                warn!("warning: invalid '{FUEL_TOOLCHAIN_TOML_FILE}' in project root: {error}");
                None
            }
        }
    }

    pub fn toolchain(&self) -> Result<Toolchain> {
        let description = self.description()?;
        let toolchain = Toolchain::from_path(&description.to_string());
        toolchain.install_if_nonexistent(&description)?;
        Ok(toolchain)
    }

    pub fn description(&self) -> Result<DistToolchainDescription> {
        DistToolchainDescription::from_str(&self.cfg.toolchain.channel.to_string())
    }

    pub fn to_toml(&self) -> Document {
        let mut document = Document::new();
        document["toolchain"]["channel"] = value(self.cfg.toolchain.channel.to_string());

        if let Some(components) = &self.cfg.components {
            for (name, version) in components {
                document["components"][name] = value(version.to_string());
            }
        }

        document
    }
}

impl OverrideCfg {
    pub fn new(toolchain: ToolchainCfg, components: Option<HashMap<String, Version>>) -> Self {
        Self {
            toolchain,
            components,
        }
    }

    pub(crate) fn from_toml(toml: &str) -> Result<Self> {
        let cfg: Self = de::from_str(toml)?;
        if cfg.description().is_err() {
            bail!("invalid channel '{}'", cfg.toolchain.channel);
        }

        if matches!(&cfg.components, Some(components) if components.is_empty()) {
            bail!("'[components]' table is declared with no components");
        }

        Ok(cfg)
    }

    #[cfg(test)]
    pub(crate) fn to_string_pretty(&self) -> Result<String, toml_edit::ser::Error> {
        toml_edit::ser::to_string_pretty(self)
    }

    fn description(&self) -> Result<DistToolchainDescription> {
        DistToolchainDescription::from_str(&self.toolchain.channel.to_string())
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::channel::{BETA_1, BETA_2};

    #[test]
    fn parses_beta_channel() {
        let toml = r#"[toolchain]
channel = "beta-1"
"#;

        let cfg = OverrideCfg::from_toml(toml).unwrap();
        assert_eq!(cfg.toolchain.channel.to_string(), BETA_1);
        assert_eq!(cfg.to_string_pretty().unwrap(), toml);
    }

    #[test]
    fn parses_dated_channel() {
        let toml = r#"[toolchain]
channel = "nightly-2023-01-09"
"#;

        let cfg = OverrideCfg::from_toml(toml).unwrap();
        assert_eq!(cfg.toolchain.channel.to_string(), "nightly-2023-01-09");
    }

    #[test]
    fn rejects_channels_without_pinned_date() {
        assert!(OverrideCfg::from_toml(
            r#"[toolchain]
channel = "latest"
"#
        )
        .is_err());
        assert!(OverrideCfg::from_toml(
            r#"[toolchain]
channel = "nightly"
"#
        )
        .is_err());
    }

    #[test]
    fn rejects_invalid_override_toml() {
        assert!(OverrideCfg::from_toml("").is_err());
        assert!(OverrideCfg::from_toml("[toolchain]\n").is_err());
        assert!(OverrideCfg::from_toml(
            r#"[toolchain]
channel = "invalid"
"#
        )
        .is_err());
        assert!(OverrideCfg::from_toml(
            r#"[toolchain]
channel = "beta-2"

[components]
"#
        )
        .is_err());
    }

    #[test]
    fn parses_supported_channels() {
        assert!(Channel::from_str(BETA_1).is_ok());
        assert!(Channel::from_str(BETA_2).is_ok());
        assert!(Channel::from_str("latest-2023-01-09").is_ok());
        assert!(Channel::from_str("nightly-2023-01-09").is_ok());
        assert!(Channel::from_str(LATEST).is_err());
        assert!(Channel::from_str(NIGHTLY).is_err());
    }
}
