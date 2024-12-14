import torch
import clip
from PIL import Image
import pika
import os
from io import BytesIO
import psycopg2


db_conn = psycopg2.connect(
    user=os.getenv("DB_USER"),
    password=os.getenv("DB_PASS"),
    host=os.getenv("DB_HOST"),
    port=os.getenv("DB_PORT"),
    database=os.getenv("DB_NAME")
)


def callback(ch, method, properties, body):
    hash = body.decode()
    with open(os.getenv("BUCKET_PATH")+"/"+hash, "rb") as f:
        image = preprocess(
            Image.open(BytesIO(f.read()))
        ).unsqueeze(0).to(device)
        with torch.no_grad():
            image_features = model.encode_image(image)
            cursor = db_conn.cursor()
            cursor.execute(
                "UPDATE image SET embedding = %s WHERE hash = %s;",
                (image_features.squeeze().tolist(), hash)
            )
            db_conn.commit()
            cursor.close()
        ch.basic_ack(delivery_tag=method.delivery_tag)
        return

    ch.basic_nack(delivery_tag=method.delivery_tag, requeue=True)


device = "cuda" if torch.cuda.is_available() else "cpu"
model, preprocess = clip.load("ViT-L/14", device=device)

# Establish connection to RabbitMQ
connection = pika.BlockingConnection(
    pika.ConnectionParameters(
        host=os.getenv("QUEUE_HOST"),
        port=os.getenv("QUEUE_PORT"),
    )
)
# Create a channel
channel = connection.channel()
# Declare the queue (ensures it exists)
channel.queue_declare(queue=os.getenv("QUEUE_NAME"))
# Set up consumer with callback
channel.basic_consume(
    queue=os.getenv("QUEUE_NAME"),
    on_message_callback=callback
)
# Start consuming
channel.start_consuming()
channel.close()
connection.close()
# Close DB connection
db_conn.close()
