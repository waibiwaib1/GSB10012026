from typing import Any, Dict, List, Optional

from moto.core.base_backend import BackendDict, BaseBackend
from moto.core.common_models import BaseModel
from moto.core.utils import unix_time
from moto.moto_api._internal import mock_random
from moto.utilities.tagging_service import TaggingService

from .exceptions import FileSystemNotFound, ResourceNotFound


class FileSystem(BaseModel):
    def __init__(
        self,
        account_id: str,
        region_name: str,
        file_system_type: str,
        storage_capacity: int,
        subnet_ids: List[str],
        storage_type: Optional[str],
        security_group_ids: Optional[List[str]],
        kms_key_id: Optional[str],
        tags: Optional[List[Dict[str, str]]],
        lustre_configuration: Optional[Dict[str, Any]],
        ontap_configuration: Optional[Dict[str, Any]],
        open_zfs_configuration: Optional[Dict[str, Any]],
        windows_configuration: Optional[Dict[str, Any]],
    ):
        self.file_system_id = f"fs-{mock_random.get_random_hex(17)}"
        self.owner_id = account_id
        self.region_name = region_name
        self.file_system_type = file_system_type
        self.storage_capacity = storage_capacity
        self.storage_type = storage_type or "SSD"
        self.subnet_ids = subnet_ids
        self.security_group_ids = security_group_ids or []
        self.kms_key_id = kms_key_id
        self.lustre_configuration = lustre_configuration
        self.ontap_configuration = ontap_configuration
        self.open_zfs_configuration = open_zfs_configuration
        self.windows_configuration = windows_configuration
        self.lifecycle = "AVAILABLE"
        self.creation_time = unix_time()
        self.vpc_id = f"vpc-{mock_random.get_random_hex(8)}"
        self.dns_name = (
            f"{self.file_system_id}.fsx.{region_name}.amazonaws.com"
        )
        self.resource_arn = (
            f"arn:aws:fsx:{region_name}:{account_id}"
            f":file-system/{self.file_system_id}"
        )
        self.tags = tags or []

    def to_dict(self) -> Dict[str, Any]:
        dct: Dict[str, Any] = {
            "CreationTime": self.creation_time,
            "DNSName": self.dns_name,
            "FileSystemId": self.file_system_id,
            "FileSystemType": self.file_system_type,
            "KmsKeyId": self.kms_key_id,
            "Lifecycle": self.lifecycle,
            "OwnerId": self.owner_id,
            "ResourceARN": self.resource_arn,
            "StorageCapacity": self.storage_capacity,
            "StorageType": self.storage_type,
            "SubnetIds": self.subnet_ids,
            "Tags": self.tags,
            "VpcId": self.vpc_id,
        }
        if self.lustre_configuration is not None:
            dct["LustreConfiguration"] = self.lustre_configuration
        if self.ontap_configuration is not None:
            dct["OntapConfiguration"] = self.ontap_configuration
        if self.open_zfs_configuration is not None:
            dct["OpenZFSConfiguration"] = self.open_zfs_configuration
        if self.windows_configuration is not None:
            dct["WindowsConfiguration"] = self.windows_configuration
        return dct


class FSxBackend(BaseBackend):
    def __init__(self, region_name: str, account_id: str):
        super().__init__(region_name, account_id)
        self.file_systems: Dict[str, FileSystem] = {}
        self.tagger = TaggingService()

    def create_file_system(
        self,
        file_system_type: str,
        storage_capacity: int,
        subnet_ids: List[str],
        storage_type: Optional[str],
        security_group_ids: Optional[List[str]],
        kms_key_id: Optional[str],
        tags: Optional[List[Dict[str, str]]],
        lustre_configuration: Optional[Dict[str, Any]],
        ontap_configuration: Optional[Dict[str, Any]],
        open_zfs_configuration: Optional[Dict[str, Any]],
        windows_configuration: Optional[Dict[str, Any]],
    ) -> FileSystem:
        file_system = FileSystem(
            account_id=self.account_id,
            region_name=self.region_name,
            file_system_type=file_system_type,
            storage_capacity=storage_capacity,
            subnet_ids=subnet_ids,
            storage_type=storage_type,
            security_group_ids=security_group_ids,
            kms_key_id=kms_key_id,
            tags=tags,
            lustre_configuration=lustre_configuration,
            ontap_configuration=ontap_configuration,
            open_zfs_configuration=open_zfs_configuration,
            windows_configuration=windows_configuration,
        )
        self.file_systems[file_system.file_system_id] = file_system
        self.tagger.tag_resource(file_system.resource_arn, tags)
        return file_system

    def describe_file_systems(
        self, file_system_ids: Optional[List[str]]
    ) -> List[FileSystem]:
        if file_system_ids:
            file_systems = []
            for file_system_id in file_system_ids:
                if file_system_id not in self.file_systems:
                    raise FileSystemNotFound(file_system_id)
                file_systems.append(self.file_systems[file_system_id])
            return file_systems
        return list(self.file_systems.values())

    def delete_file_system(self, file_system_id: str) -> FileSystem:
        if file_system_id not in self.file_systems:
            raise FileSystemNotFound(file_system_id)
        file_system = self.file_systems.pop(file_system_id)
        self.tagger.delete_all_tags_for_resource(file_system.resource_arn)
        file_system.lifecycle = "DELETING"
        return file_system

    def tag_resource(self, resource_arn: str, tags: List[Dict[str, str]]) -> None:
        file_system = self._get_file_system_by_arn(resource_arn)
        self.tagger.tag_resource(resource_arn, tags)
        file_system.tags = self.tagger.list_tags_for_resource(resource_arn)["Tags"]

    def untag_resource(self, resource_arn: str, tag_keys: List[str]) -> None:
        file_system = self._get_file_system_by_arn(resource_arn)
        self.tagger.untag_resource_using_names(resource_arn, tag_keys)
        file_system.tags = self.tagger.list_tags_for_resource(resource_arn)["Tags"]

    def _get_file_system_by_arn(self, resource_arn: str) -> FileSystem:
        for file_system in self.file_systems.values():
            if file_system.resource_arn == resource_arn:
                return file_system
        raise ResourceNotFound(resource_arn)


fsx_backends = BackendDict(FSxBackend, "fsx")
