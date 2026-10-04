from typing import Any, Dict, List, Optional

from moto.core import BackendDict, BaseBackend, BaseModel
from moto.moto_api._internal import mock_random

from .exceptions import ResourceNotFound


class Application(BaseModel):
    def __init__(self, name: str, description: Optional[str]):
        self.id = mock_random.get_random_hex(7)
        self.name = name
        self.description = description
        self.configuration_profiles: Dict[str, ConfigurationProfile] = {}

    def to_dict(self) -> Dict[str, Any]:
        return {
            "Id": self.id,
            "Name": self.name,
            "Description": self.description,
        }


class ConfigurationProfile(BaseModel):
    def __init__(
        self,
        application_id: str,
        name: str,
        description: Optional[str],
        location_uri: str,
        retrieval_role_arn: Optional[str],
        validators: Optional[List[Dict[str, str]]],
        profile_type: Optional[str],
        kms_key_identifier: Optional[str],
    ):
        self.application_id = application_id
        self.id = mock_random.get_random_hex(7)
        self.name = name
        self.description = description
        self.location_uri = location_uri
        self.retrieval_role_arn = retrieval_role_arn
        self.validators = validators or []
        self.type = profile_type
        self.kms_key_identifier = kms_key_identifier
        self.versions: Dict[int, HostedConfigurationVersion] = {}

    def to_dict(self) -> Dict[str, Any]:
        return {
            "ApplicationId": self.application_id,
            "Id": self.id,
            "Name": self.name,
            "Description": self.description,
            "LocationUri": self.location_uri,
            "RetrievalRoleArn": self.retrieval_role_arn,
            "Validators": self.validators,
            "Type": self.type,
            "KmsKeyIdentifier": self.kms_key_identifier,
        }


class HostedConfigurationVersion(BaseModel):
    def __init__(
        self,
        application_id: str,
        configuration_profile_id: str,
        version_number: int,
        description: Optional[str],
        content: bytes,
        content_type: str,
        version_label: Optional[str],
    ):
        self.application_id = application_id
        self.configuration_profile_id = configuration_profile_id
        self.version_number = version_number
        self.description = description
        self.content = content
        self.content_type = content_type
        self.version_label = version_label

    def headers(self) -> Dict[str, str]:
        headers = {
            "Application-Id": self.application_id,
            "Configuration-Profile-Id": self.configuration_profile_id,
            "Version-Number": str(self.version_number),
            "Content-Type": self.content_type,
        }
        if self.description is not None:
            headers["Description"] = self.description
        if self.version_label is not None:
            headers["VersionLabel"] = self.version_label
        return headers


class AppConfigBackend(BaseBackend):
    def __init__(self, region_name: str, account_id: str):
        super().__init__(region_name, account_id)
        self.applications: Dict[str, Application] = {}

    def create_application(self, name: str, description: Optional[str]) -> Application:
        application = Application(name=name, description=description)
        self.applications[application.id] = application
        return application

    def create_configuration_profile(
        self,
        application_id: str,
        name: str,
        description: Optional[str],
        location_uri: str,
        retrieval_role_arn: Optional[str],
        validators: Optional[List[Dict[str, str]]],
        profile_type: Optional[str],
        kms_key_identifier: Optional[str],
    ) -> ConfigurationProfile:
        application = self._get_application(application_id)
        profile = ConfigurationProfile(
            application_id=application_id,
            name=name,
            description=description,
            location_uri=location_uri,
            retrieval_role_arn=retrieval_role_arn,
            validators=validators,
            profile_type=profile_type,
            kms_key_identifier=kms_key_identifier,
        )
        application.configuration_profiles[profile.id] = profile
        return profile

    def create_hosted_configuration_version(
        self,
        application_id: str,
        configuration_profile_id: str,
        description: Optional[str],
        content: bytes,
        content_type: str,
        version_label: Optional[str],
    ) -> HostedConfigurationVersion:
        profile = self._get_configuration_profile(
            application_id, configuration_profile_id
        )
        version_number = len(profile.versions) + 1
        version = HostedConfigurationVersion(
            application_id=application_id,
            configuration_profile_id=configuration_profile_id,
            version_number=version_number,
            description=description,
            content=content,
            content_type=content_type,
            version_label=version_label,
        )
        profile.versions[version_number] = version
        return version

    def get_hosted_configuration_version(
        self,
        application_id: str,
        configuration_profile_id: str,
        version_number: int,
    ) -> HostedConfigurationVersion:
        profile = self._get_configuration_profile(
            application_id, configuration_profile_id
        )
        try:
            return profile.versions[version_number]
        except KeyError:
            raise ResourceNotFound(
                f"Requested configuration version {version_number} not found"
            )

    def _get_application(self, application_id: str) -> Application:
        try:
            return self.applications[application_id]
        except KeyError:
            raise ResourceNotFound(f"Requested application {application_id} not found")

    def _get_configuration_profile(
        self, application_id: str, configuration_profile_id: str
    ) -> ConfigurationProfile:
        application = self._get_application(application_id)
        try:
            return application.configuration_profiles[configuration_profile_id]
        except KeyError:
            raise ResourceNotFound(
                f"Requested configuration profile {configuration_profile_id} not found"
            )


appconfig_backends = BackendDict(AppConfigBackend, "appconfig")
