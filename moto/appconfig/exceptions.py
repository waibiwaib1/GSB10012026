from moto.core.exceptions import JsonRESTError


class ResourceNotFound(JsonRESTError):
    def __init__(self, message: str):
        super().__init__("ResourceNotFoundException", message)
        self.code = 404
