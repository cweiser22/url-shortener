from contextlib import asynccontextmanager
from datetime import timedelta, datetime
from motor.motor_asyncio import AsyncIOMotorDatabase
from fastapi import FastAPI, Depends
from fastapi.responses import RedirectResponse, JSONResponse
from app import database, utils, schemas
from app.cache import get_cache
from redis.exceptions import ConnectionError as RedisConnectionError

import logging

app = FastAPI()
logging.basicConfig(level=logging.INFO)
log = logging.getLogger(__name__)

app.add_event_handler("shutdown", database.close_connection)


@app.get("/health")
async def health_check():
    return {"status": "ok"}


@app.post("/urls/", )
async def create_url_mapping(data: schemas.CreateURLMappingRequest, db: AsyncIOMotorDatabase = Depends(database.get_database)):
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


@app.get("/{code}")
async def redirect_to_long_url(code: str, db: AsyncIOMotorDatabase = Depends(database.get_database),
                               redis_client=Depends(get_cache)):

    # try and see if the url is cached
    cache_available = True
    try:
        cached_url = await redis_client.get(f"url-{code}")
        if cached_url:
            cached_url = cached_url.decode('utf-8')
            log.info(f"{code} was cached.")
            await redis_client.expire(f"url-{code}", 5)
            return RedirectResponse(url=cached_url, status_code=301, headers={
                'Location': cached_url
            })
    except (RedisConnectionError, TimeoutError):
        log.warning("Redis cache is not available. Performance may degrade.")
        cache_available = False

    print(f"{code} was not cached.")

    url_mapping = await db.get_collection('url_mappings').find_one({"short_code": code})
    if url_mapping is None:
        return {"message": "URL not found"}
    else:
        cache_exp_time = (datetime.now() + timedelta(minutes=5)).strftime("%a, %d %b %Y %H:%M:%S GMT")
        if cache_available:
            print("Caching url")
            await redis_client.set(f"url-{code}", url_mapping["long_url"], ex=5)
        return RedirectResponse(url=url_mapping["long_url"],
                                headers={
                                    "Cache-Control": "public, max-age=300",
                                    "Expires": cache_exp_time,
                                    "Location": url_mapping["long_url"]
                                },
                                status_code=301,
                                )

