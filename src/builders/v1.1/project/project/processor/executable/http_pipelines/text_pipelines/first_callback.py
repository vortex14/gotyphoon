from typhoon.components.processor.executable.text_pipelines.base_pipeline import BasePipeline
import urllib.parse
from datetime import datetime
from typhoon.extensions.elogger import TyphoonLogger, ProcessorLogger, SchedulerLogger, FetcherLogger


class FirstCallback(BasePipeline):

    async def run(self):
        # PLOG = ProcessorLogger(ppl='first_group')
        # PLOG.debug(f"response.code:{self.response.code},test_at:{datetime.utcnow()}" )
        # self.switch_pipelines_group(test_product, "second_group")

        self.LOG.debug("++++++++++++++++++++++++")
        await self.finish({"status_code": self.response.code, "test_at": datetime.utcnow()})
        # callback = self.handler.first_group
