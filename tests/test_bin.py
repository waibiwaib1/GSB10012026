import os
import shutil
import sys
import tempfile

from numpy.testing import assert_raises

from fuel.bin import import_external_module, resolve_dataset_argument

EXTERNAL_MODULE = '''
def fill_subparser(subparser):
    subparser.set_defaults(func=lambda: None)
    return subparser
'''


class TestImportExternalModule(object):
    def setup(self):
        self.tempdir = tempfile.mkdtemp()
        self.module_path = os.path.join(self.tempdir, 'my_dataset.py')
        with open(self.module_path, 'w') as f:
            f.write(EXTERNAL_MODULE)
        sys.path.insert(0, self.tempdir)

    def teardown(self):
        sys.path.remove(self.tempdir)
        shutil.rmtree(self.tempdir)
        sys.modules.pop('my_dataset', None)

    def test_import_from_file_path(self):
        module, name = import_external_module(self.module_path)
        assert name == 'my_dataset'
        assert hasattr(module, 'fill_subparser')

    def test_import_from_dotted_path(self):
        module, name = import_external_module('my_dataset')
        assert name == 'my_dataset'
        assert hasattr(module, 'fill_subparser')

    def test_import_nonexistent_raises(self):
        assert_raises(IOError, import_external_module,
                      os.path.join(self.tempdir, 'nonexistent.py'))
        assert_raises(ImportError, import_external_module, 'nonexistent')


class TestResolveDatasetArgument(object):
    def setup(self):
        self.tempdir = tempfile.mkdtemp()
        self.module_path = os.path.join(self.tempdir, 'my_dataset.py')
        with open(self.module_path, 'w') as f:
            f.write(EXTERNAL_MODULE)

    def teardown(self):
        shutil.rmtree(self.tempdir)
        sys.modules.pop('my_dataset', None)

    def test_built_in_dataset_is_untouched(self):
        datasets = {'iris': lambda subparser: None}
        args = ['iris']
        resolve_dataset_argument(args, datasets)
        assert args == ['iris']
        assert list(datasets.keys()) == ['iris']

    def test_external_module_from_file_path(self):
        datasets = {}
        args = [self.module_path]
        resolve_dataset_argument(args, datasets)
        assert args == ['my_dataset']
        assert 'my_dataset' in datasets

    def test_no_arguments(self):
        datasets = {'iris': lambda subparser: None}
        resolve_dataset_argument([], datasets)
        assert list(datasets.keys()) == ['iris']

    def test_option_first_argument(self):
        datasets = {'iris': lambda subparser: None}
        args = ['-h']
        resolve_dataset_argument(args, datasets)
        assert args == ['-h']
        assert list(datasets.keys()) == ['iris']
