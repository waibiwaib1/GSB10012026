"""
General utilities.
"""
import os
from inspect import signature


def get_home():
    """
    Get the path of the verde home directory.

    Defaults to ``$HOME/.verde``.

    If the folder doesn't already exist, it will be created.

    Returns
    -------
    path : str
        The path of the home directory.

    """
    home = os.path.abspath(os.environ.get('HOME'))
    verde_home = os.path.join(home, '.verde')
    os.makedirs(verde_home, exist_ok=True)
    return verde_home


def get_data_dir():
    """
    Get the path of the verde data directory.

    Defaults to ``get_home()/data``.

    If the folder doesn't already exist, it will be created.

    Returns
    -------
    path : str
        The path of the data directory.

    """
    data_dir = os.path.join(get_home(), 'data')
    os.makedirs(data_dir, exist_ok=True)
    return data_dir


def filter_args(func, args):
    """
    Filter the arguments in args to only those that func can receive.

    Assumes that the arguments of *func* are in the same order as in *args*.

    Parameters
    ----------
    func : callable
        The callable (function, method, etc).
    args : tuple
        The full set of arguments.

    Returns
    -------
    filtered_args : tuple
        The arguments that *func* can receive.

    Examples
    --------

    >>> def func1(data, weights):
    ...     return data
    >>> def func2(data, weights, extra=None):
    ...     return data
    >>> args = (1, 2, 3)
    >>> print(filter_args(func1, args))
    (1, 2)
    >>> print(filter_args(func2, args))
    (1, 2, 3)

    """
    params = list(signature(func).parameters.values())
    n_args = len([param for param in params
                  if param.kind in (param.POSITIONAL_ONLY,
                                    param.POSITIONAL_OR_KEYWORD)])
    return tuple(args[:n_args])
