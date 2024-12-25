import torch
import clip
from PIL import Image
import os
from io import BytesIO
import psycopg2

device = "cuda" if torch.cuda.is_available() else "cpu"
model, preprocess = clip.load("ViT-L/14", device=device)

db_conn = psycopg2.connect(
    user="postgres",
    password=os.getenv("DB_PASS"),
    host=os.getenv("DB_HOST"),
    port="5432",
    database="postgres"
)

cursor = db_conn.cursor()
cursor.execute("SELECT hash FROM image WHERE clip_embedding IS NULL;")
rows = cursor.fetchall()
result = [row[0] for row in rows]
cursor.close()

for hash in result:
    with open(os.getenv("BUCKET_PATH")+"/"+hash, "rb") as f:
        image = preprocess(
            Image.open(BytesIO(f.read()))
        ).unsqueeze(0).to(device)
        with torch.no_grad():
            image_features = model.encode_image(image)
            cursor = db_conn.cursor()
            cursor.execute(
                "UPDATE image SET clip_embedding = %s WHERE hash = %s;",
                (image_features.squeeze().tolist(), hash)
            )
            db_conn.commit()
            cursor.close()

db_conn.close()
