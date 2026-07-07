import os
import sys
from typing import Annotated, ClassVar

from fastapi import Depends
from pydantic import BaseModel, ConfigDict
import supabase
from supabase.client import AsyncClient
from supabase_auth import AuthResponse


def getenv_with_err(key: str) -> str:
    val = os.getenv(key)
    if val is None:
        print(f"{key} environment variable not set")
        sys.exit(1)
    return val


async def create_supabase_async_client() -> AsyncClient:
    supabase_url: str = getenv_with_err("SUPABASE_URL")
    supabase_pub_key: str = getenv_with_err("SUPABASE_PUB_KEY")
    db_client: AsyncClient = await supabase.acreate_client(
        supabase_url,
        supabase_pub_key,
        options=supabase.AClientOptions(
            auto_refresh_token=False,
            persist_session=False,
        ),
    )
    return db_client


SupabaseClient = Annotated[AsyncClient, Depends(create_supabase_async_client)]


class Tokens(BaseModel):
    access_token: str
    refresh_token: str


class User(BaseModel):
    model_config: ClassVar[ConfigDict] = ConfigDict(arbitrary_types_allowed=True)

    db_client: AsyncClient
    tokens: Tokens


async def auth_db_with_tokens(db_client: SupabaseClient, tokens: Tokens) -> User:
    res: AuthResponse = await db_client.auth.set_session(
        tokens.access_token, tokens.refresh_token
    )
    new_access_token: str = (
        tokens.access_token if res.session is None else res.session.access_token
    )
    new_refresh_token: str = (
        tokens.refresh_token if res.session is None else res.session.refresh_token
    )
    return User(
        db_client=db_client,
        tokens=Tokens(
            access_token=new_access_token,
            refresh_token=new_refresh_token,
        ),
    )


AuthedUser = Annotated[User, Depends(auth_db_with_tokens)]
