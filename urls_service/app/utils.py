import hashlib

from snowflake import SnowflakeGenerator
from app.settings import settings

node_id = int(hashlib.md5(settings.node_name.encode()).hexdigest(), 16) % 1024
snowflake_generator = SnowflakeGenerator(node_id)


def generate_short_code():
    snowflake_id = next(snowflake_generator)
    return hex(snowflake_id)[2:]

