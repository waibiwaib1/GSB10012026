import pickle


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

    def save(self, filename):
        """Save the internal state of a copula in the specified file.

        Args:
            filename: 'str', path to the file where the copula will be serialized.
        """
        with open(filename, 'wb') as f:
            pickle.dump(self, f)

    @classmethod
    def load(cls, filename):
        """Load a copula previously saved with the save method.

        Args:
            filename: 'str', path to the file containing the serialized copula.

        Returns:
            The loaded copula instance, with its internal state restored.
        """
        with open(filename, 'rb') as f:
            return pickle.load(f)
