# Licensed to Modin Development Team under one or more contributor license agreements.
# See the NOTICE file distributed with this work for additional information regarding
# copyright ownership.  The Modin Development Team licenses this file to you under the
# Apache License, Version 2.0 (the "License"); you may not use this file except in
# compliance with the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software distributed under
# the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF
# ANY KIND, either express or implied. See the License for the specific language
# governing permissions and limitations under the License.

import pytest

from modin.core.storage_formats.base.query_compiler import BaseQueryCompiler
from modin.core.storage_formats.pandas.query_compiler_caster import (
    cast_nested_args_to_current_qc_type,
)


class FakeFrame:
    pass


class FakeQueryCompiler(BaseQueryCompiler):
    storage_format = "fake"
    engine = "fake"

    def __init__(self):
        self._modin_frame = FakeFrame()

    def to_pandas(self):
        return "pandas-data"

    @classmethod
    def from_pandas(cls, df, data_cls):
        qc = cls()
        qc.cast_from = df
        return qc

    @classmethod
    def from_arrow(cls, at, data_cls):
        return cls()

    @classmethod
    def from_interchange_dataframe(cls, df, data_cls):
        return cls()

    def to_interchange_dataframe(self, nan_as_null):
        pass

    def free(self):
        pass

    def finalize(self):
        pass

    def execute(self):
        return self


class CurrentQC(FakeQueryCompiler):
    pass


class ArgumentQC(FakeQueryCompiler):
    pass


class BlockingQC(FakeQueryCompiler):
    def allow_coercion_to(self, other_qc):
        return False


def test_cast_argument_when_coercion_allowed():
    current_qc = CurrentQC()
    casted = cast_nested_args_to_current_qc_type((ArgumentQC(),), current_qc)
    assert isinstance(casted[0], CurrentQC)
    assert casted[0].cast_from == "pandas-data"


def test_cast_argument_blocked_by_allow_coercion_to():
    current_qc = CurrentQC()
    with pytest.raises(TypeError, match="allow_coercion_to"):
        cast_nested_args_to_current_qc_type((BlockingQC(),), current_qc)
