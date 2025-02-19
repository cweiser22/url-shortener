import logging

from motor.motor_asyncio import AsyncIOMotorClient, AsyncIOMotorDatabase
from .settings import settings

async_client = AsyncIOMotorClient(settings.mongo_uri)
async_db = async_client[settings.url_db_name]

log = logging.getLogger(__name__)


def get_database() -> AsyncIOMotorDatabase:
    return async_db


def close_connection():
    async_client.close()
    log.info("Connection to MongoDB closed")
