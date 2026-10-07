package de.komoot.photon.nominatim;

import de.komoot.photon.ReflectionTestUtil;
import de.komoot.photon.Updater;
import de.komoot.photon.nominatim.testdb.H2DataAdapter;
import de.komoot.photon.nominatim.testdb.OsmlineTestRow;
import de.komoot.photon.nominatim.testdb.PlacexTestRow;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.jdbc.datasource.embedded.EmbeddedDatabase;
import org.springframework.jdbc.datasource.embedded.EmbeddedDatabaseBuilder;
import org.springframework.jdbc.datasource.embedded.EmbeddedDatabaseType;

import java.sql.Timestamp;

import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.never;
import static org.mockito.Mockito.times;
import static org.mockito.Mockito.verify;

class NominatimUpdaterDBTest {
    private EmbeddedDatabase db;
    private JdbcTemplate jdbc;
    private NominatimUpdater updater;
    private Updater elasticsearchUpdater;

    @BeforeEach
    void setup() {
        db = new EmbeddedDatabaseBuilder()
                .setType(EmbeddedDatabaseType.H2)
                .generateUniqueName(true)
                .addScript("/test-schema.sql")
                .build();

        jdbc = new JdbcTemplate(db);
        updater = new NominatimUpdater(null, 0, null, null, null, new H2DataAdapter());
        elasticsearchUpdater = mock(Updater.class);
        updater.setUpdater(elasticsearchUpdater);

        ReflectionTestUtil.setFieldValue(updater, "template", jdbc);
        ReflectionTestUtil.setFieldValue(updater, "exporter", "template", jdbc);
    }

    @Test
    void updateIndexedPlace() {
        PlacexTestRow place = new PlacexTestRow("amenity", "cafe").name("Spot").add(jdbc);
        addUpdate("placex", place.getPlaceId(), "UPDATE", 1);

        updater.update();

        verify(elasticsearchUpdater).delete(place.getPlaceId());
        verify(elasticsearchUpdater).create(any());
        verify(elasticsearchUpdater).finish();
        assertNoQueuedUpdates();
    }

    @Test
    void deletePlaceWhenNominatimDeletesIt() {
        long placeId = 99999L;
        addUpdate("placex", placeId, "DELETE", 1);

        updater.update();

        verify(elasticsearchUpdater).delete(placeId);
        verify(elasticsearchUpdater, never()).create(any());
        assertNoQueuedUpdates();
    }

    @Test
    void useNewestChangeWhenPlaceIsQueuedMoreThanOnce() {
        PlacexTestRow place = new PlacexTestRow("amenity", "cafe").name("Spot").add(jdbc);
        addUpdate("placex", place.getPlaceId(), "DELETE", 1);
        addUpdate("placex", place.getPlaceId(), "UPDATE", 2);

        updater.update();

        verify(elasticsearchUpdater).delete(place.getPlaceId());
        verify(elasticsearchUpdater).create(any());
        assertNoQueuedUpdates();
    }

    @Test
    void updateIndexedInterpolation() {
        PlacexTestRow street = PlacexTestRow.make_street("Main Street").add(jdbc);
        OsmlineTestRow osmline = new OsmlineTestRow().number(1, 5, "all").parent(street).add(jdbc);
        addUpdate("location_property_osmline", osmline.getPlaceId(), "UPDATE", 1);

        updater.update();

        verify(elasticsearchUpdater).delete(osmline.getPlaceId());
        verify(elasticsearchUpdater, times(3)).create(any());
        assertNoQueuedUpdates();
    }

    private void addUpdate(String relation, Long placeId, String operation, int seconds) {
        jdbc.update("INSERT INTO photon_updates (rel, place_id, operation, indexed_date) VALUES (?, ?, ?, ?)",
                relation, placeId, operation, new Timestamp(seconds * 1000L));
    }

    private void assertNoQueuedUpdates() {
        Integer count = jdbc.queryForObject("SELECT count(*) FROM photon_updates", Integer.class);
        org.junit.jupiter.api.Assertions.assertEquals(Integer.valueOf(0), count);
    }
}
