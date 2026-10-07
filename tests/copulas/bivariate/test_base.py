from unittest import TestCase, mock

from copulas.bivariate.base import Bivariate, CopulaTypes
from tests import compare_nested_dicts


class TestBivariate(TestCase):

    def test_from_dict(self):
        """from_dict builds the right copula subclass with the given parameters."""
        parameters = {
            'copula_type': 'FRANK',
            'tau': 0.15,
            'theta': 0.8
        }

        instance = Bivariate.from_dict(parameters)

        assert instance.copula_type == CopulaTypes.FRANK
        assert instance.tau == 0.15
        assert instance.theta == 0.8

    def test_to_dict(self):
        """to_dict returns the defining parameters of the copula."""
        instance = Bivariate('frank')
        U = [0.1, 0.5, 0.8]
        V = [0.2, 0.7, 0.1]
        instance.fit(U, V)

        expected_result = {
            'copula_type': 'FRANK',
            'tau': -0.33333333333333337,
            'theta': -3.305771759329249
        }

        result = instance.to_dict()

        assert result == expected_result

    @mock.patch('copulas.bivariate.base.json.dump')
    def test_save(self, json_mock):
        """save serializes the internal state as json into a file."""
        instance = Bivariate('frank')
        U = [0.1, 0.5, 0.8]
        V = [0.2, 0.7, 0.1]
        instance.fit(U, V)

        expected_content = {
            'copula_type': 'FRANK',
            'theta': -3.305771759329249,
            'tau': -0.33333333333333337
        }

        instance.save('test.json')

        assert json_mock.called
        compare_nested_dicts(json_mock.call_args[0][0], expected_content)

    @mock.patch('builtins.open', new_callable=mock.mock_open)
    @mock.patch('copulas.bivariate.base.json.load')
    def test_load(self, json_mock, file_mock):
        """load recreates an instance from a saved file."""
        json_mock.return_value = {
            'copula_type': 'FRANK',
            'tau': -0.33333333333333337,
            'theta': -3.305771759329249
        }

        instance = Bivariate.load('somefile.json')

        assert instance.copula_type == CopulaTypes.FRANK
        assert instance.tau == -0.33333333333333337
        assert instance.theta == -3.305771759329249

    def test_save_load_roundtrip(self, tmp_path=None):
        """An instance saved and loaded keeps its internal state."""
        import tempfile

        instance = Bivariate('clayton')
        U = [0.1, 0.2, 0.3, 0.4]
        V = [0.5, 0.6, 0.5, 0.8]
        instance.fit(U, V)

        with tempfile.NamedTemporaryFile(suffix='.json') as temp_file:
            instance.save(temp_file.name)
            loaded = Bivariate.load(temp_file.name)

        assert loaded.copula_type == CopulaTypes.CLAYTON
        assert loaded.tau == instance.tau
        assert loaded.theta == instance.theta
