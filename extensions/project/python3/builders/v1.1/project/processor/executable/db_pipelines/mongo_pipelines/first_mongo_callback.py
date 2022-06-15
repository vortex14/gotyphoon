from typhoon.components.processor.executable.text_pipelines.base_pipeline import BasePipeline


class FirstCallback(BasePipeline):

    async def run(self):
        for doc in self.response.content:
            self.finish(doc, force_update=True)
        # self.switch_pipelines_group(test_product, "second_group")
        # self.finish(test_product)
        # callback = self.handler.first_group
