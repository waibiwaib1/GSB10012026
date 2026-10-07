===================
Python Code Binding
===================

The ``PycodeSerializer`` renders an object tree back to executable python code.
This is useful for converting XML or JSON samples into fixtures generated with
your xsdata models.

.. code-block:: python

    from pathlib import Path

    from xsdata.formats.dataclass.context import XmlContext
    from xsdata.formats.dataclass.parsers import XmlParser
    from tests.fixtures.books import Books
    from xsdata.formats.dataclass.serializers import PycodeSerializer

    xml_string = Path("books.xml").read_text(encoding="utf-8")
    parser = XmlParser(context=XmlContext())
    books = parser.from_string(xml_string, Books)

    serializer = PycodeSerializer(context=XmlContext())
    pycode_string = serializer.render(books, var_name="books")

The generated output includes the imports needed to construct the object tree.

.. testsetup:: *

    from tests.fixtures.books.fixtures import books


Serialize python code to string
===============================

.. doctest::

    >>> from xsdata.formats.dataclass.serializers import PycodeSerializer
    >>> serializer = PycodeSerializer()
    >>> print(serializer.render(books, var_name="books"))
    from tests.fixtures.books.books import BookForm
    from tests.fixtures.books.books import Books
    from xsdata.models.datatype import XmlDate
    <BLANKLINE>
    <BLANKLINE>
    books = Books(
        book=[
            BookForm(
                author='Hightower, Kim',
                title='The First Book',
                genre='Fiction',
                price=44.95,
                pub_date=XmlDate(2000, 10, 1),
                review='An amazing story of nothing.',
                id='bk001'
            ),
            BookForm(
                author='Nagata, Suanne',
                title='Becoming Somebody',
                genre='Biography',
                price=33.95,
                pub_date=XmlDate(2001, 1, 10),
                review='A masterpiece of the fine art of gossiping.',
                id='bk002'
            ),
        ]
    )


Serialize python code to stream
================================

.. doctest::

    >>> from io import StringIO
    >>> output = StringIO()
    >>> serializer.write(output, books, "books")
    >>> output.getvalue() == serializer.render(books, "books")
    True


.. meta::
    :keywords: python, code, pycode, serialize, fixtures
