from sqlalchemy import Column, BigInteger, TIMESTAMP, Text, BOOLEAN, ForeignKey, text
from sqlalchemy.dialects.postgresql import INET
from db_utils import Base as BaseModel
from db_utils import UUID
import uuid


class Request(BaseModel):
    """CREATE TABLE requests (
        uuid uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
        tracker_uuid uuid NOT NULL REFERENCES trackers (uuid),
        url text NOT NULL,
        ip inet NOT NULL,
        subnet_uuid uuid NOT NULL REFERENCES subnets (uuid),
        user_agent text NOT NULL,
        created_at timestamp DEFAULT now() NOT NULL
    );
    """
    __tablename__ = 'requests'
    uuid = Column(UUID, primary_key=True, server_default=text("uuid_generate_v4()"))
    tracker_uuid = Column(UUID, ForeignKey('trackers.uuid'), nullable=False)
    url = Column(Text, nullable=False)
    ip = Column(INET, nullable=False)
    subnet_uuid = Column(UUID, ForeignKey('subnets.uuid'), nullable=False)
    user_agent = Column(Text, nullable=False)
    created_at = Column(TIMESTAMP, nullable=False, server_default=text("now()"))


class Subnet(BaseModel):
    """
        CREATE TABLE subnets (
            uuid uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
            name text NOT NULL,
            tag_uuid uuid NOT NULL REFERENCES subnets_tags (uuid)
        );
    """
    __tablename__ = 'subnets'
    uuid = Column(UUID, primary_key=True, server_default=text("uuid_generate_v4()"))
    name = Column(Text, nullable=False)
    tag_uuid = Column(UUID, ForeignKey('subnets_tags.uuid'), nullable=False)
