from typhoon.api.components.base import BaseComponent
from project.donor.v1.components import projects
from project.donor.v1.schemas import schemas

class V1(BaseComponent):
    def __init__(self, request, state, paths):
        super().__init__(request, state, paths)

        self.events = {
            "test": {
                "type": "POST",
                "method": self.test
            },
            "ping": {
                "type": "GET",
                "method": self.ping
            }
        }

        self.components = {
            "projects": projects.Projects
        }

        self.schemas = {
            "test": schemas.test_schema
        }


    async def ping(self):
        import time
        return {
            "status": True,
            "time": time.time()
        }

    async def test(self):
        return {
            "test": 12345678910
        }