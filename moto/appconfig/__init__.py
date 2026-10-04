from ..core.models import base_decorator
from .models import appconfig_backends


mock_appconfig = base_decorator(appconfig_backends)
