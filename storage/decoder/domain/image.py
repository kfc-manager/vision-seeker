import io
import time
import hashlib
import numpy as np
import torch
import clip
from torchvision import transforms
from PIL import Image as Img


class Image:
    device = "cuda" if torch.cuda.is_available() else "cpu"
    dinov2 = {
        "model": None,
        "transform": None
    }
    clip = {
        "model": None,
        "preprocess": None,
    }

    @staticmethod
    def load_dinov2():
        Image.dinov2["model"] = torch.hub.load(
            "facebookresearch/dinov2",
            "dinov2_vitg14_reg",
        )
        Image.dinov2["model"].eval()
        Image.dinov2["transform"] = transforms.Compose([
            transforms.Resize((224, 224)),
            transforms.ToTensor(),
            transforms.Normalize(mean=[0.5, 0.5, 0.5],
                                 std=[0.5, 0.5, 0.5])
        ])

    @staticmethod
    def load_clip():
        Image.clip["model"], Image.clip["preprocess"] = clip.load(
            "ViT-L/14",
            device=Image.device,
        )

    def __init__(self, bytes):
        image = Img.open(io.BytesIO(bytes))
        self.format = image.format.lower()
        self.width, self.height = image.size
        self.size = len(bytes)
        self.bytes = bytes
        self.image = image
        self._hash = None
        self._entropy = None
        self._dinov2_embedding = None
        self._clip_embedding = None

    def hash(self):
        if self._hash is None:
            self._hash = hashlib.sha256(self.bytes).hexdigest()
        return self._hash

    def is_transparent(self):
        if self.image.mode in ("RGBA", "LA"):
            alpha = self.image.getchannel("A")
            if alpha.getextrema()[0] < 255:
                return True
            return False
        return False

    def entropy(self):
        if self._entropy is None:
            grayscale = np.array(self.image.convert('L'))
            flat = grayscale.flatten()
            hist, bins = np.histogram(
                flat,
                bins=256,
                range=(0, 255),
                density=True
            )
            # remove zero values (log(0) is undefined)
            hist = hist[hist > 0]
            entropy = -np.sum(hist * np.log2(hist))
            self._entropy = entropy if entropy >= 0.0 else 0.0
        return self._entropy

    def dinov2_embedding(self):
        if Image.dinov2["model"] is None or Image.dinov2["transform"] is None:
            raise RuntimeError("model or preprocessor of dinov2 not loaded")
        if self._dinov2_embedding is None:
            image = self.image.convert("RGB")
            input_tensor = Image.dinov2["transform"](image).unsqueeze(0)
            with torch.no_grad():
                self._dinov2_embedding = Image.dinov2["model"](input_tensor).\
                    squeeze().tolist()
        return self._dinov2_embedding

    def clip_embedding(self):
        if Image.clip["model"] is None or Image.clip["preprocess"] is None:
            raise RuntimeError("model or preprocessor of clip not loaded")
        if self._clip_embedding is None:
            input_tensor = Image.clip["preprocess"](self.image).\
                unsqueeze(0).to(Image.device)
            with torch.no_grad():
                self._clip_embedding = Image.clip["model"].\
                    encode_image(input_tensor).squeeze().tolist()
        return self._clip_embedding

    def to_dict(self):
        return {
            "hash": self.hash(),
            "format": self.format,
            "size": self.size,
            "width": self.width,
            "height": self.height,
            "entropy": self.entropy(),
            "dinov2_embedding": self.dinov2_embedding(),
            "clip_embedding": self.clip_embedding(),
            "created_at": time.time(),
        }
