#!/usr/bin/env python3

# (C) Copyright 2020 ECMWF.
#
# This software is licensed under the terms of the Apache Licence Version 2.0
# which can be obtained at http://www.apache.org/licenses/LICENSE-2.0.
# In applying this licence, ECMWF does not waive the privileges and immunities
# granted to it by virtue of its status as an intergovernmental organisation
# nor does it submit to any jurisdiction.
#

import pytest

from earthkit.data import from_source
from earthkit.data.core.fieldlist import FieldList
from earthkit.data.testing import earthkit_examples_file


HIDDEN_KEYS = [
    "min",
    "max",
    "avg",
    "sd",
    "skew",
    "kurt",
    "const",
    "isConstant",
    "numberOfMissing",
    "numberOfCodedValues",
    "bitmapPresent",
    "offsetValuesBy",
    "packingError",
    "referenceValue",
    "referenceValueError",
    "unpackedError",
]

HIDDEN_ALIASES = [
    "minimum",
    "maximum",
    "average",
    "standardDeviation",
    "skewness",
    "kurtosis",
]


def test_numpy_fs_hides_value_related_metadata():
    ds = from_source("file", earthkit_examples_file("test.grib"))
    r = FieldList.from_numpy(
        ds[0].values + 1, ds[0].metadata().override(shortName="pt")
    )

    md = r[0].metadata()
    for key in HIDDEN_KEYS + HIDDEN_ALIASES:
        with pytest.raises(KeyError):
            r[0].metadata(key)
        with pytest.raises(KeyError):
            md[key]
        assert r[0].metadata(key, default=None) is None
        assert md.get(key) is None
        assert key not in md

    with pytest.raises(KeyError):
        r[0].metadata("statistics.max")


def test_numpy_fs_hides_statistics_namespace():
    ds = from_source("file", earthkit_examples_file("test.grib"))
    r = FieldList.from_numpy(
        ds[0].values + 1, ds[0].metadata().override(shortName="pt")
    )
    md = r[0].metadata()

    assert r[0].metadata(namespace="statistics") == {}
    assert "statistics" not in r[0].metadata(namespace=all)
    assert "statistics" not in md.namespaces()

    default_namespace = r[0].metadata(namespace=None)
    for key in HIDDEN_KEYS + HIDDEN_ALIASES:
        assert key not in default_namespace

    dump = r[0].dump(_as_raw=True)
    assert "statistics" not in [namespace["title"] for namespace in dump]

    geography = r[0].metadata(namespace="geography")
    assert "bitmapPresent" not in geography


if __name__ == "__main__":
    from earthkit.data.testing import main

    main(__file__)
