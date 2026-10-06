from typing import Any, Union, Optional, NamedTuple
from sqlalchemy import inspect
from sqlalchemy.orm import DeclarativeMeta
from sqlalchemy.orm import DeclarativeBase, aliased
from sqlalchemy.sql import ColumnElement
from sqlalchemy.sql.schema import Column
from sqlalchemy.sql.elements import Label

from pydantic import BaseModel


class JoinConfig(NamedTuple):
    model: Any
    join_on: Any
    join_prefix: Optional[str] = None
    schema_to_select: Optional[type[BaseModel]] = None
    join_type: str = "left"


def _extract_matching_columns_from_schema(
    model: type[DeclarativeBase], schema: Optional[Union[type[BaseModel], list]]
) -> list[Any]:
    """
    Retrieves a list of ORM column objects from a SQLAlchemy model that match the field names in a given Pydantic schema.

    Args:
        model: The SQLAlchemy ORM model containing columns to be matched with the schema fields.
        schema: The Pydantic schema containing field names to be matched with the model's columns.

    Returns:
        A list of ORM column objects from the model that correspond to the field names defined in the schema.
    """
    column_list = list(inspect(model).selectable.columns)
    if schema is not None:
        if isinstance(schema, list):
            schema_fields = schema
        else:
            schema_fields = schema.model_fields.keys()

        column_list = []
        for column_name in schema_fields:
            if hasattr(model, column_name):
                column_list.append(getattr(model, column_name))

    return column_list


def _extract_matching_columns_from_kwargs(
    model: type[DeclarativeBase], kwargs: dict[str, Any]
) -> list[Any]:
    """
    Extracts matching ORM column objects from a SQLAlchemy model based on provided keyword arguments.

    Args:
        model: The SQLAlchemy ORM model.
        kwargs: A dictionary containing field names as keys.

    Returns:
        A list of ORM column objects from the model that correspond to the field names provided in kwargs.
    """
    if kwargs is not None:
        kwargs_fields = kwargs.keys()
        column_list = []
        for column_name in kwargs_fields:
            if hasattr(model, column_name):
                column_list.append(getattr(model, column_name))

    return column_list


def _extract_matching_columns_from_column_names(
    model: type[DeclarativeBase], column_names: list[str]
) -> list[Any]:
    """
    Extracts ORM column objects from a SQLAlchemy model based on a list of column names.

    Args:
        model: The SQLAlchemy ORM model.
        column_names: A list of column names to extract.

    Returns:
        A list of ORM column objects from the model that match the provided column names.
    """
    column_list = []
    for column_name in column_names:
        if hasattr(model, column_name):
            column_list.append(getattr(model, column_name))

    return column_list


def _auto_alias_joins(joins: list[JoinConfig]) -> list[JoinConfig]:
    """
    Assigns aliases to join configurations that reference the same model more than once,
    adapting their join conditions to the aliased model. This prevents SQL errors caused
    by the same table name appearing multiple times in a query.

    Args:
        joins: The list of JoinConfig instances to process.

    Returns:
        A list of JoinConfig instances where duplicated models are replaced by aliases.
    """
    model_counts: dict[Any, int] = {}
    for join in joins:
        model_counts[join.model] = model_counts.get(join.model, 0) + 1

    seen_counts: dict[Any, int] = {}
    aliased_joins = []
    for join in joins:
        if model_counts[join.model] > 1:
            index = seen_counts.get(join.model, 0)
            seen_counts[join.model] = index + 1
            alias = aliased(
                join.model, name=f"{join.model.__tablename__}_{index}"
            )
            join_on = join.join_on
            if join_on is not None:
                join_on = inspect(alias)._adapter.traverse(join_on)
            join = join._replace(model=alias, join_on=join_on)
        aliased_joins.append(join)

    return aliased_joins


def _auto_detect_join_condition(
    base_model: type[DeclarativeMeta], join_model: type[DeclarativeMeta]
) -> Optional[ColumnElement]:
    """
    Automatically detects the join condition for SQLAlchemy models based on foreign key relationships.

    Args:
        base_model: The base SQLAlchemy model from which to join.
        join_model: The SQLAlchemy model to join with the base model.

    Returns:
        A SQLAlchemy ColumnElement representing the join condition, if successfully detected.

    Raises:
        ValueError: If the join condition cannot be automatically determined.

    Example:
        # Assuming User has a foreign key reference to Tier:
        join_condition = auto_detect_join_condition(User, Tier)
    """
    fk_columns = [col for col in inspect(base_model).c if col.foreign_keys]
    join_on = next(
        (
            base_model.__table__.c[col.name]
            == join_model.__table__.c[list(col.foreign_keys)[0].column.name]
            for col in fk_columns
            if list(col.foreign_keys)[0].column.table == join_model.__table__
        ),
        None,
    )

    if join_on is None:
        raise ValueError(
            "Could not automatically determine join condition. Please provide join_on."
        )

    return join_on


def _add_column_with_prefix(column: Column, prefix: Optional[str]) -> Label:
    """
    Creates a SQLAlchemy column label with an optional prefix.

    Args:
        column: The SQLAlchemy Column object to be labeled.
        prefix: An optional prefix to prepend to the column's name.

    Returns:
        A labeled SQLAlchemy Column object.
    """
    column_label = f"{prefix}{column.name}" if prefix else column.name
    return column.label(column_label)
