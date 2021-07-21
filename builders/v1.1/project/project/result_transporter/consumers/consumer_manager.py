from project.result_transporter.consumers import mongo_production
from typhoon.extensions.result_transporter import exceptions, remove_exceptions


class ConsumerManager:

    async def counts(self):
        self.col = self.config.cache_storage.get_mongo_collection(collection_name="test", client_name="main")
        count = await self.col.count_documents({})
        # print(count)


    def __init__(self, config, loop):
        self.config = config
        self.loop = loop
        self.mongo_production = mongo_production.MongoProduction(self.config, self.loop)
        self.exceptions = exceptions.Exceptions(self.config, self.loop)
        self.remove_exceptions = remove_exceptions.RemoveExceptions(self.config, self.loop)
        self.loop.create_task(self.counts())


        self.consumers = {
            "mongo_production": {
                "consumer_instance": self.mongo_production,
                "bucket_required": False,
                "bucket_limit": 8
            },
            "remove_exceptions": {
                "bucket_required": False,
                "consumer_instance": self.remove_exceptions,
            },
            "exceptions": {
                "consumer_instance": self.exceptions,
                "bucket_required": False,
                "bucket_limit": 8
            },
            "R1": {
                "consumer_instance": None,
                "bucket_required": True
            }
        }
