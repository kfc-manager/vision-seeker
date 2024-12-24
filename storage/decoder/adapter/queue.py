import time
import pika


class Queue:
    def __init__(self, host, port, name, poll_interval=5):
        self.poll_interval = poll_interval
        self.name = name
        self.connection = pika.BlockingConnection(
            pika.ConnectionParameters(host=host, port=port)
        )
        self.channel = self.connection.channel()
        self.channel.queue_declare(queue=self.name)

    def close(self):
        self.channel.close()
        self.connection.close()

    def pull(self):
        while True:
            method, header, body = self.channel.basic_get(
                queue=self.name,
                auto_ack=True,
            )
            if body is None:
                time.sleep(self.poll_interval)
                continue
            return body.decode('utf-8')
