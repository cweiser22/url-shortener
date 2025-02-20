from redis.asyncio import BlockingConnectionPool, ConnectionPool, StrictRedis
from .settings import settings

redis_pool = BlockingConnectionPool.from_url(settings.redis_uri, max_connections=20, timeout=1)
redis_client = StrictRedis.from_pool(redis_pool)


async def get_cache():
    return redis_client

