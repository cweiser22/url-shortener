from datetime import timedelta, datetime
from motor.motor_asyncio import AsyncIOMotorDatabase
from fastapi import FastAPI, Depends, APIRouter
from fastapi.responses import RedirectResponse, JSONResponse
from app import database, utils, schemas
from app.cache import get_cache
from redis.exceptions import ConnectionError as RedisConnectionError
from app.settings import settings
import logging

app = FastAPI()
logging.basicConfig(level=logging.INFO)
log = logging.getLogger(__name__)

app.add_event_handler("shutdown", database.close_connection)

log.info(f"Node: {settings.node_name}")

router = APIRouter()


@router.get("/health")
async def health_check():
    return {"status": "ok"}


@router.post("/mappings/", )
async def create_url_mapping(data: schemas.CreateURLMappingRequest, db: AsyncIOMotorDatabase = Depends(database.get_database)):
    print("making url")
    try:
        url_mapping_document = {
            "long_url": data.long_url,
            "short_code": utils.generate_short_code(),
        }
        result = await db.get_collection('url_mappings').insert_one(url_mapping_document)
        if result.acknowledged:
            return JSONResponse(content={"_id": str(result.inserted_id), "shortCode":
                url_mapping_document['short_code'], "longUrl": url_mapping_document['long_url']}, status_code=201,
                                headers={
                                    "Cache-Control": "no-cache"
                                }
                                )
        else:
            raise Exception("Failed to create URL mapping")
    except Exception as e:
        return JSONResponse(content={"message": "Something went wrong"}, status_code=400)


@router.get("/{code}/redirect")
async def redirect_to_long_url(code: str, db: AsyncIOMotorDatabase = Depends(database.get_database),
                               redis_client=Depends(get_cache)):

    # first, try and return the URL from the cache
    # cache_available is True if the cache if functioning properly
    # if the cache is down or unable to handle connections, we skip caching logic and serve the request from the db
    cache_available = True
    try:
        cached_url = await redis_client.get(f"url-{code}")
        if cached_url:
            cached_url = cached_url.decode('utf-8')
            log.debug(f"{code} was cached.")

            # we use a sliding TTL to keep frequently used URLs in the cache
            # so, we need to reset the TTL to 7 days upon each request
            await redis_client.expire(f"url-{code}", 60 * 60 * 24 * 7)
            return RedirectResponse(url=cached_url, status_code=301, headers={
                'Location': cached_url
            })
    # RedisConnectionError can be raised if the cache is getting too much traffic and hits the max_connections limit
    except (RedisConnectionError, TimeoutError):
        # we don't want to stop serving URLs if the cache is down, but we need to log a warning
        log.warning("Redis cache is not available. Performance may degrade.")
        cache_available = False

    log.debug(f"{code} was not cached.")


    # short_code is indexed, so it's efficient to query
    url_mapping = await db.get_collection('url_mappings').find_one({"short_code": code})
    if url_mapping is None:
        return {"message": "URL not found"}
    else:
        # cache for 7 days if the Redis is working
        cache_exp_time = (datetime.now() + timedelta(days=7)).strftime("%a, %d %b %Y %H:%M:%S GMT")
        if cache_available:
            log.debug(f"Caching url {code}")
            await redis_client.set(f"url-{code}", url_mapping["long_url"], ex=5)
        return RedirectResponse(url=url_mapping["long_url"],
                                headers={
                                    "Cache-Control": "public, max-age=300",
                                    "Expires": cache_exp_time,
                                    "Location": url_mapping["long_url"]
                                },
                                status_code=301,
                                )
app.include_router(router, prefix="/urls/api/v1")
