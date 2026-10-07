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

import java.util.Arrays;
import java.util.Comparator;

import org.hipparchus.CalculusFieldElement;
import org.hipparchus.Field;
import org.hipparchus.analysis.CalculusFieldUnivariateFunction;
import org.hipparchus.exception.LocalizedCoreFormats;
import org.hipparchus.exception.MathIllegalArgumentException;
import org.hipparchus.util.MathArrays;

/**
 * Implements the representation of a polynomial function in
 * <a href="http://mathworld.wolfram.com/LagrangeInterpolatingPolynomial.html">
 * Lagrange Form</a>. For reference, see <b>Introduction to Numerical
 * Analysis</b>, ISBN 038795452X, chapter 2.
 * <p>
 * The approximated function should be smooth enough for Lagrange polynomial
 * to work well. Otherwise, consider using splines instead.</p>
 *
 * @param <T> the type of the field elements
 * @since 4.0
 */
public class FieldPolynomialFunctionLagrangeForm<T extends CalculusFieldElement<T>>
    implements CalculusFieldUnivariateFunction<T> {

    /** Interpolating points (abscissas). */
    private final T[] x;

    /** Function values at interpolating points. */
    private final T[] y;

    /** Coefficients of the polynomial. */
    private T[] coefficients;

    /** Whether the polynomial coefficients are available. */
    private boolean coefficientsComputed;

    /**
     * Construct a Lagrange polynomial with the given abscissas and function
     * values. The order of interpolating points is not important.
     * <p>
     * The constructor makes copies of the input arrays.</p>
     *
     * @param x interpolating points
     * @param y function values at interpolating points
     * @throws MathIllegalArgumentException if the array lengths are different
     * @throws MathIllegalArgumentException if the number of points is less than 2
     * @throws MathIllegalArgumentException if two abscissae have the same value
     */
    public FieldPolynomialFunctionLagrangeForm(final T[] x, final T[] y)
        throws MathIllegalArgumentException {
        this.x = x.clone();
        this.y = y.clone();
        coefficientsComputed = false;

        if (!verifyInterpolationArray(this.x, this.y, false)) {
            sortInPlace(this.x, this.y);
            verifyInterpolationArray(this.x, this.y, true);
        }
    }

    /** Get the {@link Field} to which the instance belongs.
     * @return {@link Field} to which the instance belongs
     */
    public Field<T> getField() {
        return x[0].getField();
    }

    /**
     * Calculate the function value at the given point.
     *
     * @param z point at which the function value is to be computed
     * @return the function value
     */
    public T value(final double z) {
        return value(getField().getZero().add(z));
    }

    /**
     * Calculate the function value at the given point.
     *
     * @param z point at which the function value is to be computed
     * @return the function value
     */
    @Override
    public T value(final T z) {
        return evaluateInternal(x, y, z);
    }

    /**
     * Returns the degree of the polynomial.
     *
     * @return the degree of the polynomial
     */
    public int degree() {
        return x.length - 1;
    }

    /**
     * Returns a copy of the interpolating points array.
     *
     * @return a fresh copy of the interpolating points array
     */
    public T[] getInterpolatingPoints() {
        return x.clone();
    }

    /**
     * Returns a copy of the interpolating values array.
     *
     * @return a fresh copy of the interpolating values array
     */
    public T[] getInterpolatingValues() {
        return y.clone();
    }

    /**
     * Returns a copy of the coefficients array.
     * <p>
     * Note that coefficients computation can be ill-conditioned. Use with caution
     * and only when it is necessary.</p>
     *
     * @return a fresh copy of the coefficients array
     */
    public T[] getCoefficients() {
        if (!coefficientsComputed) {
            computeCoefficients();
        }
        return coefficients.clone();
    }

    /**
     * Evaluate the Lagrange polynomial using
     * <a href="http://mathworld.wolfram.com/NevillesAlgorithm.html">
     * Neville's Algorithm</a>. It takes O(n^2) time.
     *
     * @param x interpolating points array
     * @param y interpolating values array
     * @param z point at which the function value is to be computed
     * @param <T> the type of the field elements
     * @return the function value
     * @throws MathIllegalArgumentException if the array lengths are different
     * @throws MathIllegalArgumentException if the number of points is less than 2
     * @throws MathIllegalArgumentException if two abscissae have the same value
     */
    public static <T extends CalculusFieldElement<T>> T evaluate(final T[] x, final T[] y, final T z)
        throws MathIllegalArgumentException {
        if (verifyInterpolationArray(x, y, false)) {
            return evaluateInternal(x, y, z);
        }

        final T[] sortedX = x.clone();
        final T[] sortedY = y.clone();
        sortInPlace(sortedX, sortedY);
        verifyInterpolationArray(sortedX, sortedY, true);
        return evaluateInternal(sortedX, sortedY, z);
    }

    /**
     * Evaluate the Lagrange polynomial using Neville's Algorithm.
     *
     * @param x interpolating points array
     * @param y interpolating values array
     * @param z point at which the function value is to be computed
     * @param <T> the type of the field elements
     * @return the function value
     */
    private static <T extends CalculusFieldElement<T>> T evaluateInternal(final T[] x, final T[] y, final T z) {
        int nearest = 0;
        final int n = x.length;
        final T[] c = MathArrays.buildArray(x[0].getField(), n);
        final T[] d = MathArrays.buildArray(x[0].getField(), n);
        double minDistance = Double.POSITIVE_INFINITY;
        for (int i = 0; i < n; i++) {
            c[i] = y[i];
            d[i] = y[i];

            final double distance = z.subtract(x[i]).norm();
            if (distance < minDistance) {
                nearest = i;
                minDistance = distance;
            }
        }

        T value = y[nearest];

        for (int i = 1; i < n; i++) {
            for (int j = 0; j < n - i; j++) {
                final T tc = x[j].subtract(z);
                final T td = x[i + j].subtract(z);
                final T divider = x[j].subtract(x[i + j]);
                final T w = c[j + 1].subtract(d[j]).divide(divider);
                c[j] = tc.multiply(w);
                d[j] = td.multiply(w);
            }
            if (nearest < 0.5 * (n - i + 1)) {
                value = value.add(c[nearest]);
            } else {
                nearest--;
                value = value.add(d[nearest]);
            }
        }

        return value;
    }

    /**
     * Calculate the coefficients of Lagrange polynomial from the interpolation data.
     */
    protected void computeCoefficients() {
        final int n = degree() + 1;
        final Field<T> field = getField();
        coefficients = MathArrays.buildArray(field, n);
        Arrays.fill(coefficients, field.getZero());

        final T[] c = MathArrays.buildArray(field, n + 1);
        Arrays.fill(c, field.getZero());
        c[0] = field.getOne();
        for (int i = 0; i < n; i++) {
            for (int j = i; j > 0; j--) {
                c[j] = c[j - 1].subtract(c[j].multiply(x[i]));
            }
            c[0] = c[0].multiply(x[i].negate());
            c[i + 1] = field.getOne();
        }

        final T[] tc = MathArrays.buildArray(field, n);
        for (int i = 0; i < n; i++) {
            T d = field.getOne();
            for (int j = 0; j < n; j++) {
                if (i != j) {
                    d = d.multiply(x[i].subtract(x[j]));
                }
            }
            final T t = y[i].divide(d);
            tc[n - 1] = c[n];
            coefficients[n - 1] = coefficients[n - 1].add(t.multiply(tc[n - 1]));
            for (int j = n - 2; j >= 0; j--) {
                tc[j] = c[j + 1].add(tc[j + 1].multiply(x[i]));
                coefficients[j] = coefficients[j].add(t.multiply(tc[j]));
            }
        }

        coefficientsComputed = true;
    }

    /**
     * Check that the interpolation arrays are valid.
     *
     * @param x interpolating points array
     * @param y interpolating values array
     * @param abort whether to throw an exception if {@code x} is not sorted
     * @param <T> the type of the field elements
     * @return {@code false} if {@code x} is not sorted in increasing order,
     * {@code true} otherwise
     * @throws MathIllegalArgumentException if the array lengths are different
     * @throws MathIllegalArgumentException if the number of points is less than 2
     * @throws MathIllegalArgumentException if {@code x} is not sorted in strictly
     * increasing order and {@code abort} is {@code true}
     */
    public static <T extends CalculusFieldElement<T>> boolean verifyInterpolationArray(final T[] x,
                                                                                        final T[] y,
                                                                                        final boolean abort)
        throws MathIllegalArgumentException {
        MathArrays.checkEqualLength(x, y);
        if (x.length < 2) {
            throw new MathIllegalArgumentException(LocalizedCoreFormats.WRONG_NUMBER_OF_POINTS, 2, x.length, true);
        }

        return MathArrays.checkOrder(x, MathArrays.OrderDirection.INCREASING, true, abort);
    }

    /** Sort interpolating arrays by the real value of their abscissas. */
    private static <T extends CalculusFieldElement<T>> void sortInPlace(final T[] x, final T[] y) {
        final Integer[] indices = new Integer[x.length];
        for (int i = 0; i < indices.length; i++) {
            indices[i] = i;
        }
        Arrays.sort(indices, Comparator.comparingDouble(index -> x[index].getReal()));

        final T[] sortedX = x.clone();
        final T[] sortedY = y.clone();
        for (int i = 0; i < indices.length; i++) {
            x[i] = sortedX[indices[i]];
            y[i] = sortedY[indices[i]];
        }
    }
}
