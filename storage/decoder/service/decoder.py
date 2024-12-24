from domain.image import Image


class Decoder:
    def __init__(self, db, bucket, queue):
        self.db = db
        self.bucket = bucket
        self.queue = queue
        Image.load_dinov2()
        Image.load_clip()

    def decode(self):
        while True:
            try:
                msg = self.queue.pull()
            except Exception:
                print("error on pulling message from queue")
                return

            b = self.bucket.get(msg)
            try:
                image = Image(b)
                if image.is_transparent():
                    raise Exception("image is transparent")
            except Exception:
                try:
                    self.bucket.delete(msg)
                except Exception:
                    print(f"error on deleting {msg}")
                    return
                continue

            try:
                self.db.insert_image(image.to_dict())
            except Exception:
                print(f"error on inserting image {msg} in database")
