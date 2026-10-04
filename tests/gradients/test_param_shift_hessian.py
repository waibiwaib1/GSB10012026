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
"""Tests for the gradients.param_shift_hessian module."""
import pytest

import pennylane as qml
from pennylane import numpy as np


class TestParamShiftHessian:
    """Tests for the parameter-shift Hessian transform."""

    def test_two_term_probability_hessian(self, tol):
        """Test RX/RY gates with a two-term shift rule and probability output."""
        dev = qml.device("default.qubit", wires=2)

        @qml.qnode(dev)
        def circuit(weights):
            qml.RX(weights[0], wires=0)
            qml.RY(weights[1], wires=0)
            return qml.probs(wires=0)

        weights = np.array([0.1, 0.2], requires_grad=True)

        res = qml.gradients.param_shift_hessian(circuit)(weights)
        expected = qml.jacobian(qml.jacobian(circuit))(weights)

        assert res.shape == (2, 2, 2)
        assert np.allclose(res, expected, atol=tol, rtol=0)

    def test_two_term_scalar_hessian(self, tol):
        """Test expectation values for gates with two-term shift rules."""
        dev = qml.device("default.qubit", wires=2)

        @qml.qnode(dev)
        def circuit(weights):
            qml.RX(weights[0], wires=0)
            qml.RY(weights[1], wires=1)
            qml.CNOT(wires=[0, 1])
            return qml.expval(qml.PauliZ(1))

        weights = np.array([0.4, -0.3], requires_grad=True)

        res = qml.gradients.param_shift_hessian(circuit)(weights)
        expected = qml.jacobian(qml.jacobian(circuit))(weights)

        assert res.shape == (2, 2)
        assert np.allclose(res, expected, atol=tol, rtol=0)

    def test_four_term_hessian(self, tol):
        """Test a CRX gate with a four-term shift rule."""
        dev = qml.device("default.qubit", wires=2)

        @qml.qnode(dev)
        def circuit(weight):
            qml.PauliX(0)
            qml.CRX(weight, wires=[0, 1])
            return qml.expval(qml.PauliZ(1))

        weight = np.array(0.37, requires_grad=True)

        res = qml.gradients.param_shift_hessian(circuit)(weight)
        expected = qml.grad(qml.grad(circuit))(weight)

        assert np.allclose(res, expected, atol=tol, rtol=0)

    def test_two_four_term_parameters(self, tol):
        """Test multiple gates with four-term shift rules."""
        dev = qml.device("default.qubit", wires=3)

        @qml.qnode(dev)
        def circuit(weights):
            qml.PauliX(0)
            qml.CRX(weights[0], wires=[0, 1])
            qml.CRX(weights[1], wires=[1, 2])
            return qml.expval(qml.PauliZ(2))

        weights = np.array([0.4, -0.2], requires_grad=True)

        res = qml.gradients.param_shift_hessian(circuit)(weights)
        expected = qml.jacobian(qml.jacobian(circuit))(weights)

        assert np.allclose(res, expected, atol=tol, rtol=0)

    def test_two_term_minimal_tape_count(self):
        """Two two-term gates use two tapes per diagonal and four cross tapes."""
        dev = qml.device("default.qubit", wires=2)

        @qml.qnode(dev)
        def circuit(weights):
            qml.RX(weights[0], wires=0)
            qml.RY(weights[1], wires=0)
            return qml.expval(qml.PauliZ(0))

        weights = np.array([0.1, 0.2], requires_grad=True)
        tapes, _ = qml.gradients.param_shift_hessian(circuit).construct((weights,), {})

        assert len(tapes) == 8

    def test_four_term_minimal_tape_count(self):
        """Period folding and coefficient cancellation give four evaluations."""
        dev = qml.device("default.qubit", wires=2)

        @qml.qnode(dev)
        def circuit(weight):
            qml.PauliX(0)
            qml.CRX(weight, wires=[0, 1])
            return qml.expval(qml.PauliZ(1))

        weight = np.array(0.37, requires_grad=True)
        tapes, _ = qml.gradients.param_shift_hessian(circuit).construct((weight,), {})

        assert len(tapes) == 4

    def test_two_four_term_tape_count(self):
        """Two four-term gates use four tapes per diagonal and 16 cross tapes."""
        dev = qml.device("default.qubit", wires=3)

        @qml.qnode(dev)
        def circuit(weights):
            qml.PauliX(0)
            qml.CRX(weights[0], wires=[0, 1])
            qml.CRX(weights[1], wires=[1, 2])
            return qml.expval(qml.PauliZ(2))

        weights = np.array([0.4, -0.2], requires_grad=True)
        tapes, _ = qml.gradients.param_shift_hessian(circuit).construct((weights,), {})

        assert len(tapes) == 24

    def test_state_measurement_error(self):
        """State measurements are unsupported."""
        dev = qml.device("default.qubit", wires=1)

        @qml.qnode(dev)
        def circuit(weight):
            qml.RX(weight, wires=0)
            return qml.state()

        with pytest.raises(ValueError, match="state is not supported"):
            qml.gradients.param_shift_hessian(circuit)(0.4)


class TestParamShiftHessianDifferentiation:
    """Tests for third derivatives through the Hessian transform."""

    def test_autograd(self, tol):
        """The Hessian transform is differentiable with Autograd."""
        dev = qml.device("default.qubit.autograd", wires=1)

        @qml.qnode(dev)
        def circuit(weight):
            qml.RX(weight, wires=0)
            return qml.expval(qml.PauliZ(0))

        weight = np.array(0.37, requires_grad=True)
        res = qml.grad(qml.gradients.param_shift_hessian(circuit))(weight)

        assert np.allclose(res, np.sin(weight), atol=tol, rtol=0)

    @pytest.mark.slow
    def test_tf(self, tol):
        """The Hessian transform is differentiable with TensorFlow."""
        tf = pytest.importorskip("tensorflow")
        dev = qml.device("default.qubit.tf", wires=1)

        @qml.qnode(dev, interface="tf")
        def circuit(weight):
            qml.RX(weight, wires=0)
            return qml.expval(qml.PauliZ(0))

        weight = tf.Variable(0.37, dtype=tf.float64)
        with tf.GradientTape() as tape:
            hessian = qml.gradients.param_shift_hessian(circuit)(weight)
        res = tape.gradient(hessian, weight)

        assert np.allclose(res, np.sin(0.37), atol=tol, rtol=0)

    def test_torch(self, tol):
        """The Hessian transform is differentiable with PyTorch."""
        torch = pytest.importorskip("torch")
        dev = qml.device("default.qubit.torch", wires=1)

        @qml.qnode(dev, interface="torch")
        def circuit(weight):
            qml.RX(weight, wires=0)
            return qml.expval(qml.PauliZ(0))

        weight = torch.tensor(0.37, dtype=torch.float64, requires_grad=True)
        hessian = qml.gradients.param_shift_hessian(circuit)(weight)
        hessian.backward()

        assert np.allclose(weight.grad.detach().numpy(), np.sin(0.37), atol=tol, rtol=0)
