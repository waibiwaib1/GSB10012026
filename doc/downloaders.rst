.. _downloaders:

Downloaders: Customizing the download
=====================================

By default, :meth:`pooch.Pooch.fetch` and :meth:`pooch.retrieve` will detect
the download protocol from the given URL (HTTP, FTP, or SFTP) and use the
appropriate download method.
Sometimes this is not enough: some servers require logins, redirections, or
other non-standard operations.
To get around this, use the ``downloader`` argument of
:meth:`~pooch.Pooch.fetch` and :meth:`~pooch.retrieve`.

Downloaders are Python *callable objects*  (like functions or classes with a
``__call__`` method) and must have the following format:

.. code:: python

    def mydownloader(url, output_file, pooch):
        '''
        Download a file from the given URL to the given local file.

        The function **must** take as arguments (in order):

        url : str
            The URL to the file you want to download.
        output_file : str or file-like object
            Path (and file name) to which the file will be downloaded.
        pooch : pooch.Pooch
            The instance of the Pooch class that is calling this function.

        No return value is required.
        '''
        ...

Pooch provides downloaders for HTTP, FTP, and SFTP that support authentication
and optionally printing progress bars.
See :ref:`api` for a list of available downloaders.

Common uses of downloaders include:

* Passing :ref:`login credentials <authentication>` to HTTP and FTP servers
* Printing :ref:`progress bars <progressbars>`
* Downloading files from :ref:`DOI-based repositories <doi_downloads>`


.. _doi_downloads:

Digital Object Identifiers (DOIs)
---------------------------------

Open-access data repositories often issue Digital Object Identifiers (DOIs)
for data which provide a stable link and citation point.
Files hosted in these repositories can be downloaded by Pooch using the
:class:`pooch.DOIDownloader`, which resolves the DOI to the actual download
URL through the repository's public API.
To use it, specify the file in the registry using the format
``doi:{DOI}/{file name}``:

.. code:: python

    POOCH = pooch.create(
        path=pooch.os_cache("myproject"),
        base_url="doi:10.5281/zenodo.4924875/",
        registry={"tiny-data.txt": "baee0894dba14b12085eacb204284b97e362f4f3e5a5807693cc90ef415c1b2d"},
    )
    fname = POOCH.fetch("tiny-data.txt")

Currently supported repositories are `figshare <https://www.figshare.com>`__
and `Zenodo <https://www.zenodo.org>`__.


Creating your own downloaders
-----------------------------

If your use case is not covered by our downloaders, you can implement your own.
:meth:`pooch.Pooch.fetch` and :func:`pooch.retrieve` will accept any *callable
obejct* that has the signature specified above. As an example, consider the
case in which the login credentials need to be provided to a site that is
redirected from the original download URL:

.. code:: python

    import requests


    def redirect_downloader(url, output_file, pooch):
        """
        Download after following a redirection.
        """
        # Get the credentials from the user's environment
        username = os.environ.get("SOMESITE_USERNAME")
        password = os.environ.get("SOMESITE_PASSWORD")
        # Make a request that will redirect to the login page
        login = requests.get(url)
        # Provide the credentials and download from the new URL
        download = HTTPDownloader(auth=(username, password))
        download(login.url, output_file, mypooch)


    def fetch_protected_data():
        """
        Fetch a file from a server that requires authentication
        """
        fname = GOODBOY.fetch("some-data.csv", downloader=redirect_downloader)
        data = pandas.read_csv(fname)
        return data
