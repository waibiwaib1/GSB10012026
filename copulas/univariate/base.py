import json


class Univariate(object):
    """ Abstract class for representing univariate distributions """

    def __init__(self):
        pass

    def fit(self):
        """ fits a univariate model and updates parameters """
        raise NotImplementedError

    def get_pdf(self, value):
        """ given a value, returns corresponding pdf value """
        raise NotImplementedError

    def get_cdf(self, value):
        """ given a value returns corresponding cdf value """
        raise NotImplementedError

    def inverse_cdf(self, value):
        """ given a cdf value, returns a value in original space """
        raise NotImplementedError

    def sample(self):
        """ returns new data point based on model """
        raise NotImplementedError

    def to_dict(self):
        """Return a `dict` with the parameters needed to replicate this distribution."""
        raise NotImplementedError

    @classmethod
    def from_dict(cls, param_dict):
        """Create a new instance from a `dict` of parameters."""
        raise NotImplementedError

    def save(self, filename):
        """Save the internal state of the distribution in the given file as JSON."""
        content = self.to_dict()
        with open(filename, 'w') as file_handle:
            json.dump(content, file_handle)

    @classmethod
    def load(cls, distribution_path):
        """Create a new instance from a JSON file."""
        with open(distribution_path) as file_handle:
            distribution_dict = json.load(file_handle)

        return cls.from_dict(distribution_dict)
