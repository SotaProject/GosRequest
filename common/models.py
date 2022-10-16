from sqlalchemy import Column, TIMESTAMP, Text, BigInteger, Boolean, ForeignKey, text
from sqlalchemy.dialects.postgresql import INET, CIDR
from sqlalchemy.orm import relationship

from common.db_utils import Base as BaseModel, UUID


class Request(BaseModel):
    __tablename__ = 'requests'

    uuid = Column(UUID, primary_key=True, server_default=text("uuid_generate_v4()"))
    tracker_uuid = Column(UUID, ForeignKey('trackers.uuid'), nullable=False)
    url = Column(Text, nullable=False)
    ip = Column(INET, nullable=False)
    subnet_uuid = Column(UUID, ForeignKey('subnets.uuid'), nullable=False)
    user_agent = Column(Text, nullable=False)
    created_at = Column(TIMESTAMP, nullable=False, server_default=text("now()"))


class SubnetRange(BaseModel):
    __tablename__ = 'subnet_ranges'

    cidr = Column(CIDR, primary_key=True)
    subnet_uuid = Column(UUID, ForeignKey('subnets.uuid'), nullable=False)


class SubnetTag(BaseModel):
    __tablename__ = 'subnet_tags'

    uuid = Column(UUID, primary_key=True, server_default=text("uuid_generate_v4()"))
    name = Column(Text, nullable=False)
    subnet_uuid = Column(UUID, ForeignKey('subnets.uuid'), nullable=False)


class Subnet(BaseModel):
    __tablename__ = 'subnets'

    uuid = Column(UUID, primary_key=True, server_default=text("uuid_generate_v4()"))
    name = Column(Text, nullable=False)

    ranges = relationship(SubnetRange)
    tags = relationship(SubnetTag)


class Tracker(BaseModel):
    __tablename__ = 'trackers'
    uuid = Column(UUID, nullable=False, unique=True, primary_key=True)
    name = Column(Text, nullable=False)
    owner_id = Column(BigInteger, nullable=False)
    created_at = Column(TIMESTAMP, nullable=False)


class Notification(BaseModel):
    __tablename__ = 'notifications'
    uuid = Column(UUID, nullable=False, unique=True, primary_key=True, default=text("uuid_generate_v4()"))
    tracker_uuid = Column(UUID, ForeignKey('trackers.uuid'), nullable=False)
    chat_id = Column(BigInteger, nullable=False)
    enable = Column(Boolean, default=True)


class Users(BaseModel):
    __tablename__ = 'users'
    telegram_id = Column(BigInteger, nullable=False, unique=True, primary_key=True)
    action = Column(Text)
    state = Column(Text)
    created_at = Column(TIMESTAMP, nullable=False)
