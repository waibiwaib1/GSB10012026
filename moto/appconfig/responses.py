import json
from typing import Any, Dict, cast

from moto.core.common_types import TYPE_RESPONSE
from moto.core.responses import BaseResponse

from .models import appconfig_backends, AppConfigBackend


class AppConfigResponse(BaseResponse):
    def __init__(self) -> None:
        super().__init__(service_name="appconfig")

    def _dispatch(self, request: Any, full_url: str, headers: Any) -> TYPE_RESPONSE:
        self.setup_class(request, full_url, headers, use_raw_body=True)
        return self.call_action()

    @property
    def appconfig_backend(self) -> AppConfigBackend:
        return appconfig_backends[self.current_account][self.region]

    @property
    def payload(self) -> Dict[str, Any]:
        if not hasattr(self, "_payload"):
            self._payload = json.loads(self.body) if self.body else {}
        return cast(Dict[str, Any], self._payload)

    def _uri_param(self, name: str) -> str:
        assert self.uri_match is not None
        return self.uri_match.groupdict()[name]

    def create_application(self) -> TYPE_RESPONSE:
        application = self.appconfig_backend.create_application(
            name=cast(str, self.payload.get("Name")),
            description=self.payload.get("Description"),
        )
        return 201, {"status": 201}, json.dumps(application.to_dict())

    def create_configuration_profile(self) -> TYPE_RESPONSE:
        profile = self.appconfig_backend.create_configuration_profile(
            application_id=self._uri_param("ApplicationId"),
            name=cast(str, self.payload.get("Name")),
            description=self.payload.get("Description"),
            location_uri=cast(str, self.payload.get("LocationUri")),
            retrieval_role_arn=self.payload.get("RetrievalRoleArn"),
            validators=self.payload.get("Validators"),
            profile_type=self.payload.get("Type"),
            kms_key_identifier=self.payload.get("KmsKeyIdentifier"),
        )
        return 201, {"status": 201}, json.dumps(profile.to_dict())

    def create_hosted_configuration_version(self) -> TYPE_RESPONSE:
        version = self.appconfig_backend.create_hosted_configuration_version(
            application_id=self._uri_param("ApplicationId"),
            configuration_profile_id=self._uri_param("ConfigurationProfileId"),
            description=self.headers.get("Description"),
            content=self.body,
            content_type=self.headers["Content-Type"],
            version_label=self.headers.get("VersionLabel"),
        )
        headers: Dict[str, Any] = {"status": 201}
        headers.update(version.headers())
        return 201, headers, version.content

    def get_hosted_configuration_version(self) -> TYPE_RESPONSE:
        version = self.appconfig_backend.get_hosted_configuration_version(
            application_id=self._uri_param("ApplicationId"),
            configuration_profile_id=self._uri_param("ConfigurationProfileId"),
            version_number=int(self._uri_param("VersionNumber")),
        )
        return 200, version.headers(), version.content
