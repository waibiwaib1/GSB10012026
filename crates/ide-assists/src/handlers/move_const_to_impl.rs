use ide_db::{
    assists::{AssistId, AssistKind},
    defs::{Definition, NameRefClass},
};
use syntax::{
    ast::{self, edit_in_place::Indent, HasName},
    AstNode, SyntaxKind,
};

use crate::assist_context::{AssistContext, Assists};

// Assist: move_const_to_impl
//
// Moves a local constant to an impl block, qualifying its references with `Self`.
//
// ```
// struct Foo;
//
// impl Foo {
//     fn test(&self) -> usize {
//         /// A doc comment attached.
//         const XY$0Z: usize = 1234;
//
//         XYZ
//     }
// }
// ```
// ->
// ```
// struct Foo;
//
// impl Foo {
//     /// A doc comment attached.
//     const XYZ: usize = 1234;
//
//     fn test(&self) -> usize {
//         Self::XYZ
//     }
// }
// ```
pub(crate) fn move_const_to_impl(acc: &mut Assists, ctx: &AssistContext<'_>) -> Option<()> {
    let const_item = ctx.find_node_at_offset::<ast::Const>()?;
    let name = const_item.name()?;

    // The const has to be a statement directly inside a function body.
    let stmt_list = const_item.syntax().parent().and_then(ast::StmtList::cast)?;
    let block = stmt_list.syntax().parent().and_then(ast::BlockExpr::cast)?;
    let fn_ = block.syntax().parent().and_then(ast::Fn::cast)?;

    // The function has to be an associated item of an inherent impl.
    let impl_ = fn_
        .syntax()
        .parent()
        .and_then(ast::AssocItemList::cast)?
        .syntax()
        .parent()
        .and_then(ast::Impl::cast)?;
    if impl_.trait_().is_some() {
        return None;
    }
    let assoc_item_list = impl_.assoc_item_list()?;

    // The impl must not already contain an item with the same name.
    let name_text = name.text();
    let name_conflicts = assoc_item_list.assoc_items().any(|item| {
        let item_name = match &item {
            ast::AssocItem::Const(it) => it.name(),
            ast::AssocItem::Fn(it) => it.name(),
            ast::AssocItem::TypeAlias(it) => it.name(),
            ast::AssocItem::MacroCall(_) => None,
        };
        item_name.map_or(false, |it| it.text() == name_text)
    });
    if name_conflicts {
        return None;
    }

    // The const must not depend on anything local to the function body, like
    // local variables or the function's generic parameters.
    let refers_to_fn_locals =
        const_item.syntax().descendants().filter_map(ast::NameRef::cast).any(|name_ref| {
            match NameRefClass::classify(&ctx.sema, &name_ref) {
                // `Self` keeps referring to the same impl after the move.
                Some(NameRefClass::Definition(Definition::Local(_)))
                | Some(NameRefClass::Definition(Definition::GenericParam(_)))
                | Some(NameRefClass::FieldShorthand { .. })
                | None => true,
                _ => false,
            }
        });
    if refers_to_fn_locals {
        return None;
    }

    let const_def = ctx.sema.to_def(&const_item)?;
    let target = const_item.syntax().text_range();

    acc.add(
        AssistId("move_const_to_impl", AssistKind::Refactor),
        "Move const to impl block",
        target,
        |builder| {
            // Qualify all references to the const with `Self`.
            let usages = Definition::Const(const_def).usages(&ctx.sema).all();
            if let Some(usages) = usages.references.get(&ctx.file_id()) {
                for usage in usages {
                    builder.replace(usage.range, format!("Self::{name_text}"));
                }
            }

            // Remove the const statement, including the leading whitespace of its
            // line and the line break after it.
            let delete_start = match const_item.syntax().prev_sibling_or_token() {
                Some(ws) if ws.kind() == SyntaxKind::WHITESPACE => ws.text_range().start(),
                _ => target.start(),
            };
            let delete_end = match const_item.syntax().next_sibling_or_token() {
                Some(ws) if ws.kind() == SyntaxKind::WHITESPACE => match ws.into_token() {
                    Some(ws) => match ws.text().find('\n') {
                        Some(idx) => {
                            ws.text_range().start() + syntax::TextSize::from(idx as u32 + 1)
                        }
                        None => ws.text_range().end(),
                    },
                    None => target.end(),
                },
                _ => target.end(),
            };
            let delete_range = syntax::TextRange::new(delete_start, delete_end);
            builder.delete(delete_range);

            // Insert the const into the impl block, right before the function.
            let fn_indent = syntax::ast::edit::IndentLevel::from_node(fn_.syntax());
            let moved_const = const_item.clone_for_update();
            moved_const.reindent_to(fn_indent);
            builder
                .insert(fn_.syntax().text_range().start(), format!("{moved_const}\n\n{fn_indent}"));
        },
    )
}

#[cfg(test)]
mod tests {
    use crate::tests::{check_assist, check_assist_not_applicable};

    use super::*;

    #[test]
    fn simple() {
        check_assist(
            move_const_to_impl,
            r#"
struct Foo;

impl Foo {
    fn test(&self) -> usize {
        const XY$0Z: usize = 1234;

        XYZ
    }
}
"#,
            r#"
struct Foo;

impl Foo {
    const XYZ: usize = 1234;

    fn test(&self) -> usize {
        Self::XYZ
    }
}
"#,
        );
    }

    #[test]
    fn moves_doc_comments_and_attrs() {
        check_assist(
            move_const_to_impl,
            r#"
struct Foo;

impl Foo {
    fn test(&self) -> usize {
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

    fn test(&self) -> usize {
        Self::XYZ
    }
}
"#,
        );
    }

    #[test]
    fn qualifies_multiple_usages() {
        check_assist(
            move_const_to_impl,
            r#"
struct Foo;

impl Foo {
    fn test(&self) -> usize {
        const XY$0Z: usize = 1234;

        XYZ + XYZ
    }
}
"#,
            r#"
struct Foo;

impl Foo {
    const XYZ: usize = 1234;

    fn test(&self) -> usize {
        Self::XYZ + Self::XYZ
    }
}
"#,
        );
    }

    #[test]
    fn multiline_const() {
        check_assist(
            move_const_to_impl,
            r#"
struct Foo;

impl Foo {
    fn test(&self) -> usize {
        const XY$0Z: usize = [
            1, 2, 3, 4,
        ][0];

        XYZ
    }
}
"#,
            r#"
struct Foo;

impl Foo {
    const XYZ: usize = [
        1, 2, 3, 4,
    ][0];

    fn test(&self) -> usize {
        Self::XYZ
    }
}
"#,
        );
    }

    #[test]
    fn not_applicable_outside_impl() {
        check_assist_not_applicable(
            move_const_to_impl,
            r#"
fn test() -> usize {
    const XY$0Z: usize = 1234;

    XYZ
}
"#,
        );
    }

    #[test]
    fn not_applicable_in_trait_impl() {
        check_assist_not_applicable(
            move_const_to_impl,
            r#"
struct Foo;
trait Bar {}

impl Bar for Foo {
    fn test(&self) -> usize {
        const XY$0Z: usize = 1234;

        XYZ
    }
}
"#,
        );
    }

    #[test]
    fn not_applicable_when_using_local() {
        check_assist_not_applicable(
            move_const_to_impl,
            r#"
struct Foo;

impl Foo {
    fn test(&self) -> usize {
        let local = 1234;
        const XY$0Z: usize = local;

        XYZ
    }
}
"#,
        );
    }

    #[test]
    fn not_applicable_when_name_conflicts() {
        check_assist_not_applicable(
            move_const_to_impl,
            r#"
struct Foo;

impl Foo {
    const XYZ: usize = 0;

    fn test(&self) -> usize {
        const XY$0Z: usize = 1234;

        XYZ
    }
}
"#,
        );
    }
}
