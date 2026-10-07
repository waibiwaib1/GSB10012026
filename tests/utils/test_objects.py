from types import SimpleNamespace
from unittest import TestCase
from xml.etree.ElementTree import QName

from xsdata.utils import objects


class ObjectsTests(TestCase):
    def test_update(self):

        obj = SimpleNamespace()
        obj.foo = SimpleNamespace()
        obj.foo.bar = 2
        obj.bar = 1

        kwargs = {"foo.bar": 1, "bar": 2}
        objects.update(obj, **kwargs)
        self.assertEqual(1, obj.foo.bar)
        self.assertEqual(2, obj.bar)

    def test_literal_value(self):
        self.assertEqual("1", objects.literal_value(1))
        self.assertEqual("'foo'", objects.literal_value("foo"))
        self.assertEqual(
            'QName("{ns}name")',
            objects.literal_value(QName("{ns}name")),
        )
        self.assertEqual('float("nan")', objects.literal_value(float("nan")))
        self.assertEqual('float("inf")', objects.literal_value(float("inf")))
