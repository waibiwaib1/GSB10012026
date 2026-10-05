"""Shared helpers for Fuel's command-line scripts."""
import imp
import importlib
import os


def import_external_module(name_or_path):
    """Import an external dataset module.

    External dataset modules allow 'fuel-download' and 'fuel-convert'
    to work with datasets that are not part of the Fuel distribution.
    Such a module follows the same convention as the built-in download
    and conversion modules: it defines a 'fill_subparser' function
    which accepts an 'argparse.ArgumentParser' instance, fills it with
    the arguments it needs and sets a 'func' default argument to be
    called with the parsed command-line arguments.

    Parameters
    ----------
    name_or_path : str
        Either a dotted module path (e.g. 'my_package.my_module') or
        a path to a Python file (e.g. 'my_module.py').

    Returns
    -------
    module : module
        The imported module.
    name : str
        A name for the dataset, derived from the module name or the
        file name.

    """
    if name_or_path.endswith('.py') or os.sep in name_or_path:
        path = os.path.abspath(name_or_path)
        name = os.path.splitext(os.path.basename(path))[0]
        module = imp.load_source(name, path)
    else:
        module = importlib.import_module(name_or_path)
        name = name_or_path.rsplit('.', 1)[-1]
    return module, name


def resolve_dataset_argument(args, built_in_datasets):
    """Resolve the dataset argument to a subparser-filling function.

    Inspects the first positional argument. If it refers to a built-in
    dataset, nothing happens. Otherwise, it is interpreted as an
    external dataset module (a dotted module path or a path to a
    Python file), which is imported and registered alongside the
    built-in datasets.

    Parameters
    ----------
    args : list of str
        The command-line arguments. The first element is assumed to be
        the dataset argument and is replaced in-place with the dataset
        name if an external module was given.
    built_in_datasets : dict
        Maps dataset names to their subparser-filling functions.
        External datasets are added to this mapping in-place.

    """
    if (args and not args[0].startswith('-') and
            args[0] not in built_in_datasets):
        module, name = import_external_module(args[0])
        built_in_datasets[name] = module.fill_subparser
        args[0] = name
