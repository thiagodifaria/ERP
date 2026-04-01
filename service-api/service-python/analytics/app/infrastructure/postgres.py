"""Acesso relacional basico para os relatórios do analytics."""

from collections.abc import Iterator
from contextlib import contextmanager

import psycopg
from psycopg.rows import dict_row

from app.config.settings import settings


@contextmanager
def connect() -> Iterator[psycopg.Connection]:
    connection = psycopg.connect(
        host=settings.postgres_host,
        port=settings.postgres_port,
        dbname=settings.postgres_db,
        user=settings.postgres_user,
        password=settings.postgres_password,
        sslmode=settings.postgres_ssl_mode,
        row_factory=dict_row,
    )

    try:
        yield connection
    finally:
        connection.close()


def postgres_ready() -> bool:
    if settings.repository_driver != "postgres":
        return False

    try:
        with connect() as connection:
            with connection.cursor() as cursor:
                cursor.execute("SELECT 1")
                cursor.fetchone()
        return True
    except psycopg.Error:
        return False
