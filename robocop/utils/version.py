import re
from functools import total_ordering


VERSION_RE = re.compile(
    r"v?\s*(?:(\d+)!)?"
    r"(?:(\d+)(?:\.(\d+))?(?:\.(\d+))?)?"
    r"(?:(?:[-_\.]?(a|alpha|b|beta|c|rc|preview))[-_\.]?(\d*))?"
    r"(?:(?:[-_\.]?(dev))[-_\.]?(\d*))?"
    r"(?:\+([a-z0-9]+(?:[-_\.][a-z0-9]+)*))?",
    re.IGNORECASE,
)
SPECIFIER_RE = re.compile(r"\s*(===|==|!=|~=|>=|<=|>|<)\s*([^,\s]+)\s*", re.IGNORECASE)
PRE_RELEASE_VALUES = {"a": 0, "alpha": 0, "b": 1, "beta": 1, "c": 2, "rc": 2, "preview": 2}


def _number(value, default=0):
    return int(value) if value else default


@total_ordering
class Version:
    def __init__(self, value):
        self.raw = str(value).strip()
        match = VERSION_RE.fullmatch(self.raw)
        if not match:
            raise ValueError(f"Invalid version: '{value}'")
        epoch, major, minor, patch, pre_release, pre_number, dev, dev_number, _local = match.groups()
        if major is None:
            raise ValueError(f"Invalid version: '{value}'")
        self.epoch = _number(epoch)
        self.major = _number(major)
        self.minor = _number(minor)
        self.patch = _number(patch)
        self.release_parts = tuple(
            _number(part) for part in (major, minor, patch) if part is not None
        )
        if dev:
            self.pre_release = (-1, _number(dev_number))
        elif pre_release:
            self.pre_release = (0, PRE_RELEASE_VALUES[pre_release.lower()], _number(pre_number))
        else:
            self.pre_release = None

    @property
    def release(self):
        return self.epoch, self.major, self.minor, self.patch

    @property
    def _sort_key(self):
        pre_release = self.pre_release if self.pre_release is not None else (1,)
        return self.epoch, self.major, self.minor, self.patch, pre_release

    def __eq__(self, other):
        return self._sort_key == self._coerce(other)._sort_key

    def __lt__(self, other):
        return self._sort_key < self._coerce(other)._sort_key

    def __hash__(self):
        return hash(self._sort_key)

    def __str__(self):
        return self.raw

    @staticmethod
    def _coerce(other):
        return other if isinstance(other, Version) else Version(other)


def parse_version(value):
    return Version(value)


def _parse_specifier_part(part):
    match = SPECIFIER_RE.fullmatch(part)
    if not match:
        raise ValueError(f"Invalid version specifier: '{part}'")
    return match.group(1).lower(), match.group(2)


def _padded_release_parts(version_object, length):
    parts = version_object.release_parts
    return parts + (0,) * max(0, length - len(parts))


def _matches_specifier(version_object, operator, expected):
    if operator == "===":
        return version_object.raw == expected

    wildcard = expected.endswith(".*")
    expected_value = expected[:-2] if wildcard else expected
    expected_version = Version(expected_value)
    parts = expected_version.release_parts
    actual = _padded_release_parts(version_object, len(parts))
    if expected_version.epoch:
        raise ValueError("Epoch versions are not supported in version specifiers")

    if wildcard:
        prefix_matches = actual[: len(parts)] == parts
        if version_object.pre_release is not None:
            prefix_matches = prefix_matches and expected_version.pre_release is not None
        if operator == "==":
            return prefix_matches
        if operator == "!=":
            return not prefix_matches
        raise ValueError(f"Invalid version specifier: '{operator}{expected}'")

    if operator == "~=":
        if len(parts) < 2:
            raise ValueError(f"Invalid version specifier: '~={expected}'")
        lower_bound = parts
        upper_bound = parts[:-2] + (parts[-2] + 1,)
        actual_release = actual[: len(lower_bound)]
        return actual_release >= lower_bound and actual[: len(upper_bound)] < upper_bound

    return {
        "==": version_object == expected_version,
        "!=": version_object != expected_version,
        ">=": version_object >= expected_version,
        "<=": version_object <= expected_version,
        ">": version_object > expected_version,
        "<": version_object < expected_version,
    }[operator]


def matches_specifiers(value, specifiers):
    version_object = value if isinstance(value, Version) else Version(value)
    return all(
        _matches_specifier(version_object, *_parse_specifier_part(part))
        for part in str(specifiers).split(",")
        if part.strip()
    )
