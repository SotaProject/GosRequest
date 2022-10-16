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
    tracker_name: str
    chat_ids: list[str]

    last_updated: int


class NewRequest(BaseModel):
    tracker_uuid: str
    url: str
    ip: str
    user_agent: str
    subnet_uuid: str
