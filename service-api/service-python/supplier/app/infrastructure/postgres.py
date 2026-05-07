from collections.abc import Iterator
from contextlib import contextmanager

import psycopg
from psycopg.rows import dict_row

from app.config.settings import settings


def connection_string() -> str:
    return (
        f"host={settings.postgres_host} "
        f"port={settings.postgres_port} "
        f"dbname={settings.postgres_db} "
        f"user={settings.postgres_user} "
        f"password={settings.postgres_password} "
        f"sslmode={settings.postgres_ssl_mode}"
    )


@contextmanager
def connect() -> Iterator[psycopg.Connection]:
    with psycopg.connect(connection_string(), row_factory=dict_row) as connection:
        yield connection


def postgres_ready() -> bool:
    try:
        with connect() as connection:
            with connection.cursor() as cursor:
                cursor.execute("SELECT 1")
                cursor.fetchone()
        return True
    except Exception:
        return False
