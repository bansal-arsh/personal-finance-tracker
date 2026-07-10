FROM ghcr.io/astral-sh/uv:python3.14-trixie-slim

ENV UV_NO_DEV=1

WORKDIR /python-supabase-api

COPY .python-version pyproject.toml uv.lock ./

RUN uv sync --locked

COPY . .

CMD ["uv", "run", "fastapi", "run"]
