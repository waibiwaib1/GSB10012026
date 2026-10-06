import boto3
import pytest
from botocore.exceptions import ClientError

from moto import mock_aws


@mock_aws
def test_create_file_system():
    client = boto3.client("fsx", region_name="us-east-1")
    resp = client.create_file_system(
        FileSystemType="LUSTRE",
        StorageCapacity=1200,
        SubnetIds=["subnet-12345678"],
        Tags=[{"Key": "Name", "Value": "test-fs"}],
    )
    file_system = resp["FileSystem"]
    assert file_system["FileSystemId"].startswith("fs-")
    assert file_system["FileSystemType"] == "LUSTRE"
    assert file_system["StorageCapacity"] == 1200
    assert file_system["StorageType"] == "SSD"
    assert file_system["Lifecycle"] == "AVAILABLE"
    assert file_system["SubnetIds"] == ["subnet-12345678"]
    assert file_system["Tags"] == [{"Key": "Name", "Value": "test-fs"}]
    assert file_system["ResourceARN"].endswith(
        f"file-system/{file_system['FileSystemId']}"
    )


@mock_aws
def test_describe_file_systems():
    client = boto3.client("fsx", region_name="us-east-1")
    fs1 = client.create_file_system(
        FileSystemType="WINDOWS", StorageCapacity=300, SubnetIds=["subnet-12345678"]
    )["FileSystem"]
    fs2 = client.create_file_system(
        FileSystemType="LUSTRE", StorageCapacity=1200, SubnetIds=["subnet-87654321"]
    )["FileSystem"]

    resp = client.describe_file_systems()
    assert len(resp["FileSystems"]) == 2

    resp = client.describe_file_systems(FileSystemIds=[fs2["FileSystemId"]])
    assert len(resp["FileSystems"]) == 1
    assert resp["FileSystems"][0]["FileSystemId"] == fs2["FileSystemId"]
    assert resp["FileSystems"][0]["FileSystemType"] == "LUSTRE"

    with pytest.raises(ClientError) as exc:
        client.describe_file_systems(FileSystemIds=["fs-00000000000000000"])
    assert exc.value.response["Error"]["Code"] == "FileSystemNotFound"

    assert fs1["FileSystemId"] != fs2["FileSystemId"]


@mock_aws
def test_delete_file_system():
    client = boto3.client("fsx", region_name="us-east-1")
    file_system_id = client.create_file_system(
        FileSystemType="LUSTRE", StorageCapacity=1200, SubnetIds=["subnet-12345678"]
    )["FileSystem"]["FileSystemId"]

    resp = client.delete_file_system(FileSystemId=file_system_id)
    assert resp["FileSystemId"] == file_system_id
    assert resp["Lifecycle"] == "DELETING"

    resp = client.describe_file_systems()
    assert resp["FileSystems"] == []

    with pytest.raises(ClientError) as exc:
        client.delete_file_system(FileSystemId=file_system_id)
    assert exc.value.response["Error"]["Code"] == "FileSystemNotFound"


@mock_aws
def test_tag_and_untag_resource():
    client = boto3.client("fsx", region_name="us-east-1")
    file_system = client.create_file_system(
        FileSystemType="LUSTRE",
        StorageCapacity=1200,
        SubnetIds=["subnet-12345678"],
        Tags=[{"Key": "Name", "Value": "test-fs"}],
    )["FileSystem"]
    arn = file_system["ResourceARN"]

    client.tag_resource(
        ResourceARN=arn,
        Tags=[{"Key": "Environment", "Value": "test"}],
    )
    tags = client.describe_file_systems(
        FileSystemIds=[file_system["FileSystemId"]]
    )["FileSystems"][0]["Tags"]
    assert {"Key": "Name", "Value": "test-fs"} in tags
    assert {"Key": "Environment", "Value": "test"} in tags

    client.untag_resource(ResourceARN=arn, TagKeys=["Name"])
    tags = client.describe_file_systems(
        FileSystemIds=[file_system["FileSystemId"]]
    )["FileSystems"][0]["Tags"]
    assert tags == [{"Key": "Environment", "Value": "test"}]

    with pytest.raises(ClientError) as exc:
        client.tag_resource(
            ResourceARN="arn:aws:fsx:us-east-1:123456789012:file-system/fs-unknown",
            Tags=[{"Key": "k", "Value": "v"}],
        )
    assert exc.value.response["Error"]["Code"] == "ResourceNotFound"
