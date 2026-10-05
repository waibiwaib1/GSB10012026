from mpmath import inf, mp, mpf, nan, ninf


def test_mpf_format_basic():
    value = mpf('1.0')

    assert format(value) == '1.0'
    assert '{:10.0f}'.format(value) == '         1'
    assert format(value, '.2f') == '1.00'
    assert format(value, '08.2f') == '00001.00'
    assert format(-value, '+.2f') == '-1.00'
    assert format(value, ' .2f') == ' 1.00'
    assert format(value, '^12.3f') == '   1.000    '


def test_mpf_format_rounding_and_precision():
    assert format(mpf('0.75'), '.0f') == '1'
    assert format(mpf('2.5'), '.0f') == '2'
    assert format(mpf('3.5'), '.0f') == '4'
    assert format(mpf('1.5'), '.2e') == '1.50e+0'
    assert format(mpf('0.0015'), '.3g') == '0.00150'
    assert format(mpf('0.0015'), '.2%') == '0.15%'


def test_mpf_format_arbitrary_precision():
    mp.dps = 50
    try:
        assert format(mp.pi, '.50f') == (
            '3.14159265358979323846264338327950288419716939937511')
        assert format(mpf('1.22'), '.25f') == '1.2200000000000000000000000'
    finally:
        mp.dps = 15


def test_mpf_format_special_values():
    assert format(inf, 'f') == 'inf'
    assert format(inf, 'F') == 'INF'
    assert format(inf, '+f') == '+inf'
    assert format(ninf, 'f') == '-inf'
    assert format(nan, 'f') == 'nan'
    assert format(nan, 'F') == 'NAN'
