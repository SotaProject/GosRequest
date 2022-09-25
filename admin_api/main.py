from datetime import datetime
from typing import Union
from fastapi import FastAPI
from sqlalchemy.sql import select
from sqlalchemy import Date, cast, func, text
import httpagentparser
import logging

from db_utils import session_scope
import models

app = FastAPI()

logger = logging.getLogger(__name__)
logging.getLogger('sqlalchemy.engine').setLevel(logging.INFO)


@app.get("/")
async def root():
    return {"message": "Yep, another microservice"}


@app.get("/statistics")
async def statistics(start_datetime: Union[datetime, None],
                     end_datetime: Union[datetime, None]):
    async with session_scope() as session:
        q = select(models.Request.user_agent, models.Request.subnet_uuid) \
            .where(models.Request.created_at > start_datetime, models.Request.created_at < end_datetime)
        r = (await session.execute(q)).all()
        browsers = {}
        oss = {}
        for req in r:
            ua = httpagentparser.detect(req[0])
            browser_name = ua['browser']['name']
            os_name = ua['os']['name']
            if browser_name not in browsers.keys():
                browsers[browser_name] = 1
            else:
                browsers[browser_name] += 1
            if os_name not in oss.keys():
                oss[os_name] = 1
            else:
                oss[os_name] += 1

        browsers_p = {}
        os_p = {}
        for k, v in browsers.items():
            browsers_p[k] = v / sum(browsers.values())
        for k, v in oss.items():
            os_p[k] = v / sum(oss.values())

        q = select(
            cast(models.Request.created_at, Date).label('date'),
            func.count().label('count'),
        ).filter(
            models.Request.created_at > start_datetime, models.Request.created_at < end_datetime
        ).group_by(
            cast(models.Request.created_at, Date),
        ).order_by(
            text('date desc')
        )
        dates = (await session.execute(q)).all()

        q = select(models.Subnet.name, func.count(models.Request.uuid).label('count'))\
            .group_by(
            models.Request.subnet_uuid, models.Subnet.name
        ).join(
            models.Subnet.name, models.Subnet.uuid == models.Request.subnet_uuid
        )
        subnets = (await session.execute(q)).all()

    return {
        "dates": dates,
        "subnets": subnets,
        "browsers": browsers_p,
        "os": os_p,
        "start_datetime": start_datetime,
        "end_datetime": end_datetime
    }