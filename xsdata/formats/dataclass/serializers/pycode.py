from dataclasses import dataclass
from dataclasses import field
from enum import Enum
from io import StringIO
from typing import Any
from typing import Iterator
from typing import List
from typing import Mapping
from typing import Set
from typing import TextIO
from typing import Tuple
from typing import Union

from xsdata.formats.bindings import AbstractSerializer
from xsdata.formats.dataclass.context import XmlContext
from xsdata.formats.dataclass.serializers.config import SerializerConfig
from xsdata.utils import collections
from xsdata.utils.objects import literal_value

SPACES = "    "
UNSET = object()


@dataclass
class PycodeSerializer(AbstractSerializer):
    """Serialize binding instances to python representation code."""

    config: SerializerConfig = field(default_factory=SerializerConfig)
    context: XmlContext = field(default_factory=XmlContext)

    def render(self, obj: Any, var_name: str = "obj") -> str:
        """Convert the given object tree to executable python code."""
        output = StringIO()
        self.write(output, obj, var_name)
        return output.getvalue()

    def write(self, out: TextIO, obj: Any, var_name: str = "obj"):
        """Write the given object tree to the output text stream."""
        types: Set[type] = set()
        body = StringIO()

        for chunk in self.repr_object(obj, 0, types):
            body.write(chunk)

        out.write(self.build_imports(types))
        out.write("\n\n")
        out.write(f"{var_name} = ")
        out.write(body.getvalue())
        out.write("\n")

    @classmethod
    def build_imports(cls, types: Set[type]) -> str:
        """Build import statements for all non-builtin referenced types."""
        imports = set()

        for tp in types:
            module = tp.__module__
            if module == "builtins":
                continue

            name = tp.__qualname__.split(".")[0]
            imports.add(f"from {module} import {name}\n")

        return "".join(sorted(imports))

    def repr_object(
        self, obj: Any, level: int, types: Set[type]
    ) -> Iterator[str]:
        """Yield the python representation of the given object."""
        types.add(type(obj))

        if collections.is_array(obj):
            yield from self.repr_array(obj, level, types)
        elif isinstance(obj, dict):
            yield from self.repr_mapping(obj, level, types)
        elif self.context.class_type.is_model(obj):
            yield from self.repr_model(obj, level, types)
        elif isinstance(obj, Enum):
            yield str(obj)
        else:
            yield literal_value(obj)

    def repr_array(
        self,
        obj: Union[List, Tuple],
        level: int,
        types: Set[type],
    ) -> Iterator[str]:
        """Yield the python representation of an array."""
        if not obj:
            yield str(obj)
            return

        next_level = level + 1
        yield "[\n"
        for value in obj:
            yield SPACES * next_level
            yield from self.repr_object(value, next_level, types)
            yield ",\n"

        yield f"{SPACES * level}]"

    def repr_mapping(
        self, obj: Mapping, level: int, types: Set[type]
    ) -> Iterator[str]:
        """Yield the python representation of a mapping."""
        if not obj:
            yield str(obj)
            return

        next_level = level + 1
        yield "{\n"
        for key, value in obj.items():
            yield SPACES * next_level
            yield from self.repr_object(key, next_level, types)
            yield ": "
            yield from self.repr_object(value, next_level, types)
            yield ",\n"

        yield f"{SPACES * level}}}"

    def repr_model(self, obj: Any, level: int, types: Set[type]) -> Iterator[str]:
        """Yield a model constructor expression."""
        yield f"{obj.__class__.__qualname__}(\n"

        next_level = level + 1
        index = 0
        for model_field in self.context.class_type.get_fields(obj):
            if not model_field.init:
                continue

            value = getattr(obj, model_field.name)
            default = self.default_value(model_field)
            if self.is_default(value, default):
                continue

            if index:
                yield f",\n{SPACES * next_level}{model_field.name}="
            else:
                yield f"{SPACES * next_level}{model_field.name}="

            yield from self.repr_object(value, next_level, types)
            index += 1

        yield f"\n{SPACES * level})"

    def default_value(self, model_field: Any) -> Any:
        try:
            return self.context.class_type.default_value(
                model_field, default=UNSET
            )
        except TypeError:
            return self.context.class_type.default_value(model_field)

    @staticmethod
    def is_default(value: Any, default: Any) -> bool:
        if default is UNSET:
            return False

        if callable(default):
            try:
                return default() == value
            except Exception:
                return False

        return default == value
