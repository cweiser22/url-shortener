from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    influx_url: str
    influx_token: str
    influx_org: str
    influx_bucket: str

    class Config:
        env_file = ".env"


settings = Settings()
