import math
from typing import Any
from xml.etree.ElementTree import QName


def update(obj: Any, **kwargs: Any):
    """Update an object from keyword arguments with dotted keys."""

    for key, value in kwargs.items():
        attrsetter(obj, key, value)


def attrsetter(obj: Any, attr: str, value: Any):
    names = attr.split(".")
    last = names.pop()
    for name in names:
        obj = getattr(obj, name)

    setattr(obj, last, value)


def literal_value(value: Any) -> str:
    """Return a python literal expression for the given value."""
    if isinstance(value, float) and not math.isfinite(value):
        return f'float("{value}")'

    if isinstance(value, QName):
        return f'QName("{value.text}")'

    return repr(value)
