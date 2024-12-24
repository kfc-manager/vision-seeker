import os


class Bucket:
    def __init__(self, path):
        self.path = path

    def get(self, key):
        with open(f"{self.path}/{key}", "rb") as f:
            return f.read()

    def delete(self, key):
        os.remove(f"{self.path}/{key}")
