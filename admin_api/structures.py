from pydantic import BaseModel


class Subnet(BaseModel):
    id: str
    name: str

    ranges: list[str]
    tags: list[str]


class SubnetsDataResponse(BaseModel):
    subnets: list[Subnet]

    last_updated: int


class FetchNotificationsResponse(BaseModel):
    chat_ids: list[str]
