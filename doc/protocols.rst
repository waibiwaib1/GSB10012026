.. _protocols:

Download protocols
==================

Pooch supports the HTTP, FTP, and SFTP protocols by default.
It will **automatically detect** the correct protocol from the URL and use the
appropriate download method.

.. note::

    To download files over SFTP,
    `paramiko <https://github.com/paramiko/paramiko>`__ needs to be installed.

For example, if our data were hosted on an FTP server, we could use the
following setup:

.. code:: python

    POOCH = pooch.create(
        path=pooch.os_cache("plumbus"),
        # Use an FTP server instead of HTTP. The rest is all the same.
        base_url="ftp://garage-basement.org/{version}/",
        version=version,
        version_dev="master",
        registry={
            "c137.csv": "19uheidhlkjdwhoiwuhc0uhcwljchw9ochwochw89dcgw9dcgwc",
            "cronen.csv": "1upodh2ioduhw9celdjhlfvhksgdwikdgcowjhcwoduchowjg8w",
        },
    )


    def fetch_c137():
        """
        Load the C-137 sample data as a pandas.DataFrame (over FTP this time).
        """
        fname = POOCH.fetch("c137.csv")
        data = pandas.read_csv(fname)
        return data

You can even specify custom functions for the download or login credentials for
**authentication**. See :ref:`downloaders` for more information.

Downloading from Zenodo using a DOI
-----------------------------------

Pooch can also download files from data repositories hosted on
`Zenodo <https://zenodo.org>`__ using their DOI instead of an absolute URL.
To do so, register the file with a URL in the format
``"doi:<doi>/<file_name>"`` where ``<doi>`` is the DOI of the Zenodo record
and ``<file_name>`` is the name of the file attached to the record.
For example:

.. code:: python

    POOCH = pooch.create(
        path=pooch.os_cache("plumbus"),
        base_url="",
        registry={
            "c137.csv": "19uheidhlkjdwhoiwuhc0uhcwljchw9ochwochw89dcgw9dcgwc",
        },
        urls={
            "c137.csv": "doi:10.5281/zenodo.3629437/c137.csv",
        },
    }

Pooch will query the Zenodo API to find the download link of the file and then
fetch it over HTTPS.
The same can be done when using :func:`pooch.retrieve`:

.. code:: python

    pooch.retrieve(
        url="doi:10.5281/zenodo.3629437/c137.csv",
        known_hash="19uheidhlkjdwhoiwuhc0uhcwljchw9ochwochw89dcgw9dcgwc",
    )

To pass extra arguments (like authentication or a progress bar), use
:class:`pooch.DOIDownloader` explicitly with :meth:`pooch.Pooch.fetch` or
:func:`pooch.retrieve`.
