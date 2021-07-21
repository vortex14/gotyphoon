from typhoon.components.processor.executable.text_pipelines.base_pipeline import BasePipeline


class SecondCallback(BasePipeline):

    async def run(self):
        pass
        # if not self.response.save["project"]:
        #     await self.second_callback(some_dict, saving)
        # else:
        #     print("ELSE in FIRST CALLBACK")

    async def second_callback(self, some_dict, saving):
        pass
