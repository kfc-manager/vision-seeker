import os
from adapter.database import Database
from adapter.bucket import Bucket
from adapter.queue import Queue
from service.decoder import Decoder


def env_or_panic(key):
    val = os.getenv(key)
    if val is None:
        raise IOError(f"missing env variable '{key}'")
    return val


def main():
    db = Database(
        host=env_or_panic("DB_HOST"),
        port=env_or_panic("DB_PORT"),
        name=env_or_panic("DB_NAME"),
        user=env_or_panic("DB_USER"),
        pasw=env_or_panic("DB_PASS"),
    )
    bucket = Bucket(path=env_or_panic("BUCKET_PATH"))
    queue = Queue(
        host=env_or_panic("QUEUE_HOST"),
        port=env_or_panic("QUEUE_PORT"),
        name=env_or_panic("QUEUE_NAME"),
    )

    Decoder(
        db=db,
        bucket=bucket,
        queue=queue,
    ).decode()


if __name__ == "__main__":
    main()
