import json
import warnings
import platform


def deprecated(func):
    """This is a decorator which can be used to mark functions
    as deprecated. It will result in a warning being emmitted
    when the function is used."""
    def newFunc(*args, **kwargs):
        warnings.warn("Call to deprecated function {}.".format(func.__name__),
                      category=DeprecationWarning, stacklevel=2)
        return func(*args, **kwargs)
    newFunc.__name__ = func.__name__
    newFunc.__doc__ = func.__doc__
    newFunc.__dict__.update(func.__dict__)
    return newFunc


def deprecated_option(option):
    warnings.warn("Call to deprecated option {}.".format(option),
                  category=DeprecationWarning, stacklevel=2)


def java_version():
    import subprocess

    # TODO: Remove this Python 2 compatibility code if possible
    try:
        FileNotFoundError
    except NameError:
        FileNotFoundError = IOError

    try:
        res = subprocess.check_output(["java", "-version"], stderr=subprocess.STDOUT)
        res = res.decode()

    except FileNotFoundError as e:
        res = "`java -version` faild. `java` command is not found from this Python process. Please ensure Java is installed and PATH is set for `java`"

    return res


def environment_info():
    import sys
    import distro
    from .__version__ import __version__

    print("""Python version:
    {}
Java version:
    {}
tabula-py version: {}
platform: {}
uname:
    {}
linux_distribution: {}
mac_ver: {}
    """.format(
        sys.version,
        java_version().strip(),
        __version__,
        platform.platform(),
        str(platform.uname()),
        distro.linux_distribution(),
        platform.mac_ver(),
    ))


def load_template(path_or_buffer):
    '''Build tabula-py option from template file

    Args:
        path_or_buffer:
            File-like object of tabula app template or its path

    Returns:
        `obj`:list: tabula-py options
    '''

    from .wrapper import _is_file_like, _stringify_path

    path_or_buffer = _stringify_path(path_or_buffer)

    if _is_file_like(path_or_buffer):
        templates = json.load(path_or_buffer)

    else:
        with open(path_or_buffer, 'r') as f:
            templates = json.load(f)

    options = []

    for template in templates:
        if template.get('extraction_method') == 'stream':
            template['stream'] = True

        elif template.get('extraction_method') == 'lattice':
            template['lattice'] = True

        template['pages'] = template.pop('page')
        template['area'] = [template['y1'], template['x1'],
                            template['y2'], template['x2']]
        options.append(template)

    return options
