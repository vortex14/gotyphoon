from typhoon.components.processor.executable.pipeline_group import PipelinesGroup
from typhoon.components.processor.executable.base_handler import BaseHandler
from project.processor.executable.http_pipelines.text_pipelines import first_callback
from project.processor.executable.http_pipelines.text_pipelines import second_callback



class Handler(BaseHandler):

    async def init_attributes(self):
        self.url = "http://httpstat.us/200"
        self.first_group = PipelinesGroup("first_group", [
            first_callback.FirstCallback,
        ], ("mongo_production", ))


        self.second_group = PipelinesGroup("second_group", [
            second_callback.SecondCallback,
        ], ("mongo_production", ))

        self.loop.create_task(self.on_start({}))

    # @BaseHandler.execute_at("23:00", "UTC")
    @BaseHandler.every(minutes=1)
    async def on_start(self, task):


        await self.crawl(self.url, callback=self.first_group, force_update=True)

