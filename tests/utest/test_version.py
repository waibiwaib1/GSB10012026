import pytest

from robocop.utils.version import Version, matches_specifiers, parse_version


@pytest.mark.parametrize(
    "value,major,minor",
    [
        ("3", 3, 0),
        ("4.1", 4, 1),
        ("5.0.1", 5, 0),
        ("v6.0", 6, 0),
    ],
)
def test_parse_version(value, major, minor):
    version = parse_version(value)
    assert version.major == major
    assert version.minor == minor


@pytest.mark.parametrize(
    "left,right",
    [
        ("4.0", "5.0"),
        ("5.0", "5.0.1"),
        ("6.0.dev1", "6.0a1"),
        ("6.0a1", "6.0rc1"),
        ("6.0rc1", "6.0"),
    ],
)
def test_version_ordering(left, right):
    assert Version(left) < Version(right)


@pytest.mark.parametrize(
    "value,specifier,expected",
    [
        ("4.0", ">=4.0", True),
        ("4.0", ">=5.0", False),
        ("5.1", "<5.0", False),
        ("4.1.3", "==4.*", True),
        ("4.0.0rc1", "==4.*", False),
        ("4.0rc1", "==4.0rc1", True),
        ("5.0", "==4.*", False),
        ("5.0", "!=4.*", True),
        ("4.9", "~=4.0", True),
        ("5.0", "~=4.0", False),
        ("5.0", ">=5.0,<6.0", True),
        ("6.0", ">=5.0,<6.0", False),
    ],
)
def test_matches_specifiers(value, specifier, expected):
    assert matches_specifiers(value, specifier) is expected


@pytest.mark.parametrize("value", ["", "not-a-version"])
def test_invalid_version(value):
    with pytest.raises(ValueError):
        Version(value)


@pytest.mark.parametrize("specifier", [">=not-a-version", "=4.0"])
def test_invalid_specifier(specifier):
    with pytest.raises(ValueError):
        matches_specifiers("4.0", specifier)
