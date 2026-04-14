from PIL import Image, ImageDraw
import os

def create_panorama(path, width=2048, height=1024):
    # Create a 2:1 image with a gradient
    img = Image.new('RGB', (width, height), color=(73, 109, 137))
    draw = ImageDraw.Draw(img)
    
    # Draw simple "sky" and "ground"
    draw.rectangle([0, 0, width, height // 2], fill=(135, 206, 235)) # Sky
    draw.rectangle([0, height // 2, width, height], fill=(34, 139, 34)) # Ground
    
    # Draw a "sun"
    draw.ellipse([width//2-50, height//4-50, width//2+50, height//4+50], fill=(255, 255, 0))
    
    # Draw some "clouds" or boxes to see texture
    for i in range(0, width, 200):
        draw.rectangle([i, height//2-20, i+100, height//2+20], fill=(200, 200, 200))

    img.save(path, quality=95)
    print(f"Test panorama created at {path}")

if __name__ == "__main__":
    test_dir = r"d:\lwdd\code\毕设\webGL-720yun\test\test_data"
    if not os.path.exists(test_dir):
        os.makedirs(test_dir)
    create_panorama(os.path.join(test_dir, "test_panorama_2k.jpg"))
