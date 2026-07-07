import os
import sys

from fastapi import FastAPI
from dotenv import load_dotenv

from routers import auth

env: str | None = os.getenv("ENV")
match env:
    case "DEV":
        _ = load_dotenv()
    case "PROD":
        pass
    case _:
        print("ENV environment variable not set")
        sys.exit(1)

app = FastAPI()
app.include_router(auth.router)

# @app.post("/get-data")
# async def get_data(user: AuthedUser):
#     return await user.db_client.table("test").select("*").execute()
