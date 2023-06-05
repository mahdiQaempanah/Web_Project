from locust import HttpUser, task

# host is http://localhost:6433

class GatewayUser(HttpUser):

    @task
    def auth(self):
        self.client.get(url = "/api/v1/authz/pg" , json = { "nonce" : "sample", "message_id" : 10 })

    @task
    def biz_get(self):
        self.client.get(url = "/api/v1/biz/get" , json = { "userID": "99109999", "auth_key": -1, "message_id": 2 })
