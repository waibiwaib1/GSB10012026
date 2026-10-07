"""
Meta-gridders that chain several operations.
"""
from sklearn.utils.validation import check_is_fitted

from .base import BaseGridder
from .coordinates import get_region
from .utils import filter_args


class Chain(BaseGridder):
    """
    Chain filtering operations to fit on each subsequent output.

    Each element of the chain is fitted on the residuals of the previous
    elements in the chain. When predicting, the predictions of each element
    of the chain are summed together.

    This is useful, for example, for removing a trend from the data before
    gridding it. Each step of the chain can be accessed individually through
    the ``named_steps`` attribute.

    Parameters
    ----------
    steps : list of (name, gridder) tuples
        List of ``(name, gridder)`` pairs where ``gridder`` is a verde
        gridder class (anything with ``fit`` and ``predict`` methods,
        including other chains). The names are used to access each step of
        the chain through the ``named_steps`` attribute.

    Attributes
    ----------
    named_steps : dict
        A dictionary version of *steps* where the keys are the names and the
        values are the gridder objects.
    region_ : tuple
        The boundaries (``[W, E, S, N]``) of the data used to fit the
        interpolator. Used as the default region for the
        :meth:`~verde.Chain.grid` and :meth:`~verde.Chain.scatter` methods.

    Examples
    --------

    The :class:`verde.Chain` can be used to remove a trend from the data
    before gridding. This is done by fitting a :class:`verde.Trend` on the
    data, a gridder on the residuals, and summing the predictions of both:

    >>> from verde import Trend, ScipyGridder, grid_coordinates
    >>> import numpy as np
    >>> coordinates = grid_coordinates((1, 5, -5, -1), shape=(5, 5))
    >>> data = 10 + 2*coordinates[0] - 0.4*coordinates[1]
    >>> chain = Chain([('trend', Trend(degree=1)),
    ...                ('gridder', ScipyGridder())])
    >>> chain.fit(coordinates, data)
    Chain(steps=[('trend', Trend(degree=1)), ('gridder', ScipyGridder())])
    >>> # The chain's prediction is the sum of the trend and the gridder
    >>> np.allclose(chain.predict(coordinates), data)
    True
    >>> # Each step can be accessed through named_steps
    >>> trend = chain.named_steps['trend']
    >>> print(', '.join(['{:.1f}'.format(i) for i in trend.coef_]))
    10.0, 2.0, -0.4
    >>> np.allclose(trend.predict(coordinates), data)
    True

    """

    def __init__(self, steps):
        self.steps = steps

    @property
    def named_steps(self):
        """
        A dictionary version of *steps* where the keys are the names and the
        values are the gridder objects.
        """
        return dict(self.steps)

    def fit(self, coordinates, data, weights=None):
        """
        Fit the chained operations to the given data.

        Each step in the chain is fitted on the residuals of the previous
        steps. The data region is captured and used as default for the
        :meth:`~verde.Chain.grid` and :meth:`~verde.Chain.scatter` methods.

        All input arrays must have the same shape.

        Parameters
        ----------
        coordinates : tuple of arrays
            Arrays with the coordinates of each data point. Should be in the
            following order: (easting, northing, vertical, ...). Only easting
            and northing will be used, all subsequent coordinates will be
            ignored.
        data : array
            The data values of each data point.
        weights : None or array
            If not None, then the weights assigned to each data point.
            Typically, this should be 1 over the data uncertainty squared.
            Ignored for steps that don't accept weights.

        Returns
        -------
        self : verde.Chain
            Returns this estimator instance for chaining operations.

        """
        self.region_ = get_region(coordinates[0], coordinates[1])
        for _, step in self.steps:
            step.fit(*filter_args(step.fit, (coordinates, data, weights)))
            data = data - step.predict(coordinates)
        return self

    def predict(self, coordinates):
        """
        Predict data on the given set of points.

        The prediction is the sum of the predictions of each step in the
        chain.

        Requires a fitted estimator (see :meth:`~verde.Chain.fit`).

        Parameters
        ----------
        coordinates : tuple of arrays
            Arrays with the coordinates of each data point. Should be in the
            following order: (easting, northing, vertical, ...). Only easting
            and northing will be used, all subsequent coordinates will be
            ignored.

        Returns
        -------
        data : array
            The data values predicted on the given points.

        """
        check_is_fitted(self, ['region_'])
        return sum(step.predict(coordinates) for _, step in self.steps)
