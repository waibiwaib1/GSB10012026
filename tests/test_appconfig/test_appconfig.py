import json

import boto3
import pytest
from botocore.exceptions import ClientError

from moto import mock_appconfig


NORMAL_CONTENTS = {"key": "value"}


@mock_appconfig
def test_create_and_get_hosted_configuration_version():
    client = boto3.client("appconfig", region_name="eu-west-1")

    application = client.create_application(Name="Test")
    configuration_profile = client.create_configuration_profile(
        ApplicationId=application["Id"],
        Name="Normal",
        LocationUri="hosted",
        Type="freeform",
    )
    created_version = client.create_hosted_configuration_version(
        ApplicationId=application["Id"],
        ConfigurationProfileId=configuration_profile["Id"],
        Content=json.dumps(NORMAL_CONTENTS).encode("utf-8"),
        ContentType="application/json",
    )

    assert created_version["ApplicationId"] == application["Id"]
    assert created_version["ConfigurationProfileId"] == configuration_profile["Id"]
    assert created_version["VersionNumber"] == 1
    assert created_version["ContentType"] == "application/json"
    assert json.loads(created_version["Content"].read()) == NORMAL_CONTENTS

    fetched_version = client.get_hosted_configuration_version(
        ApplicationId=application["Id"],
        ConfigurationProfileId=configuration_profile["Id"],
        VersionNumber=1,
    )

    assert fetched_version["ApplicationId"] == application["Id"]
    assert fetched_version["ConfigurationProfileId"] == configuration_profile["Id"]
    assert fetched_version["VersionNumber"] == 1
    assert fetched_version["ContentType"] == "application/json"
    assert json.loads(fetched_version["Content"].read()) == NORMAL_CONTENTS


@mock_appconfig
def test_hosted_configuration_versions_increment():
    client = boto3.client("appconfig", region_name="eu-west-1")
    application = client.create_application(Name="Test")
    configuration_profile = client.create_configuration_profile(
        ApplicationId=application["Id"],
        Name="Normal",
        LocationUri="hosted",
        Type="freeform",
    )

    first_version = client.create_hosted_configuration_version(
        ApplicationId=application["Id"],
        ConfigurationProfileId=configuration_profile["Id"],
        Content=b"first",
        ContentType="text/plain",
    )
    second_version = client.create_hosted_configuration_version(
        ApplicationId=application["Id"],
        ConfigurationProfileId=configuration_profile["Id"],
        Content=b"second",
        ContentType="text/plain",
    )

    assert first_version["VersionNumber"] == 1
    assert second_version["VersionNumber"] == 2
    assert second_version["Content"].read() == b"second"


@mock_appconfig
def test_get_missing_hosted_configuration_version():
    client = boto3.client("appconfig", region_name="eu-west-1")
    application = client.create_application(Name="Test")
    configuration_profile = client.create_configuration_profile(
        ApplicationId=application["Id"],
        Name="Normal",
        LocationUri="hosted",
        Type="freeform",
    )

    with pytest.raises(ClientError) as exc_info:
        client.get_hosted_configuration_version(
            ApplicationId=application["Id"],
            ConfigurationProfileId=configuration_profile["Id"],
            VersionNumber=1,
        )

    err = exc_info.value.response["Error"]
    assert err["Code"] == "ResourceNotFoundException"
