from pydantic import BaseModel, constr, Field
from ipaddress import IPv4Address


class SetWebhook(BaseModel):
    telegram_token: str
    route_ip: IPv4Address
    route_port: int = Field(ge=1024, le=65535)
    max_connections: int = Field(ge=1, le=100)
    drop_pending_updates: bool


class DeleteWebhook(BaseModel):
    telegram_token: str
