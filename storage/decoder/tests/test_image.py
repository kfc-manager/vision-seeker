import pytest
from domain.image import Image


@pytest.mark.parametrize(
    "input_path, expect",
    [
        ("./tests/images/non-trans.jpeg", {
            "format": "jpeg",
            "width": 16,
            "height": 12,
        }),
        ("./tests/images/non-trans.png", {
            "format": "png",
            "width": 16,
            "height": 12,
        }),
        ("./tests/images/trans.png", {
            "format": "png",
            "width": 16,
            "height": 16,
        }),
        ("./tests/images/non-trans.webp", {
            "format": "webp",
            "width": 16,
            "height": 12,
        }),
        ("./tests/images/trans.webp", {
            "format": "webp",
            "width": 16,
            "height": 16,
        }),
    ]
)
def test_init(input_path, expect):
    with open(input_path, "rb") as f:
        image = Image(f.read())
        assert image.format == expect["format"]
        assert image.width == expect["width"]
        assert image.height == expect["height"]
