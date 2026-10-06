from moto.core.exceptions import JsonRESTError


class FSxClientError(JsonRESTError):
    code = 400


class FileSystemNotFound(FSxClientError):
    def __init__(self, file_system_id: str):
        super().__init__(
            "FileSystemNotFound",
            f"File system '{file_system_id}' does not exist.",
        )


class ResourceNotFound(FSxClientError):
    def __init__(self, resource_arn: str):
        super().__init__(
            "ResourceNotFound",
            f"Resource '{resource_arn}' does not exist.",
        )
