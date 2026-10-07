/*
========================================================================
SchemaCrawler
http://www.schemacrawler.com
Copyright (c) 2000-2020, Sualeh Fatehi <sualeh@hotmail.com>.
All rights reserved.
------------------------------------------------------------------------

SchemaCrawler is distributed in the hope that it will be useful, but
WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.

SchemaCrawler and the accompanying materials are made available under
the terms of the Eclipse Public License v1.0, GNU General Public License
v3 or GNU Lesser General Public License v3.

You may elect to redistribute this code under any of these licenses.

The Eclipse Public License is available at:
http://www.eclipse.org/legal/epl-v10.html

The GNU General Public License v3 and the GNU Lesser General Public
License v3 are available at:
http://www.gnu.org/licenses/

========================================================================
*/

package schemacrawler.schemacrawler;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.emptyCollectionOf;
import static org.hamcrest.Matchers.is;

import java.util.ArrayList;

import org.junit.jupiter.api.Test;

import schemacrawler.schema.RoutineType;

public class LimitOptionsBuilderTest
{

  @Test
  public void emptyRoutineTypes()
  {
    final LimitOptions options = LimitOptionsBuilder.builder()
      .routineTypes(new ArrayList<RoutineType>())
      .toOptions();

    assertThat(options.getRoutineTypes(),
               is(emptyCollectionOf(RoutineType.class)));
  }

}
