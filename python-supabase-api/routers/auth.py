from fastapi import APIRouter
from pydantic import BaseModel
from supabase_auth import (
    SignInWithEmailAndPasswordCredentials,
    SignUpWithEmailAndPasswordCredentials,
    VerifyTokenHashParams,
)

from routers.dependencies import SupabaseClient


router = APIRouter(prefix="/auth", tags=["Auth"])


class SignInRequest(BaseModel):
    email: str
    password: str


@router.post("/sign-in")
async def sign_in(req: SignInRequest, db_client: SupabaseClient):
    return await db_client.auth.sign_in_with_password(
        SignInWithEmailAndPasswordCredentials(
            email=req.email,
            password=req.password,
        )
    )


class ConfirmEmailRequest(BaseModel):
    email: str
    token_hash: str


@router.post("/confirm-email")
async def confirm_email(req: ConfirmEmailRequest, db_client: SupabaseClient):
    return db_client.auth.verify_otp(
        VerifyTokenHashParams(
            token_hash=req.token_hash,
            type="email",
        )
    )


class SignUpRequest(BaseModel):
    email: str
    password: str


@router.post("/sign-up")
async def sign_up(req: SignUpRequest, db_client: SupabaseClient):
    return db_client.auth.sign_up(
        SignUpWithEmailAndPasswordCredentials(
            email=req.email,
            password=req.password,
        )
    )
