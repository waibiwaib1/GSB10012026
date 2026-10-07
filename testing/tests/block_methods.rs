use rinja::Template;

/// Tests that `generate_block_methods` renders each block individually.
#[test]
fn test_block_methods_source() {
    #[derive(Template)]
    #[template(
        source = "{% block greeting %}Hello {{ name }}!{% endblock %}{% block farewell %}Bye {{ name }}!{% endblock %}",
        ext = "txt",
        generate_block_methods = true
    )]
    struct Greeting<'a> {
        name: &'a str,
    }

    let tmpl = Greeting { name: "world" };

    // The whole template still renders as usual.
    assert_eq!(tmpl.render().unwrap(), "Hello world!Bye world!");
    assert_eq!(tmpl.to_string(), "Hello world!Bye world!");

    // Each block can be rendered on its own.
    assert_eq!(tmpl.render_block_greeting().unwrap(), "Hello world!");
    assert_eq!(tmpl.render_block_farewell().unwrap(), "Bye world!");

    let mut buf = String::new();
    tmpl.render_block_greeting_into(&mut buf).unwrap();
    assert_eq!(buf, "Hello world!");

    buf.clear();
    tmpl.render_block_farewell_into_with_values(&mut buf, rinja::NO_VALUES)
        .unwrap();
    assert_eq!(buf, "Bye world!");
}

/// Tests block methods with template inheritance.
#[test]
fn test_block_methods_inheritance() {
    #[derive(Template)]
    #[template(path = "fragment-simple.html", generate_block_methods = true)]
    struct FragmentSimple<'a> {
        name: &'a str,
    }

    let simple = FragmentSimple { name: "world" };

    assert_eq!(
        simple.render_block_body().unwrap(),
        "\n<p>Hello world!</p>\n"
    );
    assert_eq!(
        simple.render_block_other_body().unwrap(),
        "\n<p>Don't render me.</p>\n"
    );
}

/// Tests that `super()` works inside block methods.
#[test]
fn test_block_methods_super() {
    #[derive(Template)]
    #[template(path = "fragment-super.html", generate_block_methods = true)]
    struct FragmentSuper<'a> {
        name: &'a str,
    }

    let sup = FragmentSuper { name: "world" };

    assert_eq!(
        sup.render_block_body().unwrap(),
        "\n<p>Hello world!</p>\n\n<p>Parent body content</p>\n\n"
    );
}

/// Tests that block methods work on the base template itself.
#[test]
fn test_block_methods_base() {
    #[derive(Template)]
    #[template(path = "fragment-base.html", generate_block_methods = true)]
    struct FragmentBase;

    let base = FragmentBase;

    assert_eq!(
        base.render_block_body().unwrap(),
        "\n<p>Parent body content</p>\n"
    );
    assert_eq!(base.render_block_other_body().unwrap(), "");
}

/// Tests block methods on a generic template without `extends`.
#[test]
fn test_block_methods_generics() {
    use std::fmt::Display;

    #[derive(Template)]
    #[template(
        source = "{% block content %}<{{ value }}>{% endblock %}",
        ext = "txt",
        generate_block_methods = true
    )]
    struct Wrapper<T: Display> {
        value: T,
    }

    let tmpl = Wrapper { value: 42 };
    assert_eq!(tmpl.render_block_content().unwrap(), "<42>");
    assert_eq!(tmpl.render().unwrap(), "<42>");
}

/// Tests that runtime [`Values`] are correctly forwarded in block methods.
#[test]
fn test_block_methods_values() {
    #[derive(Template)]
    #[template(
        source = "{% block content %}{% if let Ok(user) = \"user\" | value::<&str> %}Hello {{ user }}!{% else %}Hello!{% endif %}{% endblock %}",
        ext = "txt",
        generate_block_methods = true
    )]
    struct Greeting;

    let tmpl = Greeting;
    assert_eq!(tmpl.render_block_content().unwrap(), "Hello!");

    use std::any::Any;
    use std::collections::HashMap;

    let mut values: HashMap<String, Box<dyn Any>> = HashMap::default();
    values.insert("user".to_string(), Box::new("world"));
    assert_eq!(
        tmpl.render_block_content_with_values(&values).unwrap(),
        "Hello world!"
    );
}
