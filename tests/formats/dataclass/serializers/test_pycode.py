from io import StringIO
from unittest import TestCase
from xml.etree.ElementTree import QName

from tests.fixtures.books import BookForm
from tests.fixtures.books import Books
from tests.fixtures.books.fixtures import books
from tests.fixtures.models import Parent
from xsdata.formats.dataclass.serializers import PycodeSerializer
from xsdata.formats.dataclass.serializers.config import SerializerConfig
from xsdata.models.enums import Namespace


class PycodeSerializerTests(TestCase):
    def setUp(self):
        self.serializer = PycodeSerializer()

    def test_render(self):
        result = self.serializer.render(books, var_name="books")

        expected = (
            "from tests.fixtures.books.books import BookForm\n"
            "from tests.fixtures.books.books import Books\n"
            "from xsdata.models.datatype import XmlDate\n"
            "\n"
            "\n"
            "books = Books(\n"
            "    book=[\n"
            "        BookForm(\n"
            "            author='Hightower, Kim',\n"
            "            title='The First Book',\n"
            "            genre='Fiction',\n"
            "            price=44.95,\n"
            "            pub_date=XmlDate(2000, 10, 1),\n"
            "            review='An amazing story of nothing.',\n"
            "            id='bk001'\n"
            "        ),\n"
            "        BookForm(\n"
            "            author='Nagata, Suanne',\n"
            "            title='Becoming Somebody',\n"
            "            genre='Biography',\n"
            "            price=33.95,\n"
            "            pub_date=XmlDate(2001, 1, 10),\n"
            "            review='A masterpiece of the fine art of gossiping.',\n"
            "            id='bk002'\n"
            "        ),\n"
            "    ]\n"
            ")\n"
        )

        self.assertEqual(expected, result)

    def test_write(self):
        out = StringIO()
        self.serializer.write(out, books, "books")
        self.assertEqual(out.getvalue(), self.serializer.render(books, "books"))

    def test_accepts_serializer_config(self):
        serializer = PycodeSerializer(
            context=self.serializer.context,
            config=SerializerConfig(),
        )
        self.assertEqual(
            serializer.render(books, "books"),
            self.serializer.render(books, "books"),
        )

    def test_write_class_with_default_values(self):
        obj = Books(book=[BookForm(author="me")])

        result = self.serializer.render(obj, var_name="books")

        expected = (
            "from tests.fixtures.books.books import BookForm\n"
            "from tests.fixtures.books.books import Books\n"
            "\n"
            "\n"
            "books = Books(\n"
            "    book=[\n"
            "        BookForm(\n"
            "            author='me'\n"
            "        ),\n"
            "    ]\n"
            ")\n"
        )
        self.assertEqual(expected, result)

    def test_repr_empty_collections(self):
        types = set()

        self.assertEqual("[]", "".join(self.serializer.repr_object([], 0, types)))
        self.assertEqual("()", "".join(self.serializer.repr_object((), 0, types)))
        self.assertEqual("{}", "".join(self.serializer.repr_object({}, 0, types)))

    def test_repr_mapping(self):
        result = "".join(self.serializer.repr_object({"foo": "bar"}, 0, set()))
        self.assertEqual("{\n    'foo': 'bar',\n}", result)

    def test_repr_enum(self):
        result = "".join(self.serializer.repr_object(Namespace.XML, 0, set()))
        self.assertEqual("Namespace.XML", result)

    def test_repr_qname_and_non_finite_float(self):
        result = "".join(
            self.serializer.repr_object(QName("{ns}name"), 0, set())
        )
        self.assertEqual('QName("{ns}name")', result)

        result = "".join(self.serializer.repr_object(float("nan"), 0, set()))
        self.assertEqual('float("nan")', result)

    def test_build_imports_with_nested_types(self):
        expected = "from tests.fixtures.models import Parent\n"
        self.assertEqual(expected, self.serializer.build_imports({Parent.Inner}))
