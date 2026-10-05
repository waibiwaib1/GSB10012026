# Copyright (c) 2018 The Pooch Developers.
# Distributed under the terms of the BSD 3-Clause License.
# SPDX-License-Identifier: BSD-3-Clause
#
# This code is part of the Fatiando a Terra project (https://www.fatiando.org)
#
"""
Test the downloader classes and functions separately from the Pooch core.
"""
import os
import sys
from tempfile import TemporaryDirectory

import pytest

try:
    import tqdm
except ImportError:
    tqdm = None

try:
    import paramiko
except ImportError:
    paramiko = None

from ..downloaders import (
    HTTPDownloader,
    FTPDownloader,
    SFTPDownloader,
    DOIDownloader,
    choose_downloader,
    zenodo_url_to_api,
)
from .utils import pooch_test_url, check_large_data, check_tiny_data, data_over_ftp


# FTP doesn't work on Travis CI so need to be able to skip tests there
ON_TRAVIS = bool(os.environ.get("TRAVIS", None))
BASEURL = pooch_test_url()


def test_unsupported_protocol():
    "Should raise ValueError when protocol not in {'https', 'http', 'ftp'}"
    with pytest.raises(ValueError):
        choose_downloader("httpup://some-invalid-url.com")


def test_zenodo_url_to_api():
    "Convert DOI URLs to Zenodo API URLs"
    assert (
        zenodo_url_to_api("doi:10.5281/zenodo.1234567")
        == "https://zenodo.org/api/records/1234567"
    )
    assert (
        zenodo_url_to_api("doi:https://doi.org/10.5281/zenodo.1234567")
        == "https://zenodo.org/api/records/1234567"
    )


def test_zenodo_url_to_api_invalid_doi():
    "Should raise if the DOI is not a Zenodo DOI"
    with pytest.raises(ValueError):
        zenodo_url_to_api("doi:10.6084/m9.figshare.12345")


def test_doi_downloader_choose():
    "choose_downloader should return a DOIDownloader for doi: URLs"
    downloader = choose_downloader("doi:10.5281/zenodo.1234567/file.txt")
    assert isinstance(downloader, DOIDownloader)


def test_doi_downloader_missing_file_name():
    "Should raise if no file name is given after the DOI"
    downloader = DOIDownloader()
    with TemporaryDirectory() as local_store:
        outfile = os.path.join(local_store, "file.txt")
        with pytest.raises(ValueError):
            downloader("doi:10.5281/zenodo.1234567", outfile, None)


def test_doi_downloader(monkeypatch):
    "Test the DOI downloader with a mocked Zenodo API"

    class FakeResponse:
        def __init__(self, json_data=None):
            self._json_data = json_data

        def raise_for_status(self):
            return None

        def json(self):
            return self._json_data

    def fake_get(url, **kwargs):  # pylint: disable=unused-argument
        if url == "https://zenodo.org/api/records/1234567":
            return FakeResponse(
                {
                    "files": [
                        {"key": "other-file.txt", "links": {"self": "other"}},
                        {
                            "key": "tiny-data.txt",
                            "links": {"self": BASEURL + "tiny-data.txt"},
                        },
                    ]
                }
            )
        return FakeResponse()

    # Patch only the API request; the file itself is downloaded from the
    # local test server.
    import pooch.downloaders as downloaders_module

    original_get = downloaders_module.requests.get
    call_count = {"n": 0}

    def patched_get(url, **kwargs):
        if "zenodo.org/api" in url:
            call_count["n"] += 1
            return fake_get(url, **kwargs)
        return original_get(url, **kwargs)

    monkeypatch.setattr(downloaders_module.requests, "get", patched_get)

    with TemporaryDirectory() as local_store:
        downloader = DOIDownloader()
        outfile = os.path.join(local_store, "tiny-data.txt")
        downloader(
            "doi:10.5281/zenodo.1234567/tiny-data.txt", outfile, None
        )
        check_tiny_data(outfile)
    assert call_count["n"] == 1


def test_doi_downloader_file_not_found(monkeypatch):
    "Should raise if the requested file is not in the Zenodo record"

    class FakeResponse:
        def raise_for_status(self):
            return None

        def json(self):
            return {"files": [{"key": "other-file.txt"}]}

    import pooch.downloaders as downloaders_module

    monkeypatch.setattr(
        downloaders_module.requests,
        "get",
        lambda url, **kwargs: FakeResponse(),
    )

    downloader = DOIDownloader()
    with TemporaryDirectory() as local_store:
        outfile = os.path.join(local_store, "missing.txt")
        with pytest.raises(ValueError):
            downloader("doi:10.5281/zenodo.1234567/missing.txt", outfile, None)


def test_ftp_downloader(ftpserver):
    "Test ftp downloader"
    with data_over_ftp(ftpserver, "tiny-data.txt") as url:
        with TemporaryDirectory() as local_store:
            downloader = FTPDownloader(port=ftpserver.server_port)
            outfile = os.path.join(local_store, "tiny-data.txt")
            downloader(url, outfile, None)
            check_tiny_data(outfile)


@pytest.mark.skipif(paramiko is None, reason="requires paramiko to run SFTP")
@pytest.mark.skipif(ON_TRAVIS, reason="SFTP is not allowed on Travis CI")
def test_sftp_downloader():
    "Test sftp downloader"
    with TemporaryDirectory() as local_store:
        downloader = SFTPDownloader(username="demo", password="password")
        url = "sftp://test.rebex.net/pub/example/pocketftp.png"
        outfile = os.path.join(local_store, "pocketftp.png")
        downloader(url, outfile, None)
        assert os.path.exists(outfile)


@pytest.mark.skipif(paramiko is None, reason="requires paramiko to run SFTP")
@pytest.mark.skipif(ON_TRAVIS, reason="SFTP is not allowed on Travis CI")
def test_sftp_downloader_fail_if_file_object():
    "Downloader should fail when a file object rather than string is passed"
    with TemporaryDirectory() as local_store:
        downloader = SFTPDownloader(username="demo", password="password")
        url = "sftp://test.rebex.net/pub/example/pocketftp.png"
        outfile = os.path.join(local_store, "pocketftp.png")
        with open(outfile, "wb") as outfile_obj:
            with pytest.raises(TypeError):
                downloader(url, outfile_obj, None)


@pytest.mark.skipif(paramiko is not None, reason="paramiko must be missing")
def test_sftp_downloader_fail_if_paramiko_missing():
    "test must fail if paramiko is not installed"
    with pytest.raises(ValueError) as exc:
        SFTPDownloader()
    assert "'paramiko'" in str(exc.value)


@pytest.mark.skipif(tqdm is not None, reason="tqdm must be missing")
@pytest.mark.parametrize("downloader", [HTTPDownloader, FTPDownloader, SFTPDownloader])
def test_downloader_progressbar_fails(downloader):
    "Make sure an error is raised if trying to use progressbar without tqdm"
    with pytest.raises(ValueError) as exc:
        downloader(progressbar=True)
    assert "'tqdm'" in str(exc.value)


@pytest.mark.skipif(tqdm is None, reason="requires tqdm")
def test_downloader_progressbar(capsys):
    "Setup a downloader function that prints a progress bar for fetch"
    download = HTTPDownloader(progressbar=True)
    with TemporaryDirectory() as local_store:
        fname = "large-data.txt"
        url = BASEURL + fname
        outfile = os.path.join(local_store, "large-data.txt")
        download(url, outfile, None)
        # Read stderr and make sure the progress bar is printed only when told
        captured = capsys.readouterr()
        printed = captured.err.split("\r")[-1].strip()
        assert len(printed) == 79
        if sys.platform == "win32":
            progress = "100%|####################"
        else:
            progress = "100%|████████████████████"
        # Bar size is not always the same so can't reliably test the whole bar.
        assert printed[:25] == progress
        # Check that the downloaded file has the right content
        check_large_data(outfile)


@pytest.mark.skipif(tqdm is None, reason="requires tqdm")
def test_downloader_progressbar_ftp(capsys, ftpserver):
    "Setup an FTP downloader function that prints a progress bar for fetch"
    with data_over_ftp(ftpserver, "tiny-data.txt") as url:
        download = FTPDownloader(progressbar=True, port=ftpserver.server_port)
        with TemporaryDirectory() as local_store:
            outfile = os.path.join(local_store, "tiny-data.txt")
            download(url, outfile, None)
            # Read stderr and make sure the progress bar is printed only when
            # told
            captured = capsys.readouterr()
            printed = captured.err.split("\r")[-1].strip()
            assert len(printed) == 79
            if sys.platform == "win32":
                progress = "100%|####################"
            else:
                progress = "100%|████████████████████"
            # Bar size is not always the same so can't reliably test the whole
            # bar.
            assert printed[:25] == progress
            # Check that the file was actually downloaded
            check_tiny_data(outfile)


@pytest.mark.skipif(tqdm is None, reason="requires tqdm")
@pytest.mark.skipif(paramiko is None, reason="requires paramiko")
@pytest.mark.skipif(ON_TRAVIS, reason="SFTP is not allowed on Travis CI")
def test_downloader_progressbar_sftp(capsys):
    "Setup an SFTP downloader function that prints a progress bar for fetch"
    downloader = SFTPDownloader(progressbar=True, username="demo", password="password")
    with TemporaryDirectory() as local_store:
        url = "sftp://test.rebex.net/pub/example/pocketftp.png"
        outfile = os.path.join(local_store, "pocketftp.png")
        downloader(url, outfile, None)
        # Read stderr and make sure the progress bar is printed only when told
        captured = capsys.readouterr()
        printed = captured.err.split("\r")[-1].strip()
        assert len(printed) == 79
        if sys.platform == "win32":
            progress = "100%|####################"
        else:
            progress = "100%|████████████████████"
        # Bar size is not always the same so can't reliably test the whole bar.
        assert printed[:25] == progress
        # Check that the file was actually downloaded
        assert os.path.exists(outfile)


def test_downloader_arbitrary_progressbar(capsys):
    "Setup a downloader function with an arbitrary progress bar class."

    class MinimalProgressDisplay:
        """A minimalist replacement for tqdm.tqdm"""

        def __init__(self, total):
            self.count = 0
            self.total = total

        def __repr__(self):
            """represent current completion"""
            return str(self.count) + "/" + str(self.total)

        def render(self):
            """print self.__repr__ to stderr"""
            print(f"\r{self}", file=sys.stderr, end="")

        def update(self, i):
            """modify completion and render"""
            self.count = i
            self.render()

        def reset(self):
            """set counter to 0"""
            self.count = 0

        @staticmethod
        def close():
            """print a new empty line"""
            print("", file=sys.stderr)

    pbar = MinimalProgressDisplay(total=None)
    download = HTTPDownloader(progressbar=pbar)
    with TemporaryDirectory() as local_store:
        fname = "large-data.txt"
        url = BASEURL + fname
        outfile = os.path.join(local_store, "large-data.txt")
        download(url, outfile, None)
        # Read stderr and make sure the progress bar is printed only when told
        captured = capsys.readouterr()
        printed = captured.err.split("\r")[-1].strip()

        progress = "336/336"
        assert printed == progress

        # Check that the downloaded file has the right content
        check_large_data(outfile)
