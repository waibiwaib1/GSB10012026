/*
 * Licensed to the Hipparchus project under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The Hipparchus project licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package org.hipparchus.analysis.polynomials;

import org.hipparchus.complex.Complex;
import org.hipparchus.exception.MathIllegalArgumentException;
import org.hipparchus.util.Binary64;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.fail;

/**
 * Test case for the field version of the Lagrange form of polynomial
 * function.
 * <p>
 * We use n+1 points to interpolate a polynomial of degree n. This should
 * give us the exact same polynomial as result. Thus we can use a very
 * small tolerance to account only for round-off errors.
 */
final class FieldPolynomialFunctionLagrangeFormTest {

    /**
     * Test of polynomial for the linear function.
     */
    @Test
    void testLinearFunction() {
        doTestLinearFunction(build64(new double[] {0.0, 3.0}),
                             build64(new double[] {-4.0, 0.5}));
        doTestLinearFunction(new Complex[] {Complex.ZERO, new Complex(3.0, 0.0)},
                             new Complex[] {new Complex(-4.0, 0.0), new Complex(0.5, 0.0)});
    }

    private <T extends org.hipparchus.CalculusFieldElement<T>> void doTestLinearFunction(final T[] x, final T[] y) {
        final FieldPolynomialFunctionLagrangeForm<T> p = new FieldPolynomialFunctionLagrangeForm<>(x, y);
        final double tolerance = 1E-12;

        // p(x) = 1.5x - 4
        assertEquals(-1.0, p.value(2.0).getReal(), tolerance);
        assertEquals(2.75, p.value(4.5).getReal(), tolerance);
        assertEquals(5.0, p.value(x[0].getField().getZero().add(6.0)).getReal(), tolerance);

        assertEquals(1, p.degree());

        final T[] c = p.getCoefficients();
        assertEquals(2, c.length);
        assertEquals(-4.0, c[0].getReal(), tolerance);
        assertEquals(1.5, c[1].getReal(), tolerance);
    }

    /**
     * Test of polynomial for the quadratic function.
     */
    @Test
    void testQuadraticFunction() {
        final FieldPolynomialFunctionLagrangeForm<Binary64> p =
                new FieldPolynomialFunctionLagrangeForm<>(build64(new double[] {0.0, -1.0, 0.5}),
                                                          build64(new double[] {-3.0, -6.0, 0.0}));
        final double tolerance = 1E-12;

        // p(x) = 2x^2 + 5x - 3 = (2x - 1)(x + 3)
        assertEquals(4.0, p.value(1.0).getReal(), tolerance);
        assertEquals(22.0, p.value(2.5).getReal(), tolerance);
        assertEquals(-5.0, p.value(-2.0).getReal(), tolerance);

        assertEquals(2, p.degree());

        final Binary64[] c = p.getCoefficients();
        assertEquals(3, c.length);
        assertEquals(-3.0, c[0].getReal(), tolerance);
        assertEquals(5.0, c[1].getReal(), tolerance);
        assertEquals(2.0, c[2].getReal(), tolerance);
    }

    /**
     * Test of polynomial for the quintic function.
     */
    @Test
    void testQuinticFunction() {
        final FieldPolynomialFunctionLagrangeForm<Binary64> p =
                new FieldPolynomialFunctionLagrangeForm<>(
                        build64(new double[] {1.0, -1.0, 2.0, 3.0, -3.0, 0.5}),
                        build64(new double[] {0.0, 0.0, -24.0, 0.0, -144.0, 2.34375}));
        final double tolerance = 1E-12;

        // p(x) = x^5 - x^4 - 7x^3 + x^2 + 6x = x(x^2 - 1)(x + 2)(x - 3)
        assertEquals(0.0, p.value(0.0).getReal(), tolerance);
        assertEquals(0.0, p.value(-2.0).getReal(), tolerance);
        assertEquals(360.0, p.value(4.0).getReal(), tolerance);

        assertEquals(5, p.degree());

        final Binary64[] c = p.getCoefficients();
        assertEquals(6, c.length);
        assertEquals(0.0, c[0].getReal(), tolerance);
        assertEquals(6.0, c[1].getReal(), tolerance);
        assertEquals(1.0, c[2].getReal(), tolerance);
        assertEquals(-7.0, c[3].getReal(), tolerance);
        assertEquals(-1.0, c[4].getReal(), tolerance);
        assertEquals(1.0, c[5].getReal(), tolerance);
    }

    /**
     * Test the static evaluate method, including with unsorted points.
     */
    @Test
    void testEvaluate() {
        final Binary64[] x = build64(new double[] {0.0, -1.0, 0.5});
        final Binary64[] y = build64(new double[] {-3.0, -6.0, 0.0});
        final double tolerance = 1E-12;

        assertEquals(4.0,
                     FieldPolynomialFunctionLagrangeForm.evaluate(x, y, new Binary64(1.0)).getReal(),
                     tolerance);
        assertEquals(22.0,
                     FieldPolynomialFunctionLagrangeForm.evaluate(x, y, new Binary64(2.5)).getReal(),
                     tolerance);
    }

    /**
     * Test of parameters for the polynomial.
     */
    @Test
    void testParameters() {

        try {
            // bad input array length
            final Binary64[] x = { new Binary64(1.0) };
            final Binary64[] y = { new Binary64(2.0) };
            new FieldPolynomialFunctionLagrangeForm<>(x, y);
            fail("Expecting MathIllegalArgumentException - bad input array length");
        } catch (MathIllegalArgumentException ex) {
            // expected
        }
        try {
            // mismatch input arrays
            final Binary64[] x = build64(new double[] {1.0, 2.0, 3.0, 4.0});
            final Binary64[] y = build64(new double[] {0.0, -4.0, -24.0});
            new FieldPolynomialFunctionLagrangeForm<>(x, y);
            fail("Expecting MathIllegalArgumentException - mismatch input arrays");
        } catch (MathIllegalArgumentException ex) {
            // expected
        }
        try {
            // duplicated abscissa
            final Binary64[] x = build64(new double[] {1.0, 2.0, 2.0});
            final Binary64[] y = build64(new double[] {0.0, -4.0, -24.0});
            new FieldPolynomialFunctionLagrangeForm<>(x, y);
            fail("Expecting MathIllegalArgumentException - duplicated abscissae");
        } catch (MathIllegalArgumentException ex) {
            // expected
        }
    }

    /** Build an array of {@link Binary64} from an array of doubles. */
    private static Binary64[] build64(final double[] values) {
        final Binary64[] result = new Binary64[values.length];
        for (int i = 0; i < values.length; i++) {
            result[i] = new Binary64(values[i]);
        }
        return result;
    }
}
