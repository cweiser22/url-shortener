from pydantic import BaseModel, Field


class CreateURLMappingRequest(BaseModel):
    long_url: str = Field(..., description="The long URL to be shortened", alias="longUrl")

    class Config:
        populate_by_name = True
