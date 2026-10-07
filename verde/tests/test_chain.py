# pylint: disable=redefined-outer-name
"""
Test the Chain class.
"""
import numpy as np
import numpy.testing as npt
import pytest

from ..chain import Chain
from ..trend import Trend
from ..scipy_bridge import ScipyGridder
from ..coordinates import grid_coordinates


@pytest.fixture()
def simple_data():
    "Some synthetic data with a trend"
    coords = grid_coordinates((1000, 5000, -5000, -1000), shape=(50, 50))
    coefs = [10, 2, -0.4]
    data = coefs[0] + coefs[1]*coords[0] + coefs[2]*coords[1]
    return coords, coefs, data


def test_chain(simple_data):
    "Test chaining a trend and a gridder"
    coords, coefs, data = simple_data
    chain = Chain([('trend', Trend(degree=1)),
                   ('gridder', ScipyGridder())]).fit(coords, data)
    npt.assert_allclose(chain.predict(coords), data)
    npt.assert_allclose(chain.named_steps['trend'].coef_, coefs)
    npt.assert_allclose(chain.named_steps['trend'].predict(coords), data)


def test_chain_weights(simple_data):
    "Use weights to account for outliers"
    coords, coefs, data = simple_data
    data_out = data.copy()
    outlier = data_out[20, 20]*50
    data_out[20, 20] += outlier
    weights = np.ones_like(data)
    weights[20, 20] = 1e-10
    chain = Chain([('trend', Trend(degree=1)),
                   ('gridder', ScipyGridder())]).fit(coords, data_out, weights)
    npt.assert_allclose(chain.named_steps['trend'].coef_, coefs)


def test_chain_grid(simple_data):
    "Make a grid of the chained predictions"
    coords, coefs, data = simple_data
    chain = Chain([('trend', Trend(degree=1)),
                   ('gridder', ScipyGridder())]).fit(coords, data)
    grid = chain.grid(shape=(50, 50))
    npt.assert_allclose(grid.scalars, data)
    trend_grid = chain.named_steps['trend'].grid(shape=(50, 50))
    npt.assert_allclose(trend_grid.scalars, data)


def test_chain_double(simple_data):
    "Chain two trends"
    coords, coefs, data = simple_data
    chain = Chain([('trend1', Trend(degree=1)),
                   ('trend2', Trend(degree=1))]).fit(coords, data)
    npt.assert_allclose(chain.predict(coords), data)
    npt.assert_allclose(chain.named_steps['trend1'].coef_, coefs)
    npt.assert_allclose(chain.named_steps['trend2'].coef_, 0, atol=1e-10)


def test_chain_fails(simple_data):
    "Check that predict fails on an unfitted chain"
    coords, _, _ = simple_data
    chain = Chain([('trend', Trend(degree=1))])
    with pytest.raises(Exception):
        chain.predict(coords)
