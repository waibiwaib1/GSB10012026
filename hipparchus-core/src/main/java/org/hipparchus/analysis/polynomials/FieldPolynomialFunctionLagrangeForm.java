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

import org.hipparchus.CalculusFieldElement;
import org.hipparchus.Field;
import org.hipparchus.analysis.CalculusFieldUnivariateFunction;
import org.hipparchus.exception.LocalizedCoreFormats;
import org.hipparchus.exception.MathIllegalArgumentException;
import org.hipparchus.util.FastMath;
import org.hipparchus.util.MathArrays;

/**
 * Implements the representation of a polynomial function in
 * <a href="http://mathworld.wolfram.com/LagrangeInterpolatingPolynomial.html">
 * Lagrange Form</a> for any field type. For reference, see <b>Introduction to
 * Numerical Analysis</b>, ISBN 038795452X, chapter 2.
 * <p>
 * The approximated function should be smooth enough for Lagrange polynomial
 * to work well. Otherwise, consider using splines instead.</p>
 * @param <T> the type of the field elements
 * @since 4.0
 */
public class FieldPolynomialFunctionLagrangeForm<T extends CalculusFieldElement<T>>
    implements CalculusFieldUnivariateFunction<T> {

    /**
     * The coefficients of the polynomial, ordered by degree -- i.e.
     * coefficients[0] is the constant term and coefficients[n] is the
     * coefficient of x^n where n is the degree of the polynomial.
     */
    private T[] coefficients;

    /** Interpolating points (abscissas). */
    private final T[] x;

    /** Function values at interpolating points. */
    private final T[] y;

    /** Whether the polynomial coefficients are available. */
    private boolean coefficientsComputed;

    /**
     * Construct a Lagrange polynomial with the given abscissas and function
     * values. The order of interpolating points is not important.
     * <p>
     * The constructor makes copy of the input arrays and assigns them.</p>
     *
     * @param x interpolating points
     * @param y function values at interpolating points
     * @throws MathIllegalArgumentException if the array lengths are different.
     * @throws MathIllegalArgumentException if the number of points is less than 2.
     * @throws MathIllegalArgumentException if two abscissae have the same value.
     */
    public FieldPolynomialFunctionLagrangeForm(final T[] x, final T[] y)
        throws MathIllegalArgumentException {
        this.x = MathArrays.buildArray(x[0].getField(), x.length);
        this.y = MathArrays.buildArray(y[0].getField(), y.length);
        System.arraycopy(x, 0, this.x, 0, x.length);
        System.arraycopy(y, 0, this.y, 0, y.length);
        coefficientsComputed = false;

        if (!verifyInterpolationArray(x, y, false)) {
            sortInPlace(this.x, this.y);
            // Second check in case some abscissa is duplicated.
            verifyInterpolationArray(this.x, this.y, true);
        }
    }

    /**
     * Calculate the function value at the given point.
     *
     * @param z Point at which the function value is to be computed.
     * @return the function value.
     */
    public T value(final double z) {
        return value(x[0].getField().getZero().add(z));
    }

    /**
     * Calculate the function value at the given point.
     *
     * @param z Point at which the function value is to be computed.
     * @return the function value.
     */
    @Override
    public T value(final T z) {
        return evaluateInternal(x, y, z);
    }

    /** Get the {@link Field} to which the instance belongs.
     * @return {@link Field} to which the instance belongs
     */
    public Field<T> getField() {
        return x[0].getField();
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
     * <p>
     * Changes made to the returned copy will not affect the polynomial.</p>
     *
     * @return a fresh copy of the interpolating points array
     */
    public T[] getInterpolatingPoints() {
        final T[] out = MathArrays.buildArray(getField(), x.length);
        System.arraycopy(x, 0, out, 0, x.length);
        return out;
    }

    /**
     * Returns a copy of the interpolating values array.
     * <p>
     * Changes made to the returned copy will not affect the polynomial.</p>
     *
     * @return a fresh copy of the interpolating values array
     */
    public T[] getInterpolatingValues() {
        final T[] out = MathArrays.buildArray(getField(), y.length);
        System.arraycopy(y, 0, out, 0, y.length);
        return out;
    }

    /**
     * Returns a copy of the coefficients array.
     * <p>
     * Changes made to the returned copy will not affect the polynomial.</p>
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
        final T[] out = MathArrays.buildArray(getField(), coefficients.length);
        System.arraycopy(coefficients, 0, out, 0, coefficients.length);
        return out;
    }

    /**
     * Evaluate the Lagrange polynomial using
     * <a href="http://mathworld.wolfram.com/NevillesAlgorithm.html">
     * Neville's Algorithm</a>. It takes O(n^2) time.
     *
     * @param x Interpolating points array.
     * @param y Interpolating values array.
     * @param z Point at which the function value is to be computed.
     * @param <T> the type of the field elements
     * @return the function value.
     * @throws MathIllegalArgumentException if {@code x} and {@code y} have
     * different lengths.
     * @throws MathIllegalArgumentException if {@code x} is not sorted in strictly
     * increasing order.
     * @throws MathIllegalArgumentException if the size of {@code x} is less
     * than 2.
     */
    public static <T extends CalculusFieldElement<T>> T evaluate(final T[] x, final T[] y, final T z)
        throws MathIllegalArgumentException {
        if (verifyInterpolationArray(x, y, false)) {
            return evaluateInternal(x, y, z);
        }

        // Array is not sorted.
        final T[] xNew = MathArrays.buildArray(x[0].getField(), x.length);
        final T[] yNew = MathArrays.buildArray(y[0].getField(), y.length);
        System.arraycopy(x, 0, xNew, 0, x.length);
        System.arraycopy(y, 0, yNew, 0, y.length);

        sortInPlace(xNew, yNew);
        // Second check in case some abscissa is duplicated.
        verifyInterpolationArray(xNew, yNew, true);
        return evaluateInternal(xNew, yNew, z);
    }

    /**
     * Evaluate the Lagrange polynomial using
     * <a href="http://mathworld.wolfram.com/NevillesAlgorithm.html">
     * Neville's Algorithm</a>. It takes O(n^2) time.
     *
     * @param x Interpolating points array.
     * @param y Interpolating values array.
     * @param z Point at which the function value is to be computed.
     * @param <T> the type of the field elements
     * @return the function value.
     */
    private static <T extends CalculusFieldElement<T>> T evaluateInternal(final T[] x, final T[] y, final T z) {
        int nearest = 0;
        final int n = x.length;
        final Field<T> field = z.getField();
        final T[] c = MathArrays.buildArray(field, n);
        final T[] d = MathArrays.buildArray(field, n);
        double minDist = Double.POSITIVE_INFINITY;
        for (int i = 0; i < n; i++) {
            // initialize the difference arrays
            c[i] = y[i];
            d[i] = y[i];
            // find out the abscissa closest to z
            final double dist = FastMath.abs(z.subtract(x[i]).getReal());
            if (dist < minDist) {
                nearest = i;
                minDist = dist;
            }
        }

        // initial approximation to the function value at z
        T value = y[nearest];

        for (int i = 1; i < n; i++) {
            for (int j = 0; j < n - i; j++) {
                final T tc = x[j].subtract(z);
                final T td = x[i + j].subtract(z);
                final T divider = x[j].subtract(x[i + j]);
                // update the difference arrays
                final T w = c[j + 1].subtract(d[j]).divide(divider);
                c[j] = tc.multiply(w);
                d[j] = td.multiply(w);
            }
            // sum up the difference terms to get the final value
            if (nearest < 0.5 * (n - i + 1)) {
                value = value.add(c[nearest]);    // fork down
            } else {
                nearest--;
                value = value.add(d[nearest]);    // fork up
            }
        }

        return value;
    }

    /**
     * Calculate the coefficients of Lagrange polynomial from the
     * interpolation data. It takes O(n^2) time.
     * Note that this computation can be ill-conditioned: Use with caution
     * and only when it is necessary.
     */
    protected void computeCoefficients() {
        final int n = degree() + 1;
        final Field<T> field = getField();
        coefficients = MathArrays.buildArray(field, n);
        for (int i = 0; i < n; i++) {
            coefficients[i] = field.getZero();
        }

        // c[] are the coefficients of P(x) = (x-x[0])(x-x[1])...(x-x[n-1])
        final T[] c = MathArrays.buildArray(field, n + 1);
        c[0] = field.getOne();
        for (int i = 0; i < n; i++) {
            for (int j = i; j > 0; j--) {
                c[j] = c[j - 1].subtract(c[j].multiply(x[i]));
            }
            c[0] = c[0].multiply(x[i]).negate();
            c[i + 1] = field.getOne();
        }

        final T[] tc = MathArrays.buildArray(field, n);
        for (int i = 0; i < n; i++) {
            // d = (x[i]-x[0])...(x[i]-x[i-1])(x[i]-x[i+1])...(x[i]-x[n-1])
            T d = field.getOne();
            for (int j = 0; j < n; j++) {
                if (i != j) {
                    d = d.multiply(x[i].subtract(x[j]));
                }
            }
            final T t = y[i].divide(d);
            // Lagrange polynomial is the sum of n terms, each of which is a
            // polynomial of degree n-1. tc[] are the coefficients of the i-th
            // numerator Pi(x) = (x-x[0])...(x-x[i-1])(x-x[i+1])...(x-x[n-1]).
            tc[n - 1] = c[n];     // actually c[n] = 1
            coefficients[n - 1] = coefficients[n - 1].add(t.multiply(tc[n - 1]));
            for (int j = n - 2; j >= 0; j--) {
                tc[j] = c[j + 1].add(tc[j + 1].multiply(x[i]));
                coefficients[j] = coefficients[j].add(t.multiply(tc[j]));
            }
        }

        coefficientsComputed = true;
    }

    /**
     * Sort the interpolating points and values in place according to the
     * (strictly increasing) real value of the points.
     *
     * @param x interpolating points array
     * @param y interpolating values array
     * @param <T> the type of the field elements
     */
    private static <T extends CalculusFieldElement<T>> void sortInPlace(final T[] x, final T[] y) {
        final Integer[] indices = new Integer[x.length];
        for (int i = 0; i < indices.length; i++) {
            indices[i] = i;
        }
        Arrays.sort(indices, (i, j) -> Double.compare(x[i].getReal(), x[j].getReal()));
        final T[] xSorted = MathArrays.buildArray(x[0].getField(), x.length);
        final T[] ySorted = MathArrays.buildArray(y[0].getField(), y.length);
        System.arraycopy(x, 0, xSorted, 0, x.length);
        System.arraycopy(y, 0, ySorted, 0, y.length);
        for (int k = 0; k < indices.length; k++) {
            x[k] = xSorted[indices[k]];
            y[k] = ySorted[indices[k]];
        }
    }

    /**
     * Check that the interpolation arrays are valid.
     * The arrays features checked by this method are that both arrays have the
     * same length and this length is at least 2.
     *
     * @param x Interpolating points array.
     * @param y Interpolating values array.
     * @param abort Whether to throw an exception if {@code x} is not sorted.
     * @param <T> the type of the field elements
     * @throws MathIllegalArgumentException if the array lengths are different.
     * @throws MathIllegalArgumentException if the number of points is less than 2.
     * @throws MathIllegalArgumentException if {@code x} is not sorted in strictly
     * increasing order and {@code abort} is {@code true}.
     * @return {@code false} if the {@code x} is not sorted in increasing order,
     * {@code true} otherwise.
     * @see #evaluate(CalculusFieldElement[], CalculusFieldElement[], CalculusFieldElement)
     * @see #computeCoefficients()
     */
    public static <T extends CalculusFieldElement<T>> boolean verifyInterpolationArray(final T[] x,
                                                                                        final T[] y,
                                                                                        final boolean abort)
        throws MathIllegalArgumentException {
        MathArrays.checkEqualLength(x, y);
        if (x.length < 2) {
            throw new MathIllegalArgumentException(LocalizedCoreFormats.WRONG_NUMBER_OF_POINTS,
                                                   2, x.length, true);
        }

        return MathArrays.checkOrder(x, MathArrays.OrderDirection.INCREASING, true, abort);
    }
}
