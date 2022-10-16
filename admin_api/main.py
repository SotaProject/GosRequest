import httpagentparser
import logging

from aiocache import cached
from datetime import datetime
from fastapi import FastAPI
from sqlalchemy import Date, cast, func, text, select
from sqlalchemy.orm import selectinload

from admin_api import structures
from common import models
from common.db_utils import session_scope, prepare_db

app = FastAPI()
app.on_event('startup')(prepare_db)

logger = logging.getLogger(__name__)
logging.getLogger('sqlalchemy.engine').setLevel(logging.INFO)


@app.get('/')
async def root():
    return {'message': 'Yep, another microservice'}


@app.get('/subnets_data')
@cached(ttl=60)
async def data() -> structures.SubnetsDataResponse:
    async with session_scope() as session:
        query = (
            select(models.Subnet)
            .options(
                selectinload(models.Subnet.ranges),
                selectinload(models.Subnet.tags),
            )
        )

        subnets = (await session.execute(query)).scalars().all()

        return structures.SubnetsDataResponse(
            subnets=[
                structures.Subnet(
                    id=str(subnet.uuid),
                    name=subnet.name,
                    tags=[
                        tag.name
                        for tag in subnet.tags
                    ],
                    ranges=[
                        str(subnet_range.cidr)
                        for subnet_range in subnet.ranges
                    ]
                )
                for subnet in subnets
            ],
            last_updated=datetime.now().timestamp()
        )


@app.get('/fetch_notifications')
@cached(ttl=60)
async def fetch_notifications(tracker_uuid: str) -> structures.FetchNotificationsResponse:
    async with session_scope() as session:
        query = (
            select(models.Notification)
            .where(
                models.Notification.tracker_uuid == tracker_uuid,
                models.Notification.enable.is_(True),
            )
        )

        notifications: list[models.Notification] = (await session.execute(query)).scalars().all()

        return structures.FetchNotificationsResponse(
            chat_ids=[notification.chat_id for notification in notifications],
            last_updated=datetime.now().timestamp()
        )


@app.get('/statistics')
@cached(ttl=60)
async def statistics(start_datetime: datetime | None = None, end_datetime: datetime | None = None):
    if not start_datetime:
        start_datetime = datetime.min

    if not end_datetime:
        end_datetime = datetime.now()

    async with session_scope() as session:
        q = (
            select(models.Request.user_agent, func.count(models.Request.uuid).label('count'))
            .where(models.Request.created_at > start_datetime, models.Request.created_at < end_datetime)
            .group_by(models.Request.user_agent)
        )
        r = (await session.execute(q)).all()

        browsers = {}
        oss = {}
        for req in r:
            ua = httpagentparser.detect(req[0])
            if 'browser' not in ua.keys():
                browser_name = 'Other'
            else:
                browser_name = ua['browser']['name']
            if 'os' not in ua.keys():
                os_name = 'Other'
            else:
                os_name = ua['os']['name']

            if browser_name not in browsers.keys():
                browsers[browser_name] = req[1]
            else:
                browsers[browser_name] += req[1]
            if os_name not in oss.keys():
                oss[os_name] = req[1]
            else:
                oss[os_name] += req[1]

        browsers_p = {}
        os_p = {}
        for k, v in browsers.items():
            browsers_p[k] = v / sum(browsers.values())
        for k, v in oss.items():
            os_p[k] = v / sum(oss.values())

        q = (
            select(
                cast(models.Request.created_at, Date).label('date'),
                func.count().label('count'),
            )
            .filter(
                models.Request.created_at > start_datetime,
                models.Request.created_at < end_datetime
            )
            .group_by(cast(models.Request.created_at, Date))
            .order_by(text('date desc'))
        )
        dates = (await session.execute(q)).all()

        q = (
            select(models.Subnet.name, func.count(models.Request.uuid).label('count'))
            .where(models.Request.created_at > start_datetime, models.Request.created_at < end_datetime)
            .group_by(models.Request.subnet_uuid, models.Subnet.name)
            .join(models.Subnet.name, models.Subnet.uuid == models.Request.subnet_uuid)
        )
        subnets = (await session.execute(q)).all()

    return {
        'dates': dates,
        'subnets': subnets,
        'browsers': browsers_p,
        'os': os_p,
        'start_datetime': start_datetime,
        'end_datetime': end_datetime
    }
