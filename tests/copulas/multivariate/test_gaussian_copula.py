import warnings
from unittest import TestCase, mock

import pandas as pd

from copulas.multivariate.gaussian import GaussianMultivariate
from tests import compare_nested_dicts


class TestGaussianCopula(TestCase):

    def test_deprecation_warnings(self):
        """After fitting, Gaussian copula can produce new samples warningless."""
        # Setup
        copula = GaussianMultivariate()
        data = pd.read_csv('data/iris.data.csv')

        # Run
        with warnings.catch_warnings(record=True) as warns:
            copula.fit(data)
            result = copula.sample(10)

            # Check
            assert len(warns) == 0
            assert len(result) == 10

    def _get_parameters(self):
        cov_matrix = [
            [1.006711409395973, -0.11010327176239865, 0.8776048563471857, 0.823443255069628],
            [-0.11010327176239865, 1.006711409395972, -0.4233383520816991, -0.3589370029669186],
            [0.8776048563471857, -0.4233383520816991, 1.006711409395973, 0.9692185540781536],
            [0.823443255069628, -0.3589370029669186, 0.9692185540781536, 1.0067114093959735]
        ]
        return {
            'means': [
                -3.315866100213801e-16,
                -7.815970093361102e-16,
                2.842170943040401e-16,
                -2.3684757858670006e-16
            ],
            'cov_matrix': cov_matrix,
            'distribs': {
                'feature_01': {'mean': 5.843333333333334, 'std': 0.8253012917851409},
                'feature_02': {'mean': 3.0540000000000003, 'std': 0.4321465800705435},
                'feature_03': {'mean': 3.758666666666666, 'std': 1.7585291834055212},
                'feature_04': {'mean': 1.1986666666666668, 'std': 0.7606126185881716}
            }
        }

    def test_to_dict(self):
        """to_dict returns the parameters needed to replicate the copula."""
        copula = GaussianMultivariate()
        data = pd.read_csv('data/iris.data.csv')
        copula.fit(data)

        result = copula.to_dict()

        compare_nested_dicts(result, self._get_parameters())

    def test_from_dict(self):
        """from_dict recreates a working copula from parameters."""
        parameters = self._get_parameters()

        copula = GaussianMultivariate.from_dict(parameters)

        assert copula.means == parameters['means']
        assert (copula.cov_matrix == parameters['cov_matrix']).all()
        for name, distrib in copula.distribs.items():
            assert distrib.to_dict() == parameters['distribs'][name]

        assert copula.sample(10).all().all()

    @mock.patch('copulas.multivariate.base.json.dump')
    def test_save(self, json_mock):
        """save serializes the internal state as json into a file."""
        instance = GaussianMultivariate()
        data = pd.read_csv('data/iris.data.csv')
        instance.fit(data)

        instance.save('test.json')

        assert json_mock.called
        compare_nested_dicts(json_mock.call_args[0][0], self._get_parameters())

    @mock.patch('builtins.open', new_callable=mock.mock_open)
    @mock.patch('copulas.multivariate.base.json.load')
    def test_load(self, json_mock, file_mock):
        """load recreates an instance from a saved file."""
        parameters = self._get_parameters()
        json_mock.return_value = parameters

        instance = GaussianMultivariate.load('somefile.json')

        assert instance.means == parameters['means']
        assert (instance.cov_matrix == parameters['cov_matrix']).all()
        for name, distrib in instance.distribs.items():
            assert distrib.to_dict() == parameters['distribs'][name]

    def test_save_load_roundtrip(self):
        """A fitted gaussian copula can be saved to a file and loaded back."""
        import tempfile

        copula = GaussianMultivariate()
        data = pd.read_csv('data/iris.data.csv')
        copula.fit(data)

        with tempfile.NamedTemporaryFile(suffix='.json') as temp_file:
            copula.save(temp_file.name)
            loaded = GaussianMultivariate.load(temp_file.name)

        compare_nested_dicts(loaded.to_dict(), copula.to_dict())
