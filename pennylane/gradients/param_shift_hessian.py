# Copyright 2018-2021 Xanadu Quantum Technologies Inc.

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
"""
This module contains functions for directly computing the parameter-shift
Hessian of a quantum tape.
"""
# pylint: disable=protected-access
from collections import defaultdict

import numpy as np

import pennylane as qml

from .finite_difference import generate_shifted_tapes
from .gradient_transform import gradient_transform
from .general_shift_rules import eigvals_to_frequencies, frequencies_to_period
from .parameter_shift import (
    _get_operation_recipe,
    _gradient_analysis,
    _process_gradient_recipe,
)


def _operation_period(tape, idx):
    """Return the Fourier period of a trainable operation, if available."""

    try:
        op, _ = tape.get_operation(idx)
        generator, coefficient = op.generator
        if isinstance(generator, np.ndarray):
            generator_matrix = generator
        else:
            generator_matrix = generator.matrix
        eigenvalues = tuple(
            np.round(np.linalg.eigvalsh(generator_matrix * coefficient), 12)
        )
        frequencies = eigvals_to_frequencies(eigenvalues)
        return frequencies_to_period(frequencies)
    except Exception:  # pylint: disable=broad-except
        return None


def _wrap_shift(shift, period):
    """Wrap a shift into the interval ``[0, period)``."""

    if period is None:
        return shift
    return np.mod(shift, period)


def _contract_hessian_with_cjac(qhess, cjac):
    """Contract a gate-parameter Hessian with a classical QNode Jacobian."""

    if isinstance(cjac, tuple):
        hessians = tuple(
            _contract_hessian_with_cjac(qhess, cjac_arg)
            for cjac_arg in cjac
            if cjac_arg is not None
        )
        return hessians[0] if len(hessians) == 1 else hessians

    is_square = (
        cjac.shape == (1,)
        or cjac.shape == (1, 1)
        or (cjac.ndim == 2 and cjac.shape[0] == cjac.shape[1])
    )
    if is_square and qml.math.allclose(cjac, qml.numpy.eye(cjac.shape[0])):
        return qhess

    if qml.math.ndim(qhess) == 2:
        contracted = qml.math.tensordot(qhess, cjac, axes=[[0], [0]])
        return qml.math.tensordot(cjac, contracted, axes=[[0], [0]])

    contracted = qml.math.tensordot(qhess, cjac, axes=[[1], [0]])
    contracted = qml.math.tensordot(cjac, contracted, axes=[[0], [1]])
    return qml.math.moveaxis(contracted, 0, 1)


def _generate_multiple_shifted_tapes(tape, shifts):
    """Generate tapes shifted in one or two trainable tape parameters."""

    params = list(tape.get_parameters())
    new_tape = tape.copy(copy_operations=True)

    for idx, multiplier, shift in shifts:
        param = params[idx]
        multiplier = qml.math.convert_like(multiplier, param)
        shift = qml.math.convert_like(shift, param)
        params[idx] = param * multiplier + shift

    new_tape.set_parameters(params)
    return new_tape


def _second_derivative_recipe(tape, idx1, idx2, shift):
    """Return coefficients and two-parameter shifts for a second derivative."""

    recipe1 = _process_gradient_recipe(_get_operation_recipe(tape, idx1, shift=shift))
    recipe2 = _process_gradient_recipe(_get_operation_recipe(tape, idx2, shift=shift))

    coeffs1, multipliers1, shifts1 = recipe1
    coeffs2, multipliers2, shifts2 = recipe2
    period1 = _operation_period(tape, idx1)
    period2 = _operation_period(tape, idx2) if idx2 != idx1 else period1

    terms = []
    for coeff1, multiplier1, param_shift1 in zip(coeffs1, multipliers1, shifts1):
        for coeff2, multiplier2, param_shift2 in zip(coeffs2, multipliers2, shifts2):
            if idx1 == idx2:
                multiplier = multiplier1 * multiplier2
                if np.isclose(multiplier, 0):
                    continue
                combined_shift = _wrap_shift(
                    multiplier1 * param_shift2 + param_shift1, period1
                )
                if np.isclose(combined_shift, 0) and np.isclose(multiplier, 1):
                    multiplier = 1
                    combined_shift = 0
                shift_tuple = ((idx1, multiplier, combined_shift),)
            else:
                shift_tuple = (
                    (idx1, multiplier1, _wrap_shift(param_shift1, period1)),
                    (idx2, multiplier2, _wrap_shift(param_shift2, period2)),
                )

            terms.append((coeff1 * coeff2, shift_tuple))

    combined = defaultdict(float)
    for coeff, shift_tuple in terms:
        key = tuple(
            (idx, float(round(multiplier, 12)), float(round(param_shift, 12)))
            for idx, multiplier, param_shift in shift_tuple
        )
        combined[key] += float(coeff)

    return [(coeff, key) for key, coeff in combined.items() if not np.isclose(coeff, 0)]


def _make_shifted_tape(tape, shift_spec):
    """Create a tape from a normalized shift specification."""

    if len(shift_spec) == 1:
        idx, multiplier, param_shift = shift_spec[0]
        if np.isclose(multiplier, 1):
            return generate_shifted_tapes(tape, idx, [param_shift])[0]

    return _generate_multiple_shifted_tapes(tape, shift_spec)


@gradient_transform
def param_shift_hessian(tape, argnum=None, shift=np.pi / 2, hybrid=True):
    r"""Transform a QNode to compute the parameter-shift Hessian of all gate
    parameters with respect to its inputs.

    Args:
        qnode (pennylane.QNode or .QuantumTape): quantum tape or QNode to
            differentiate
        argnum (int or list[int] or None): Trainable parameter indices to
            differentiate with respect to. If not provided, derivatives with
            respect to all trainable parameters are returned.
        shift (float): shift value used by two-term parameter-shift rules
        hybrid (bool): Whether to include classical processing inside the QNode

    Returns:
        tensor_like or tuple[list[QuantumTape], function]: The Hessian, or a
        tuple containing the Hessian tapes and post-processing function.
    """

    if any(m.return_type is qml.operation.State for m in tape.measurements):
        raise ValueError(
            "Computing the Hessian of circuits that return the state is not supported."
        )

    if any(m.return_type is qml.operation.Variance for m in tape.measurements):
        raise ValueError(
            "Computing the Hessian of variance measurements "
            "is not yet supported."
        )

    _gradient_analysis(tape)

    num_trainable_params = len(tape.trainable_params)
    if argnum is None:
        selected_argnums = list(range(num_trainable_params))
    elif isinstance(argnum, int):
        selected_argnums = [argnum]
    else:
        selected_argnums = list(argnum)

    if not selected_argnums:
        return [], lambda _: qml.math.zeros((0, 0))

    diff_methods = tape._grad_method_validation("analytic")
    selected_methods = [diff_methods[idx] for idx in selected_argnums]
    if any(method == "F" for method in selected_methods):
        raise ValueError(
            "param_shift_hessian only supports operations with analytic "
            "parameter-shift rules; unsupported operations were found."
        )

    terms_by_spec = {}
    hessian_terms = [
        [None] * num_trainable_params for _ in range(num_trainable_params)
    ]

    for row in selected_argnums:
        for col in selected_argnums:
            if col < row:
                continue

            pair_terms = []
            for coeff, spec in _second_derivative_recipe(tape, row, col, shift):
                spec = tuple(sorted(spec))
                if spec not in terms_by_spec:
                    terms_by_spec[spec] = len(terms_by_spec)
                pair_terms.append((coeff, terms_by_spec[spec]))
            hessian_terms[row][col] = pair_terms
            hessian_terms[col][row] = pair_terms

    specs = [None] * len(terms_by_spec)
    for spec, tape_idx in terms_by_spec.items():
        specs[tape_idx] = spec

    hessian_tapes = [_make_shifted_tape(tape, spec) for spec in specs]
    coefficient_matrix = np.zeros(
        (num_trainable_params, num_trainable_params, len(hessian_tapes))
    )
    for row in selected_argnums:
        for col in selected_argnums:
            for coeff, tape_idx in hessian_terms[row][col]:
                coefficient_matrix[row, col, tape_idx] = coeff

    def processing_fn(results):
        results = qml.math.stack(results)
        if len(tape.measurements) == 1:
            results = qml.math.squeeze(results, axis=1)
        coeffs = qml.math.convert_like(coefficient_matrix, results)
        qhess = qml.math.tensordot(coeffs, results, axes=[[-1], [0]])

        if qml.math.ndim(results) > 1:
            if len(tape.measurements) == 1 and not any(
                m.return_type is qml.operation.Probability
                for m in tape.measurements
            ):
                qhess = qml.math.squeeze(qhess, axis=-1)
            else:
                qhess = qml.math.moveaxis(qhess, -1, 0)

        if selected_argnums != list(range(num_trainable_params)):
            qhess = qml.math.array(
                [
                    [qhess[row, col] for col in selected_argnums]
                    for row in selected_argnums
                ]
            )

        return qhess

    return hessian_tapes, processing_fn


@param_shift_hessian.custom_qnode_wrapper
def _param_shift_hessian_qnode_wrapper(self, qnode, targs, tkwargs):
    """Include classical QNode processing in both Hessian axes."""

    hybrid = tkwargs.pop("hybrid", True)
    wrapper = super(type(self), self).default_qnode_wrapper(qnode, targs, tkwargs)
    cjac_fn = qml.transforms.classical_jacobian(qnode)

    def hessian_wrapper(*args, **kwargs):
        qhess = wrapper(*args, **kwargs)

        if not hybrid:
            return qhess

        kwargs.pop("shots", False)
        cjac = cjac_fn(*args, **kwargs)
        return _contract_hessian_with_cjac(qhess, cjac)

    return hessian_wrapper
