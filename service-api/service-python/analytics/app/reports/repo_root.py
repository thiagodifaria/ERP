"""Helpers para localizar o root do monorepo de forma resiliente."""

from pathlib import Path


def resolve_repo_root(start: Path) -> Path:
    """Resolve o root do repositorio procurando artefatos sentinela."""

    current = start.resolve()
    if current.is_file():
        current = current.parent

    for candidate in (current, *current.parents):
        if (candidate / "contracts").is_dir() and (candidate / "docs").is_dir():
            return candidate

    raise RuntimeError(f"repo root not found from {start}")


def try_resolve_repo_root(start: Path) -> Path | None:
    """Tenta resolver o root do repositorio sem estourar excecao."""

    try:
        return resolve_repo_root(start)
    except RuntimeError:
        return None
