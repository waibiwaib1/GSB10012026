import json

from moto.core.responses import BaseResponse

from .models import FSxBackend, fsx_backends


class FSxResponse(BaseResponse):
    def __init__(self) -> None:
        super().__init__(service_name="fsx")

    @property
    def fsx_backend(self) -> FSxBackend:
        return fsx_backends[self.current_account][self.region]

    def create_file_system(self) -> str:
        file_system = self.fsx_backend.create_file_system(
            file_system_type=self._get_param("FileSystemType"),
            storage_capacity=self._get_param("StorageCapacity"),
            subnet_ids=self._get_param("SubnetIds"),
            storage_type=self._get_param("StorageType"),
            security_group_ids=self._get_param("SecurityGroupIds"),
            kms_key_id=self._get_param("KmsKeyId"),
            tags=self._get_param("Tags"),
            lustre_configuration=self._get_param("LustreConfiguration"),
            ontap_configuration=self._get_param("OntapConfiguration"),
            open_zfs_configuration=self._get_param("OpenZFSConfiguration"),
            windows_configuration=self._get_param("WindowsConfiguration"),
        )
        return json.dumps({"FileSystem": file_system.to_dict()})

    def describe_file_systems(self) -> str:
        file_system_ids = self._get_param("FileSystemIds")
        file_systems = self.fsx_backend.describe_file_systems(file_system_ids)
        return json.dumps(
            {"FileSystems": [file_system.to_dict() for file_system in file_systems]}
        )

    def delete_file_system(self) -> str:
        file_system_id = self._get_param("FileSystemId")
        file_system = self.fsx_backend.delete_file_system(file_system_id)
        return json.dumps(
            {
                "FileSystemId": file_system.file_system_id,
                "Lifecycle": file_system.lifecycle,
            }
        )

    def tag_resource(self) -> str:
        resource_arn = self._get_param("ResourceARN")
        tags = self._get_param("Tags")
        self.fsx_backend.tag_resource(resource_arn, tags)
        return json.dumps({})

    def untag_resource(self) -> str:
        resource_arn = self._get_param("ResourceARN")
        tag_keys = self._get_param("TagKeys")
        self.fsx_backend.untag_resource(resource_arn, tag_keys)
        return json.dumps({})
