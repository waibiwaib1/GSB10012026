from .responses import AppConfigResponse


url_bases = [
    r"https?://appconfig\.(.+)\.amazonaws\.com",
]


url_paths = {
    "{0}/applications$": AppConfigResponse.dispatch,
    "{0}/applications/(?P<ApplicationId>[^/]+)/configurationprofiles$": AppConfigResponse.dispatch,
    "{0}/applications/(?P<ApplicationId>[^/]+)/configurationprofiles/(?P<ConfigurationProfileId>[^/]+)/hostedconfigurationversions$": (
        AppConfigResponse.dispatch
    ),
    "{0}/applications/(?P<ApplicationId>[^/]+)/configurationprofiles/(?P<ConfigurationProfileId>[^/]+)/hostedconfigurationversions/(?P<VersionNumber>[^/]+)$": (
        AppConfigResponse.dispatch
    ),
}
