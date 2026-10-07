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

import org.hipparchus.exception.MathIllegalArgumentException;
import org.hipparchus.util.Binary64;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;

/**
 * Test case for the field Lagrange form of a polynomial function.
 */
final class FieldPolynomialFunctionLagrangeFormTest {

    @Test
    void testLinearFunction() {
        final Binary64[] x = d(0.0, 3.0);
        final Binary64[] y = d(-4.0, 0.5);
        final FieldPolynomialFunctionLagrangeForm<Binary64> p =
                new FieldPolynomialFunctionLagrangeForm<>(x, y);

        assertEquals(-1.0, p.value(2.0).getReal(), 1.0e-12);
        assertEquals(2.75, p.value(new Binary64(4.5)).getReal(), 1.0e-12);
        assertEquals(5.0, p.value(6.0).getReal(), 1.0e-12);
        assertEquals(1, p.degree());

        final Binary64[] c = p.getCoefficients();
        assertEquals(2, c.length);
        assertEquals(-4.0, c[0].getReal(), 1.0e-12);
        assertEquals(1.5, c[1].getReal(), 1.0e-12);
    }

    @Test
    void testQuadraticFunction() {
        final Binary64[] x = d(0.0, -1.0, 0.5);
        final Binary64[] y = d(-3.0, -6.0, 0.0);
        final FieldPolynomialFunctionLagrangeForm<Binary64> p =
                new FieldPolynomialFunctionLagrangeForm<>(x, y);

        assertEquals(4.0, p.value(1.0).getReal(), 1.0e-12);
        assertEquals(22.0, p.value(2.5).getReal(), 1.0e-12);
        assertEquals(-5.0, p.value(-2.0).getReal(), 1.0e-12);
        assertEquals(2, p.degree());

        final Binary64[] c = p.getCoefficients();
        assertEquals(-3.0, c[0].getReal(), 1.0e-12);
        assertEquals(5.0, c[1].getReal(), 1.0e-12);
        assertEquals(2.0, c[2].getReal(), 1.0e-12);
    }

    @Test
    void testQuinticFunction() {
        final Binary64[] x = d(1.0, -1.0, 2.0, 3.0, -3.0, 0.5);
        final Binary64[] y = d(0.0, 0.0, -24.0, 0.0, -144.0, 2.34375);
        final FieldPolynomialFunctionLagrangeForm<Binary64> p =
                new FieldPolynomialFunctionLagrangeForm<>(x, y);

        assertEquals(0.0, p.value(0.0).getReal(), 1.0e-12);
        assertEquals(0.0, p.value(-2.0).getReal(), 1.0e-12);
        assertEquals(360.0, p.value(4.0).getReal(), 1.0e-12);
        assertEquals(5, p.degree());

        final Binary64[] c = p.getCoefficients();
        assertEquals(0.0, c[0].getReal(), 1.0e-12);
        assertEquals(6.0, c[1].getReal(), 1.0e-12);
        assertEquals(1.0, c[2].getReal(), 1.0e-12);
        assertEquals(-7.0, c[3].getReal(), 1.0e-12);
        assertEquals(-1.0, c[4].getReal(), 1.0e-12);
        assertEquals(1.0, c[5].getReal(), 1.0e-12);
    }

    @Test
    void testStaticEvaluateWithUnsortedPoints() {
        final Binary64[] x = d(2.0, -1.0, 0.0);
        final Binary64[] y = d(15.0, -6.0, -3.0);

        assertEquals(4.0, FieldPolynomialFunctionLagrangeForm.evaluate(x, y, new Binary64(1.0)).getReal(),
                            1.0e-12);
    }

    @Test
    void testCopiesInterpolationArrays() {
        final Binary64[] x = d(0.0, 3.0);
        final Binary64[] y = d(-4.0, 0.5);
        final FieldPolynomialFunctionLagrangeForm<Binary64> p =
                new FieldPolynomialFunctionLagrangeForm<>(x, y);

        x[0] = new Binary64(10.0);
        y[0] = new Binary64(10.0);

        assertEquals(0.0, p.getInterpolatingPoints()[0].getReal());
        assertEquals(-4.0, p.getInterpolatingValues()[0].getReal());
    }

    @Test
    void testTooFewPoints() {
        assertThrows(MathIllegalArgumentException.class,
                     () -> new FieldPolynomialFunctionLagrangeForm<>(d(1.0), d(2.0)));
    }

    @Test
    void testMismatchedArrays() {
        assertThrows(MathIllegalArgumentException.class,
                     () -> new FieldPolynomialFunctionLagrangeForm<>(d(1.0, 2.0, 3.0, 4.0),
                                                                     d(0.0, -4.0, -24.0)));
    }

    @Test
    void testDuplicateAbscissae() {
        assertThrows(MathIllegalArgumentException.class,
                     () -> new FieldPolynomialFunctionLagrangeForm<>(d(1.0, 1.0), d(0.0, 1.0)));
    }

    private static Binary64[] d(final double... values) {
        final Binary64[] result = new Binary64[values.length];
        for (int i = 0; i < values.length; i++) {
            result[i] = new Binary64(values[i]);
        }
        return result;
    }
}
