use isahc::{prelude::*, Body, Request};
use testserver::mock;

#[test]
fn expect_header_is_sent_by_default_for_unknown_size_body() {
    let m = mock!();

    isahc::post(m.url(), Body::from_reader("hello world".as_bytes())).unwrap();

    m.request().expect_header("expect", "100-continue");
}

#[test]
fn expect_header_is_not_sent_when_disabled() {
    let m = mock!();

    Request::post(m.url())
        .expect_continue(false)
        .body(Body::from_reader("hello world".as_bytes()))
        .unwrap()
        .send()
        .unwrap();

    assert_eq!(m.request().get_header("expect").count(), 0);
}
