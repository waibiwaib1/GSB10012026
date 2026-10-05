/**
 * Copyright 2018 Crawler-Commons
 * 
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * 
 *     http://www.apache.org/licenses/LICENSE-2.0
 * 
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package crawlercommons.sitemaps;

import crawlercommons.sitemaps.extension.*;
import org.junit.jupiter.api.Test;

import java.io.IOException;
import java.net.MalformedURLException;
import java.net.URL;
import java.time.ZonedDateTime;
import java.util.HashMap;
import java.util.Iterator;

import static org.junit.jupiter.api.Assertions.*;

public class SiteMapParserExtensionTest {

    @Test
    public void testVideosSitemap() throws UnknownFormatException, IOException {
        SiteMapParser parser = new SiteMapParser();
        parser.enableExtension(Extension.VIDEO);

        String contentType = "text/xml";
        byte[] content = SiteMapParserTest.getResourceAsBytes("src/test/resources/sitemaps/extension/sitemap-videos.xml");

        URL url = new URL("http://www.example.com/sitemap-video.xml");
        AbstractSiteMap asm = parser.parseSiteMap(contentType, content, url);
        assertEquals(false, asm.isIndex());
        assertEquals(true, asm instanceof SiteMap);
        SiteMap sm = (SiteMap) asm;
        assertEquals(3, sm.getSiteMapUrls().size());
        Iterator<SiteMapURL> siter = sm.getSiteMapUrls().iterator();

        // first <loc> element: nearly all video attributes
        VideoAttributes expectedVideoAttributes = new VideoAttributes(new URL("http://www.example.com/thumbs/123.jpg"), "Grilling steaks for summer",
                        "Alkis shows you how to get perfectly done steaks every time", new URL("http://www.example.com/video123.flv"), new URL("http://www.example.com/videoplayer.swf?video=123"));
        expectedVideoAttributes.setDuration(600);
        ZonedDateTime dt = ZonedDateTime.parse("2009-11-05T19:20:30+08:00");
        expectedVideoAttributes.setExpirationDate(dt);
        dt = ZonedDateTime.parse("2007-11-05T19:20:30+08:00");
        expectedVideoAttributes.setPublicationDate(dt);
        expectedVideoAttributes.setRating(4.2f);
        expectedVideoAttributes.setViewCount(12345);
        expectedVideoAttributes.setFamilyFriendly(true);
        expectedVideoAttributes.setTags(new String[] { "sample_tag1", "sample_tag2" });
        expectedVideoAttributes.setAllowedCountries(new String[] { "IE", "GB", "US", "CA" });
        expectedVideoAttributes.setGalleryLoc(new URL("http://cooking.example.com"));
        expectedVideoAttributes.setGalleryTitle("Cooking Videos");
        expectedVideoAttributes.setPrices(new VideoAttributes.VideoPrice[] { new VideoAttributes.VideoPrice("EUR", 1.99f, VideoAttributes.VideoPriceType.own) });
        expectedVideoAttributes.setRequiresSubscription(true);
        expectedVideoAttributes.setUploader("GrillyMcGrillerson");
        expectedVideoAttributes.setUploaderInfo(new URL("http://www.example.com/users/grillymcgrillerson"));
        expectedVideoAttributes.setLive(false);
        VideoAttributes attr = (VideoAttributes) siter.next().getAttributesForExtension(Extension.VIDEO)[0];
        assertNotNull(attr);
        assertEquals(expectedVideoAttributes, attr);

        // locale-specific number format in <video:price>, test #220
        // The current expected behavior is to not handle non-US locale price
        // values and set the price value to null if parsing as float value
        // fails.
        expectedVideoAttributes = new VideoAttributes(new URL("http://www.example.com/thumbs/123-2.jpg"), "Grilling steaks for summer, episode 2",
                        "Alkis shows you how to get perfectly done steaks every time", new URL("http://www.example.com/video123-2.flv"), null);
        expectedVideoAttributes.setPrices(new VideoAttributes.VideoPrice[] { new VideoAttributes.VideoPrice("EUR", null, VideoAttributes.VideoPriceType.own) });
        attr = (VideoAttributes) siter.next().getAttributesForExtension(Extension.VIDEO)[0];
        assertNotNull(attr);
        assertEquals(expectedVideoAttributes, attr);

        // empty price, only type (purchase or rent) is indicated, see #221
        expectedVideoAttributes = new VideoAttributes(new URL("http://www.example.com/thumbs/123-3.jpg"), "Grilling steaks for summer, episode 3",
                        "Alkis shows you how to get perfectly done steaks every time", new URL("http://www.example.com/video123-3.flv"), null);
        expectedVideoAttributes.setPrices(new VideoAttributes.VideoPrice[] { new VideoAttributes.VideoPrice(null, null, VideoAttributes.VideoPriceType.rent) });
        attr = (VideoAttributes) siter.next().getAttributesForExtension(Extension.VIDEO)[0];
        assertNotNull(attr);
        assertEquals(expectedVideoAttributes, attr);
    }

    @Test
    public void testImageSitemap() throws UnknownFormatException, IOException {
        SiteMapParser parser = new SiteMapParser();
        parser.enableExtension(Extension.IMAGE);

        String contentType = "text/xml";
        byte[] content = SiteMapParserTest.getResourceAsBytes("src/test/resources/sitemaps/extension/sitemap-images.xml");

        URL url = new URL("http://www.example.com/sitemap-images.xml");
        AbstractSiteMap asm = parser.parseSiteMap(contentType, content, url);
        assertEquals(false, asm.isIndex());
        assertEquals(true, asm instanceof SiteMap);
        SiteMap sm = (SiteMap) asm;
        assertEquals(1, sm.getSiteMapUrls().size());
        ImageAttributes imageAttributes1 = new ImageAttributes(new URL("http://example.com/image.jpg"));
        ImageAttributes imageAttributes2 = new ImageAttributes(new URL("http://example.com/photo.jpg"));
        imageAttributes2.setCaption("This is the caption.");
        imageAttributes2.setGeoLocation("Limerick, Ireland");
        imageAttributes2.setTitle("Example photo shot in Limerick, Ireland");
        imageAttributes2.setLicense(new URL("https://creativecommons.org/licenses/by/4.0/legalcode"));

        for (SiteMapURL su : sm.getSiteMapUrls()) {
            assertNotNull(su.getAttributesForExtension(Extension.IMAGE));
            ExtensionMetadata[] attrs = su.getAttributesForExtension(Extension.IMAGE);
            ImageAttributes attr = (ImageAttributes) attrs[0];
            assertEquals(imageAttributes1, attr);
            attr = (ImageAttributes) attrs[1];
            assertEquals(imageAttributes2, attr);
        }
    }

    @SuppressWarnings("serial")
    @Test
    public void testXHTMLLinksSitemap() throws UnknownFormatException, IOException, MalformedURLException {
        SiteMapParser parser = new SiteMapParser();
        parser.enableExtension(Extension.LINKS);

        String contentType = "text/xml";
        byte[] content = SiteMapParserTest.getResourceAsBytes("src/test/resources/sitemaps/extension/sitemap-links.xml");

        URL url = new URL("http://www.example.com/sitemap-links.xml");
        AbstractSiteMap asm = parser.parseSiteMap(contentType, content, url);
        assertEquals(false, asm.isIndex());
        assertEquals(true, asm instanceof SiteMap);
        SiteMap sm = (SiteMap) asm;
        assertEquals(3, sm.getSiteMapUrls().size());
        // all three pages share the same links attributes
        LinkAttributes[] linkAttributes = new LinkAttributes[] { new LinkAttributes(new URL("http://www.example.com/deutsch/")),
                        new LinkAttributes(new URL("http://www.example.com/schweiz-deutsch/")), new LinkAttributes(new URL("http://www.example.com/english/")) };
        linkAttributes[0].setParams(new HashMap<String, String>() {
            {
                put("rel", "alternate");
                put("hreflang", "de");
            }
        });
        linkAttributes[1].setParams(new HashMap<String, String>() {
            {
                put("rel", "alternate");
                put("hreflang", "de-ch");
            }
        });
        linkAttributes[2].setParams(new HashMap<String, String>() {
            {
                put("rel", "alternate");
                put("hreflang", "en");
            }
        });

        for (SiteMapURL su : sm.getSiteMapUrls()) {
            assertNotNull(su.getAttributesForExtension(Extension.LINKS));
            ExtensionMetadata[] attrs = su.getAttributesForExtension(Extension.LINKS);
            assertEquals(linkAttributes.length, attrs.length);
            for (int i = 0; i < linkAttributes.length; i++) {
                LinkAttributes attr = (LinkAttributes) attrs[i];
                assertEquals(linkAttributes[i], attr);
            }
        }
    }

    @Test
    public void testNewsSitemap() throws UnknownFormatException, IOException {
        SiteMapParser parser = new SiteMapParser();
        parser.enableExtension(Extension.NEWS);

        String contentType = "text/xml";
        byte[] content = SiteMapParserTest.getResourceAsBytes("src/test/resources/sitemaps/extension/sitemap-news.xml");

        URL url = new URL("http://www.example.org/sitemap-news.xml");
        AbstractSiteMap asm = parser.parseSiteMap(contentType, content, url);
        assertEquals(false, asm.isIndex());
        assertEquals(true, asm instanceof SiteMap);
        SiteMap sm = (SiteMap) asm;
        assertEquals(1, sm.getSiteMapUrls().size());
        ZonedDateTime dt = ZonedDateTime.parse("2008-11-23T00:00:00+00:00");
        NewsAttributes expectedNewsAttributes = new NewsAttributes("The Example Times", "en", dt, "Companies A, B in Merger Talks");
        expectedNewsAttributes.setKeywords(new String[] { "business", "merger", "acquisition", "A", "B" });
        expectedNewsAttributes.setGenres(new NewsAttributes.NewsGenre[] { NewsAttributes.NewsGenre.PressRelease, NewsAttributes.NewsGenre.Blog });
        expectedNewsAttributes.setStockTickers(new String[] { "NASDAQ:A", "NASDAQ:B" });
        for (SiteMapURL su : sm.getSiteMapUrls()) {
            assertNotNull(su.getAttributesForExtension(Extension.NEWS));
            NewsAttributes attr = (NewsAttributes) su.getAttributesForExtension(Extension.NEWS)[0];
            assertEquals(expectedNewsAttributes, attr);
        }
    }

    @Test
    public void testMobileSitemap() throws UnknownFormatException, IOException {
        SiteMapParser parser = new SiteMapParser();
        parser.enableExtension(Extension.MOBILE);

        String contentType = "text/xml";
        byte[] content = SiteMapParserTest.getResourceAsBytes("src/test/resources/sitemaps/extension/sitemap-mobile.xml");

        URL url = new URL("http://www.example.org/sitemap-mobile.xml");
        AbstractSiteMap asm = parser.parseSiteMap(contentType, content, url);
        assertEquals(false, asm.isIndex());
        assertEquals(true, asm instanceof SiteMap);
        SiteMap sm = (SiteMap) asm;
        for (SiteMapURL su : sm.getSiteMapUrls()) {
            URL u = su.getUrl();
            ExtensionMetadata[] attrs = su.getAttributesForExtension(Extension.MOBILE);
            if (u.getPath().contains("mobile-friendly")) {
                assertNotNull(attrs);
                MobileAttributes attr = (MobileAttributes) attrs[0];
                assertNotNull(attr);
            } else {
                assertTrue(attrs == null || attrs.length == 0);
            }
        }
    }

    @Test
    public void testShinpaideshuNewsSitemap() throws UnknownFormatException, IOException {
        SiteMapParser parser = new SiteMapParser();
        parser.setStrictNamespace(true);
        parser.enableExtension(Extension.NEWS);

        String contentType = "text/xml";
        byte[] content = SiteMapParserTest.getResourceAsBytes("src/test/resources/sitemaps/extension/shinpaideshou-news-sitemap.xml");

        URL url = new URL("https://shinpaideshou.wordpress.com/news-sitemap.xml");
        AbstractSiteMap asm = parser.parseSiteMap(contentType, content, url);
        assertEquals(false, asm.isIndex());
        assertEquals(true, asm instanceof SiteMap);
        SiteMap sm = (SiteMap) asm;
        assertEquals(3, sm.getSiteMapUrls().size());
        for (SiteMapURL su : sm.getSiteMapUrls()) {
            assertNotNull(su.getAttributesForExtension(Extension.NEWS));
            NewsAttributes attr = (NewsAttributes) su.getAttributesForExtension(Extension.NEWS)[0];
            assertNotNull(attr.getName());
            assertNotNull(attr.getPublicationDateTime());
            assertEquals(2017, attr.getPublicationDateTime().getYear());
        }
    }

    @Test
    public void testHebdenbridgetimesArticlesSitemap() throws UnknownFormatException, IOException {
        SiteMapParser parser = new SiteMapParser();
        parser.setStrictNamespace(true);
        parser.enableExtension(Extension.NEWS);
        parser.enableExtension(Extension.IMAGE);
        parser.enableExtension(Extension.VIDEO);
        parser.enableExtension(Extension.MOBILE);

        String contentType = "text/xml";
        byte[] content = SiteMapParserTest.getResourceAsBytes("src/test/resources/sitemaps/extension/hebdenbridgetimes-articles-sitemap.xml");

        URL url = new URL("http://www.hebdenbridgetimes.co.uk/sitemap-article-2015-18.xml");
        AbstractSiteMap asm = parser.parseSiteMap(contentType, content, url);
        assertEquals(false, asm.isIndex());
        assertEquals(true, asm instanceof SiteMap);
        SiteMap sm = (SiteMap) asm;
        assertEquals(74, sm.getSiteMapUrls().size());
    }

    @Test
    public void testImageAttributesAsMap() throws MalformedURLException {
        ImageAttributes attr = new ImageAttributes(new URL("http://example.com/image.jpg"));
        attr.setCaption("caption");
        attr.setGeoLocation("Limerick, Ireland");
        attr.setTitle("title");
        attr.setLicense(new URL("https://example.com/license"));

        java.util.Map<String, String[]> map = attr.asMap();
        assertArrayEquals(new String[] { "http://example.com/image.jpg" }, map.get("loc"));
        assertArrayEquals(new String[] { "caption" }, map.get("caption"));
        assertArrayEquals(new String[] { "Limerick, Ireland" }, map.get("geo_location"));
        assertArrayEquals(new String[] { "title" }, map.get("title"));
        assertArrayEquals(new String[] { "https://example.com/license" }, map.get("license"));
        assertEquals(5, map.size());

        // null attributes are not included
        ImageAttributes sparse = new ImageAttributes(new URL("http://example.com/image.jpg"));
        assertEquals(1, sparse.asMap().size());
    }

    @Test
    public void testLinkAttributesAsMap() throws MalformedURLException {
        LinkAttributes attr = new LinkAttributes(new URL("http://example.com/page"));
        java.util.Map<String, String> params = new HashMap<>();
        params.put("rel", "alternate");
        params.put("hreflang", "en");
        attr.setParams(params);

        java.util.Map<String, String[]> map = attr.asMap();
        assertArrayEquals(new String[] { "http://example.com/page" }, map.get("href"));
        assertArrayEquals(new String[] { "alternate" }, map.get("rel"));
        assertArrayEquals(new String[] { "en" }, map.get("hreflang"));
        assertEquals(3, map.size());
    }

    @Test
    public void testMobileAttributesAsMap() {
        assertTrue(new MobileAttributes().asMap().isEmpty());
    }

    @Test
    public void testNewsAttributesAsMap() {
        NewsAttributes attr = new NewsAttributes("Example Times", "en", ZonedDateTime.parse("2017-01-02T03:04:05Z"), "Title");
        attr.setGenres(new NewsAttributes.NewsGenre[] { NewsAttributes.NewsGenre.Blog, NewsAttributes.NewsGenre.OpEd });
        attr.setKeywords(new String[] { "keyword1", "keyword2" });
        attr.setStockTickers(new String[] { "GOOG" });

        java.util.Map<String, String[]> map = attr.asMap();
        assertArrayEquals(new String[] { "Example Times" }, map.get("publication/name"));
        assertArrayEquals(new String[] { "en" }, map.get("publication/language"));
        assertArrayEquals(new String[] { "Blog", "OpEd" }, map.get("genres"));
        assertArrayEquals(new String[] { "2017-01-02T03:04:05Z" }, map.get("publication_date"));
        assertArrayEquals(new String[] { "Title" }, map.get("title"));
        assertArrayEquals(new String[] { "keyword1", "keyword2" }, map.get("keywords"));
        assertArrayEquals(new String[] { "GOOG" }, map.get("stock_tickers"));
        assertEquals(7, map.size());
    }

    @Test
    public void testVideoAttributesAsMap() throws MalformedURLException {
        VideoAttributes attr = new VideoAttributes(new URL("http://example.com/thumb.jpg"), "Title", "Description", new URL("http://example.com/video.flv"),
                        new URL("http://example.com/player.swf"));
        attr.setDuration(600);
        attr.setRating(4.2f);
        attr.setFamilyFriendly(true);
        attr.setTags(new String[] { "tag1", "tag2" });
        attr.setAllowedCountries(new String[] { "IE", "US" });
        attr.setLive(false);

        java.util.Map<String, String[]> map = attr.asMap();
        assertArrayEquals(new String[] { "http://example.com/thumb.jpg" }, map.get("thumbnail_loc"));
        assertArrayEquals(new String[] { "Title" }, map.get("title"));
        assertArrayEquals(new String[] { "Description" }, map.get("description"));
        assertArrayEquals(new String[] { "http://example.com/video.flv" }, map.get("content_loc"));
        assertArrayEquals(new String[] { "http://example.com/player.swf" }, map.get("player_loc"));
        assertArrayEquals(new String[] { "600" }, map.get("duration"));
        assertArrayEquals(new String[] { "4.2" }, map.get("rating"));
        assertArrayEquals(new String[] { "true" }, map.get("family_friendly"));
        assertArrayEquals(new String[] { "tag1", "tag2" }, map.get("tags"));
        assertArrayEquals(new String[] { "IE", "US" }, map.get("allowed_countries"));
        assertArrayEquals(new String[] { "false" }, map.get("live"));
        assertEquals(11, map.size());
    }
}
