import psycopg2


class Database:
    def __init__(self, host, port, name, user, pasw):
        self.connection = psycopg2.connect(
            user=user,
            password=pasw,
            host=host,
            port=port,
            database=name,
        )

    def close(self):
        self.connection.close()

    def insert_image(self, image_dict):
        query = """INSERT INTO image (hash, size, format, created_at,
            width, height, entropy, dino_embedding, clip_embedding)
            VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s);"""
        values = (
            image_dict["hash"],
            image_dict["size"],
            image_dict["format"],
            image_dict["created_at"],
            image_dict["width"],
            image_dict["height"],
            image_dict["entropy"],
            image_dict["dino_embedding"],
            image_dict["clip_embedding"],
        )
        cursor = self.connection.cursor()
        cursor.execute(query, values)
        self.connection.commit()
        cursor.close()
