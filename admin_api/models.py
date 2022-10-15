from db_utils import Base as BaseModel, UUID
from sqlalchemy import Column, TIMESTAMP, Text, ForeignKey, text
from sqlalchemy.dialects.postgresql import INET, CIDR
from sqlalchemy.orm import relationship


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

