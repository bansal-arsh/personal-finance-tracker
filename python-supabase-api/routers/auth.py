from fastapi import APIRouter
from pydantic import BaseModel
from supabase_auth import SignInWithEmailAndPasswordCredentials

from routers.dependencies import SupabaseClient


router = APIRouter(prefix="/auth", tags=["Auth"])


class SignInRequest(BaseModel):
    email: str
    password: str


@router.post("/sign-in")
async def sign_in(req: SignInRequest, db_client: SupabaseClient):
    return await db_client.auth.sign_in_with_password(
        SignInWithEmailAndPasswordCredentials(email=req.email, password=req.password)
    )
