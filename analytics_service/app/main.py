from time import timezone

from fastapi import FastAPI
from influxdb_client import InfluxDBClient, Point
from datetime import datetime, timezone
from pydantic_settings import BaseSettings
from influxdb_client.client.write_api import SYNCHRONOUS
from app.settings import settings


# Initialize InfluxDB client
client = InfluxDBClient(url=settings.influx_url, token=settings.influx_token, org=settings.influx_org)
write_api = client.write_api(write_options=SYNCHRONOUS)
query_api = client.query_api()

app = FastAPI()


# this is a route that marks every time a short URL is visited
# WARNING: this route is NOT idempotent, despite being a GET
# the reason it's a GET to simplify the request mirroring with Contour
@app.get("/{code}")
async def track_url_visit(code: str):
    """
    Records a visit for a given URL by:
    - Logging an individual visit (`url_visits`)
    - Updating the last visit timestamp (`last_visit`)
    """

    print(f"Tracking visit for code: {code}")

    now = datetime.now(timezone.utc)

    # Store the individual visit (append new record)
    visit_point = (
        Point("url_visits")
        .tag("code", code)
        .field("visits", 1)  # Each visit is a new record
        .time(now)
    )

    last_visit_point = (
        Point("last_visit")
        .tag("code", code)
        .field("timestamp", now.timestamp())  # Store as float (Unix time)
        .time(now)
    )

    # Batch write all points
    write_api.write(bucket=settings.influx_bucket, org=settings.influx_org, record=[visit_point, last_visit_point])

    return {"message": f"Visit recorded for {code}"}


@app.get("/analytics/url/{code}/stats/")
async def get_url_stats(code: str):
    """
    Retrieves the total visit count and last visited timestamp for a given URL.
    """

    flux_query = f"""
        from(bucket: "{settings.influx_bucket}")
          |> range(start: -30d)
          |> filter(fn: (r) => r["_measurement"] == "url_visits")
          |> filter(fn: (r) => r["code"] == "{code}")
          |> last()
          
        """

    result = query_api.query(flux_query)
    stats = {}

    for table in result:
        for record in table.records:
            stats["last_visited"] = record['_time']


    flux_query = f"""
    from(bucket: "{settings.influx_bucket}")
    |> range(start: -30d)
    |> filter(fn: (r) => r["_measurement"] == "url_visits")
    |> filter(fn: (r) => r["code"] == "{code}")
    |> count()
    """

    result = query_api.query(flux_query)
    for table in result:
        for record in table.records:
            stats["visits"] = record['_value']
    
    return stats if stats else {"message": "No data found for this URL"}


@app.get("/analytics/urls/inactive/")
async def get_inactive_urls():
    """
    Retrieves a list of all URLs that have not been visited in the last 30 days.
    """

    flux_query = f"""
    from(bucket: "{settings.influx_bucket}")
      |> range(start: -90d)  // Ensures old data is checked
      |> filter(fn: (r) => r["_measurement"] == "last_visit")
      |> filter(fn: (r) => r["_field"] == "timestamp")
      |> last()
      |> filter(fn: (r) => r["_value"] < float(v: now()) - 2592000.0) // 30 days in seconds
      |> keep(columns: ["code"])
    """

    result = query_api.query(flux_query)
    inactive_urls = [record["code"] for table in result for record in table.records]

    return {"inactive_urls": inactive_urls}
