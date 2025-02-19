from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    mongo_uri: str
    redis_uri: str
    short_code_cache_ttl: int = 3600
    sliding_cache: bool = True
    url_db_name: str = "url_db"

    class Config:
        env_file = ".env.local"

settings = Settings()
