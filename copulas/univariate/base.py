import pickle


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

    def save(self, filename):
        """Save the internal state of a distribution in the specified file.

        Args:
            filename: 'str', path to the file where the distribution will be serialized.
        """
        with open(filename, 'wb') as f:
            pickle.dump(self, f)

    @classmethod
    def load(cls, filename):
        """Load a distribution previously saved with the save method.

        Args:
            filename: 'str', path to the file containing the serialized distribution.

        Returns:
            The loaded distribution instance, with its internal state restored.
        """
        with open(filename, 'rb') as f:
            return pickle.load(f)
