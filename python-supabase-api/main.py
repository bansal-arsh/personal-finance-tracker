import os
import sys

from fastapi import FastAPI
import supabase
from dotenv import load_dotenv
from supabase.client import Client

env: str | None = os.getenv("ENV")
match env:
    case "DEV":
        _ = load_dotenv()
    case "PROD":
        pass
    case _:
        print("ENV environment variable not set")
        sys.exit(1)


def getenv_with_err(key: str) -> str:
    val = os.getenv(key)
    if val is None:
        print(f"{key} environment variable not set")
        sys.exit(1)
    return val


supabase_url: str = getenv_with_err("SUPABASE_URL")
supabase_pub_key: str = getenv_with_err("SUPABASE_PUB_KEY")
db_client: Client = supabase.Client(supabase_url, supabase_pub_key)

app = FastAPI()

@app.get("/")
async def root():
    return db_client.table("test").select("*").execute()
