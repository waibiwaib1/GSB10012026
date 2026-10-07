import json


class Multivariate(object):
    """ Abstract class for a multi-variate copula object """

    def __init__(self):
        """ initialize copula object """

    def fit(self, data):
        """ Fits a model to the data and updates the parameters """
        raise NotImplementedError

    def infer(self, values):
        """ Takes in subset of values and predicts the rest """
        raise NotImplementedError

    def get_pdf(self):
        """ returns pdf of model """
        raise NotImplementedError

    def get_cdf(self):
        """ returns cdf of model """
        raise NotImplementedError

    def sample(self, num_rows=1):
        """ returns a new data point generated from model """
        raise NotImplementedError

    def to_dict(self):
        """Return a `dict` with the parameters needed to replicate this copula."""
        raise NotImplementedError

    @classmethod
    def from_dict(cls, copula_dict):
        """Create a new instance from a `dict` of parameters."""
        raise NotImplementedError

    def save(self, filename):
        """Save the internal state of the copula in the specified file as JSON.

        Args:
            filename: `str` path to the file in which the copula will be saved.
        """
        content = self.to_dict()
        with open(filename, 'w') as file_handle:
            json.dump(content, file_handle)

    @classmethod
    def load(cls, copula_path):
        """Create a new instance from a JSON file.

        Args:
            copula_path: `str` path to the file from which the copula will be loaded.

        Returns:
            Multivariate: Instance of the copula stored in the given file.
        """
        with open(copula_path) as file_handle:
            copula_dict = json.load(file_handle)

        return cls.from_dict(copula_dict)
