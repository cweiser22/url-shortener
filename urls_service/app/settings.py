import hashlib
import random

from pydantic_settings import BaseSettings
import uuid

class Settings(BaseSettings):
    mongo_uri: str
    redis_uri: str
    short_code_cache_ttl: int = 3600
    sliding_cache: bool = True
    url_db_name: str = "url_db"
    node_name: str


    class Config:
        env_file = ".env.local"


settings = Settings()
