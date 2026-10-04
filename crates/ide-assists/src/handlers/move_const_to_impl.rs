use hir::PathResolution;
use ide_db::{
    assists::{AssistId, AssistKind},
    defs::{Definition, NameRefClass},
};
use syntax::{
    ast::{self, edit::IndentLevel, HasName},
    AstNode, AstToken, TextRange,
};

use crate::assist_context::{AssistContext, Assists};

// Assist: move_const_to_impl
//
// Moves a const item declared inside a function to the impl block the function
// is defined in, replacing references to it with `Self::CONST`.
//
// ```
// struct Foo;
// impl Foo {
//     fn foo(&self) -> usize {
//         /// A doc comment attached.
//         const X$0YZ: usize = 1234;
//
//         XYZ
//     }
// }
// ```
// ->
// ```
// struct Foo;
// impl Foo {
//     /// A doc comment attached.
//     const XYZ: usize = 1234;
//
//     fn foo(&self) -> usize {
//         Self::XYZ
//     }
// }
// ```
pub(crate) fn move_const_to_impl(acc: &mut Assists, ctx: &AssistContext<'_>) -> Option<()> {
    let const_item = ctx.find_node_at_offset::<ast::Const>()?;
    let name = const_item.name()?;

    // The const has to be a statement directly inside a function body.
    let stmt_list = const_item.syntax().parent().and_then(ast::StmtList::cast)?;
    let body = stmt_list.syntax().parent().and_then(ast::BlockExpr::cast)?;
    let fn_ = body.syntax().parent().and_then(ast::Fn::cast)?;

    // The function has to be an item of an inherent impl block.
    let assoc_list = fn_.syntax().parent().and_then(ast::AssocItemList::cast)?;
    let impl_ = assoc_list.syntax().parent().and_then(ast::Impl::cast)?;
    if impl_.trait_().is_some() {
        // The const is not a member of the implemented trait.
        return None;
    }

    // Names of items declared alongside the const inside the function body.
    // Those would no longer be reachable from the impl block.
    let sibling_item_names: Vec<String> = stmt_list
        .statements()
        .filter_map(|stmt| match stmt {
            ast::Stmt::Item(item) => ast::AnyHasName::cast(item.syntax().clone()),
            _ => None,
        })
        .filter_map(|item| item.name())
        .map(|name| name.to_string())
        .collect();

    // The const must not reference anything that is only in scope inside the
    // function: locals, generic parameters or items declared in its body.
    let refers_to_fn_locals =
        const_item.syntax().descendants().filter_map(ast::Path::cast).any(|path| {
            let refers_to_sibling_item = first_segment_name(&path)
                .map_or(false, |name| sibling_item_names.iter().any(|it| *it == name));
            if refers_to_sibling_item {
                return true;
            }
            match ctx.sema.resolve_path(&path) {
                // Locals and generic parameters of the function do not resolve
                // from inside a const item, so treat unresolvable paths as
                // potential references to them.
                None => true,
                Some(res) => matches!(
                    res,
                    PathResolution::Local(_)
                        | PathResolution::TypeParam(_)
                        | PathResolution::ConstParam(_)
                ),
            }
        });
    if refers_to_fn_locals {
        return None;
    }

    // Collect the references to the const inside the function that can be
    // qualified with `Self::`. If the const is referenced from an item nested
    // inside the function body, moving it would break that reference.
    let const_def = ctx.sema.to_def(&const_item)?;
    let mut usages = Vec::new();
    for name_ref in fn_.syntax().descendants().filter_map(ast::NameRef::cast) {
        let refers_to_const = matches!(
            NameRefClass::classify(&ctx.sema, &name_ref),
            Some(NameRefClass::Definition(Definition::Const(it))) if it == const_def
        );
        if !refers_to_const {
            continue;
        }
        // `Self` is not available inside items nested in the function body.
        let is_nested = name_ref
            .syntax()
            .ancestors()
            .take_while(|it| it != fn_.syntax())
            .any(|it| ast::Item::can_cast(it.kind()));
        if is_nested {
            return None;
        }
        let segment = name_ref.syntax().parent().and_then(ast::PathSegment::cast)?;
        let path = segment.parent_path();
        if path.qualifier().is_some() || path.parent_path().is_some() {
            continue;
        }
        usages.push((path, name_ref.to_string()));
    }

    let target = const_item.syntax().text_range();
    acc.add(
        AssistId("move_const_to_impl", AssistKind::RefactorRewrite),
        format!("Move const {} to impl block", name),
        target,
        |builder| {
            // Qualify references to the const inside the function with `Self::`.
            for (path, name) in usages {
                builder.replace(path.syntax().text_range(), format!("Self::{}", name));
            }

            // Remove the const item, together with the line it sits on.
            let mut delete_range = const_item.syntax().text_range();
            if let Some(prev_ws) = const_item
                .syntax()
                .prev_sibling_or_token()
                .and_then(|it| it.into_token())
                .and_then(ast::Whitespace::cast)
            {
                delete_range =
                    TextRange::new(prev_ws.syntax().text_range().start(), delete_range.end());
            }
            builder.delete(delete_range);
            if let Some(next_ws) = const_item
                .syntax()
                .next_sibling_or_token()
                .and_then(|it| it.into_token())
                .and_then(ast::Whitespace::cast)
            {
                let text = next_ws.text();
                if let Some(pos) = text.rfind('\n') {
                    let new_text = format!("\n{}", &text[pos + 1..]);
                    if new_text != text {
                        builder.replace(next_ws.syntax().text_range(), new_text);
                    }
                }
            }

            // Insert the const at the top of the impl block.
            let item_indent = IndentLevel::from_node(impl_.syntax()) + 1;
            let const_indent = IndentLevel::from_node(const_item.syntax());
            let mut const_text = const_item.syntax().text().to_string();
            if const_indent.0 > item_indent.0 {
                const_text = const_text
                    .replace(&format!("\n{}", const_indent), &format!("\n{}", item_indent));
            }
            if let Some(l_curly) = assoc_list.l_curly_token() {
                builder.insert(
                    l_curly.text_range().end(),
                    format!("\n{}{}\n", item_indent, const_text),
                );
            }
        },
    )
}

fn first_segment_name(path: &ast::Path) -> Option<String> {
    let mut current = path.clone();
    loop {
        match current.qualifier() {
            Some(qualifier) => current = qualifier,
            None => return Some(current.segment()?.name_ref()?.to_string()),
        }
    }
}

#[cfg(test)]
mod tests {
    use crate::tests::{check_assist, check_assist_not_applicable};

    use super::*;

    #[test]
    fn move_const_to_impl_block() {
        check_assist(
            move_const_to_impl,
            r#"
struct Foo;
impl Foo {
    fn foo(&self) -> usize {
        /// A doc comment attached.
        const XY$0Z: usize = 1234;

        XYZ
    }
}
"#,
            r#"
struct Foo;
impl Foo {
    /// A doc comment attached.
    const XYZ: usize = 1234;

    fn foo(&self) -> usize {
        Self::XYZ
    }
}
"#,
        );
    }

    #[test]
    fn move_const_multiple_usages() {
        check_assist(
            move_const_to_impl,
            r#"
struct Foo;
impl Foo {
    fn foo() -> usize {
        const X$0YZ: usize = 1234;

        XYZ + XYZ
    }
}
"#,
            r#"
struct Foo;
impl Foo {
    const XYZ: usize = 1234;

    fn foo() -> usize {
        Self::XYZ + Self::XYZ
    }
}
"#,
        );
    }

    #[test]
    fn move_const_multiline_body() {
        check_assist(
            move_const_to_impl,
            r#"
struct Foo;
impl Foo {
    fn foo() -> (usize, usize) {
        const X$0YZ: (usize, usize) = (
            1234,
            5678,
        );

        XYZ
    }
}
"#,
            r#"
struct Foo;
impl Foo {
    const XYZ: (usize, usize) = (
        1234,
        5678,
    );

    fn foo() -> (usize, usize) {
        Self::XYZ
    }
}
"#,
        );
    }

    #[test]
    fn move_const_last_statement() {
        check_assist(
            move_const_to_impl,
            r#"
struct Foo;
impl Foo {
    fn foo() {
        let x = 1;
        const X$0YZ: usize = 1234;
    }
}
"#,
            r#"
struct Foo;
impl Foo {
    const XYZ: usize = 1234;

    fn foo() {
        let x = 1;
    }
}
"#,
        );
    }

    #[test]
    fn move_const_impl_in_mod() {
        check_assist(
            move_const_to_impl,
            r#"
mod m {
    struct Foo;
    impl Foo {
        fn foo() -> usize {
            const X$0YZ: usize = 1234;

            XYZ
        }
    }
}
"#,
            r#"
mod m {
    struct Foo;
    impl Foo {
        const XYZ: usize = 1234;

        fn foo() -> usize {
            Self::XYZ
        }
    }
}
"#,
        );
    }

    #[test]
    fn move_const_not_applicable_outside_impl() {
        check_assist_not_applicable(
            move_const_to_impl,
            r#"
fn foo() -> usize {
    const X$0YZ: usize = 1234;

    XYZ
}
"#,
        );
    }

    #[test]
    fn move_const_not_applicable_in_trait_impl() {
        check_assist_not_applicable(
            move_const_to_impl,
            r#"
struct Foo;
trait Tr {}
impl Tr for Foo {
    fn foo() -> usize {
        const X$0YZ: usize = 1234;

        XYZ
    }
}
"#,
        );
    }

    #[test]
    fn move_const_not_applicable_when_referencing_locals() {
        check_assist_not_applicable(
            move_const_to_impl,
            r#"
struct Foo;
impl Foo {
    fn foo() -> usize {
        let x = 1;
        const X$0YZ: usize = x;

        XYZ
    }
}
"#,
        );
    }

    #[test]
    fn move_const_not_applicable_when_referencing_generics() {
        check_assist_not_applicable(
            move_const_to_impl,
            r#"
struct Foo;
impl Foo {
    fn foo<T>() -> usize {
        const X$0YZ: usize = std::mem::size_of::<T>();

        XYZ
    }
}
"#,
        );
    }

    #[test]
    fn move_const_not_applicable_when_used_in_nested_item() {
        check_assist_not_applicable(
            move_const_to_impl,
            r#"
struct Foo;
impl Foo {
    fn foo() -> usize {
        const X$0YZ: usize = 1234;

        fn bar() -> usize {
            XYZ
        }
        bar() + XYZ
    }
}
"#,
        );
    }
}
