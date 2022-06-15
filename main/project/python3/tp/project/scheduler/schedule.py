from typhoon.components.scheduler.executions.base_schedule import BaseSchedule, ExecuteTime, ExecuteEvery
from datetime import datetime


class Schedule(BaseSchedule):

    @BaseSchedule.every(
        ExecuteEvery(minutes=1, active=True)
    )
    async def text_task(self, task):
        self.LOG.debug("project text_task", details={
            "now": str(datetime.now()),
            "taskid": task.taskid
        })


    @BaseSchedule.execute_at(
        ExecuteTime(time="02:02", timezone="MSK", active=True)
    )
    async def first_night(self, task):
        self.LOG.debug("first night", details={
            "now": str(datetime.now()),
            "taskid": task.taskid
        })


    @BaseSchedule.execute_at(
        ExecuteTime(time="08:05", timezone="MSK", active=True)
    )
    async def test(self, task):
        self.LOG.debug("test", details={
            "now": str(datetime.now()),
            "taskid": task.taskid
        })

    @BaseSchedule.execute_at(
        ExecuteTime(time="07:50", timezone="MSK", active=True)
    )
    async def test_morning(self, task):
        self.LOG.debug("test_morning", details={
            "now": str(datetime.now()),
            "taskid": task.taskid
        })

    @BaseSchedule.execute_at(
        ExecuteTime(time="04:13", timezone="MSK", active=True)
    )
    async def first_start_night(self, task):
        self.LOG.debug("first start night", details={
            "now": str(datetime.now()),
            "taskid": task.taskid
        })


    @BaseSchedule.execute_at(
        ExecuteTime(time="05:21", timezone="MSK", active=True)
    )
    async def first_start(self, task):
        self.LOG.debug("first_start", details={
            "now": str(datetime.now()),
            "taskid": task.taskid
        })


    @BaseSchedule.execute_at(
        ExecuteTime(time="07:00", timezone="MSK", active=True)
    )
    async def good_morning(self, task):
        self.LOG.debug("good_morning !", details={
            "now": str(datetime.now()),
            "taskid": task.taskid
        })

    